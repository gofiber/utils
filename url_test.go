package utils

import (
	"bytes"
	"math/rand/v2"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AppendQueryEscape(t *testing.T) {
	t.Parallel()
	// Exhaustive single-byte agreement with net/url pins the escape tables.
	for i := range 256 {
		in := string([]byte{byte(i)})
		require.Equal(t, url.QueryEscape(in), string(AppendQueryEscape(nil, in)), "byte %#x", i)
		require.Equal(t, url.PathEscape(in), string(AppendPathEscape(nil, in)), "byte %#x", i)
	}

	inputs := []string{
		"",
		"shortcuts",
		"hello world & goodbye",
		"a=1&b=2?c=/d",
		"héllo wörld",
		"100%+tax",
		"/products/42/reviews",
		strings.Repeat("clean-token.v2~x_", 8),
		"\x00\x1f\x7f\xff",
	}
	for _, in := range inputs {
		require.Equal(t, url.QueryEscape(in), string(AppendQueryEscape(nil, in)), "input %q", in)
		require.Equal(t, url.QueryEscape(in), string(AppendQueryEscape(nil, []byte(in))), "input %q", in)
		require.Equal(t, url.PathEscape(in), string(AppendPathEscape(nil, in)), "input %q", in)
		require.Equal(t, url.PathEscape(in), string(AppendPathEscape(nil, []byte(in))), "input %q", in)
	}

	// Appending must preserve existing dst content.
	got := AppendQueryEscape([]byte("q="), "a b")
	require.Equal(t, "q=a+b", string(got))
}

// escapeSegments is the reference: escape each segment on its own and rejoin.
func escapeSegments(s string) string {
	parts := strings.Split(s, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func Test_AppendPathSegmentsEscape(t *testing.T) {
	t.Parallel()
	// Exhaustive single-byte agreement pins the table, '/' included.
	for i := range 256 {
		in := string([]byte{byte(i)})
		require.Equal(t, escapeSegments(in), string(AppendPathSegmentsEscape(nil, in)), "byte %#x", i)
	}

	inputs := []string{
		"",
		"/",
		"//",
		"a/b",
		"/a/b/",
		"docs/guide/index.html",
		"a b/c&d/e=f",
		"héllo/wörld",
		"100%/tax",
		"/products/42/reviews",
		strings.Repeat("seg-1.v2~x_/", 8),
		"\x00\x1f/\x7f\xff",
	}
	for _, in := range inputs {
		want := escapeSegments(in)
		require.Equal(t, want, string(AppendPathSegmentsEscape(nil, in)), "input %q", in)
		require.Equal(t, want, string(AppendPathSegmentsEscape(nil, []byte(in))), "input %q", in)
	}

	// A slash is data for AppendPathEscape and a separator here.
	require.Equal(t, "a%2Fb", string(AppendPathEscape(nil, "a/b")))
	require.Equal(t, "a/b", string(AppendPathSegmentsEscape(nil, "a/b")))

	// Appending must preserve existing dst content.
	require.Equal(t, "/x/a%20b", string(AppendPathSegmentsEscape([]byte("/x/"), "a b")))
}

func Test_AppendQueryUnescape(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"",
		"plain",
		"a+b",
		"%41%42%43",
		"a%20b+c",
		"%2Fpath%2fmixed",
		"caf%C3%A9",
		"trailing%",
		"trailing%2",
		"%zz",
		"%1g",
		"%%41",
		"100%",
		strings.Repeat("no-escapes-here-at-all.", 8),
		strings.Repeat("%41+", 32),
	}
	for _, in := range inputs {
		wantQ, wantQErr := url.QueryUnescape(in)
		gotQ, errQ := AppendQueryUnescape(nil, in)
		gotQB, errQB := AppendQueryUnescape(nil, []byte(in))
		wantP, wantPErr := url.PathUnescape(in)
		gotP, errP := AppendPathUnescape(nil, in)
		if wantQErr != nil {
			require.Equal(t, wantQErr, errQ, "query input %q", in)
			require.Equal(t, wantQErr, errQB, "query input %q", in)
			require.Empty(t, gotQ, "query input %q", in)
		} else {
			require.NoError(t, errQ, "query input %q", in)
			require.Equal(t, wantQ, string(gotQ), "query input %q", in)
			require.Equal(t, wantQ, string(gotQB), "query input %q", in)
		}
		if wantPErr != nil {
			require.Equal(t, wantPErr, errP, "path input %q", in)
			require.Empty(t, gotP, "path input %q", in)
		} else {
			require.NoError(t, errP, "path input %q", in)
			require.Equal(t, wantP, string(gotP), "path input %q", in)
		}
	}

	// On error the original dst content survives with its original length.
	got, err := AppendQueryUnescape([]byte("keep"), "%zz")
	require.Error(t, err)
	require.Equal(t, "keep", string(got))

	// In-place decode: dst = src[:0] on the same backing array.
	buf := []byte("a%20b+c")
	out, err := AppendQueryUnescape(buf[:0], buf)
	require.NoError(t, err)
	require.Equal(t, "a b c", string(out))
}

// checkAppendOnly calls fn with a dst that holds a prefix and spare capacity full of marks, and
// fails unless the result is want appended to that prefix with every mark past it left alone:
// the helpers may not write beyond what they return.
func checkAppendOnly(tb testing.TB, what, in, want string, fn func(dst []byte) []byte) {
	tb.Helper()
	const prefix, mark = "keep:", 0xA5
	buf := bytes.Repeat([]byte{mark}, len(prefix)+3*len(in)+16)
	copy(buf, prefix)
	got := fn(buf[:len(prefix)])
	if string(got) != prefix+want {
		tb.Fatalf("%s(%q) = %q, want %q", what, in, got, prefix+want)
	}
	if &got[0] != &buf[0] {
		tb.Fatalf("%s(%q) left a dst that had the room", what, in)
	}
	if rest := buf[len(got):]; !bytes.Equal(rest, bytes.Repeat([]byte{mark}, len(rest))) {
		tb.Fatalf("%s(%q) wrote past its result: %q", what, in, rest)
	}
}

// Test_AppendEscape_Randomized compares the three escapers with net/url on inputs that escape
// rarely, often or always, into a nil dst, a full one and one with room to spare.
func Test_AppendEscape_Randomized(t *testing.T) {
	t.Parallel()
	const clean, dirty = "abcXYZ019-_.~", " /&=+%:@$;,?#\x00\x7f\x80\xc3\xa9\xff"
	rng := rand.New(rand.NewPCG(5, 6)) //nolint:gosec // a fixed sequence for a differential test
	for range 20_000 {
		in := make([]byte, rng.IntN(81))
		every := 1 << rng.IntN(7) // a byte to escape every 1 to 64 bytes
		for i := range in {
			if rng.IntN(every) == 0 {
				in[i] = dirty[rng.IntN(len(dirty))]
			} else {
				in[i] = clean[rng.IntN(len(clean))]
			}
		}
		for _, tc := range []struct {
			name string
			want string
			fn   func(dst []byte) []byte
		}{
			{"AppendQueryEscape", url.QueryEscape(string(in)), func(dst []byte) []byte { return AppendQueryEscape(dst, in) }},
			{"AppendPathEscape", url.PathEscape(string(in)), func(dst []byte) []byte { return AppendPathEscape(dst, string(in)) }},
			{"AppendPathSegmentsEscape", escapeSegments(string(in)), func(dst []byte) []byte { return AppendPathSegmentsEscape(dst, in) }},
		} {
			if got := tc.fn(nil); string(got) != tc.want {
				t.Fatalf("%s(nil, %q) = %q, want %q", tc.name, in, got, tc.want)
			}
			if got := tc.fn([]byte("full")[:4:4]); string(got) != "full"+tc.want {
				t.Fatalf("%s(full, %q) = %q, want %q", tc.name, in, got, "full"+tc.want)
			}
			checkAppendOnly(t, tc.name, string(in), tc.want, tc.fn)
		}
	}
}

// Test_AppendUnescape_Randomized compares the two decoders with net/url on inputs dense in
// escapes, valid and not, into a nil dst, one with room to spare, and in place.
func Test_AppendUnescape_Randomized(t *testing.T) {
	t.Parallel()
	pieces := []string{"%41", "%c3%A9", "%E4%BD%A0", "%2f", "+", "++", "%", "%4", "%zz", "%4g", "a", "bc", "plain-run-of-17b.", "x7~", "%25+"}
	rng := rand.New(rand.NewPCG(7, 8)) //nolint:gosec // a fixed sequence for a differential test
	for range 20_000 {
		var sb strings.Builder
		for range rng.IntN(12) {
			piece := pieces[rng.IntN(len(pieces))]
			if piece[0] == '%' && len(piece) < 3 && rng.IntN(8) != 0 {
				continue // keep most inputs valid
			}
			sb.WriteString(piece)
		}
		in := sb.String()
		for _, tc := range []struct {
			name string
			ref  func(string) (string, error)
			fn   func(dst, s []byte) ([]byte, error)
		}{
			{"AppendQueryUnescape", url.QueryUnescape, AppendQueryUnescape[[]byte]},
			{"AppendPathUnescape", url.PathUnescape, AppendPathUnescape[[]byte]},
		} {
			want, wantErr := tc.ref(in)
			got, err := tc.fn(nil, []byte(in))
			if wantErr != nil {
				if err == nil || err.Error() != wantErr.Error() || len(got) != 0 {
					t.Fatalf("%s(%q) = %q, %v, want %v", tc.name, in, got, err, wantErr)
				}
				kept, err := tc.fn([]byte("keep"), []byte(in))
				if err == nil || string(kept) != "keep" {
					t.Fatalf("%s(keep, %q) = %q, %v", tc.name, in, kept, err)
				}
				continue
			}
			if err != nil || string(got) != want {
				t.Fatalf("%s(%q) = %q, %v, want %q", tc.name, in, got, err, want)
			}
			checkAppendOnly(t, tc.name, in, want, func(dst []byte) []byte {
				out, err := tc.fn(dst, []byte(in))
				if err != nil {
					t.Fatalf("%s(%q) into a dst with room: %v", tc.name, in, err)
				}
				return out
			})
			inPlace := []byte(in)
			if out, err := tc.fn(inPlace[:0], inPlace); err != nil || string(out) != want {
				t.Fatalf("%s(%q) in place = %q, %v, want %q", tc.name, in, out, err, want)
			}
		}
	}
}

// Test_AppendUnescape_Clusters walks the byte loop that decodes an input whose first escape is
// followed at once by another: its scan after a run without escapes, with and without another
// escape behind it, a '+' in both modes, and malformed escapes in every place it reads them.
func Test_AppendUnescape_Clusters(t *testing.T) {
	t.Parallel()
	for _, in := range []string{
		"%E4%BD%A0%E5%A5%BD",
		"caf%C3%A9+cr%C3%A8me+br%C3%BBl%C3%A9e",
		"%C3%A9" + "1234567" + "%C3%A9" + "12345678" + "%C3%A9" + "123456789" + "%41",
		"%C3%A9" + "12345678", "%C3%A9" + "123456789", "%C3%A9" + "12345678" + "%", "ab%C3%A9" + "12345678" + "+%zz",
		"%C3%A9+%2B+%C3%A9",
		"%C3%A9%", "%C3%A9%4", "%C3%A9%4g", "%C3%A9%g4", "%C3%A9ab%", "%C3%A9ab%zz", "%C3%", "%C3%A", "%C3%zz",
		"%C3%A9" + strings.Repeat("x", 40) + "%zz",
		// Not clusters: the escape after the first one is not right behind it.
		"%C3a%A9", "%41+%42", "a%20b+cd%", "%20%",
	} {
		assertUnescapeParity(t, in)
		for _, fn := range []func(dst, s []byte) ([]byte, error){AppendQueryUnescape[[]byte], AppendPathUnescape[[]byte]} {
			want, wantErr := fn(nil, []byte(in))
			inPlace := []byte(in)
			got, err := fn(inPlace[:0], inPlace)
			require.Equal(t, wantErr, err, "input %q in place", in)
			require.Equal(t, string(want), string(got), "input %q in place", in)
			kept, err := fn([]byte("keep"), []byte(in))
			require.Equal(t, wantErr, err, "input %q", in)
			require.Equal(t, "keep"+string(want), string(kept), "input %q", in)
		}
	}
}

// Test_AppendEscape_GrowsOnce is not parallel: AllocsPerRun needs nothing else allocating. A
// dst too small for the input grows once, whatever the number of escapes; one with room for the
// result does not grow.
func Test_AppendEscape_GrowsOnce(t *testing.T) {
	plain := strings.Repeat("a b&c=d/\u00e9 ", 40)
	escaped := string(AppendQueryEscape(nil, plain))
	small, roomy := make([]byte, 0, 8), make([]byte, 0, 3*len(plain))
	var sink []byte
	// What growing once allocates is measured: one allocation, and two under the race detector.
	once := testing.AllocsPerRun(20, func() { sink = slices.Grow(small, 64) })
	for _, tc := range []struct {
		name string
		fn   func(dst []byte) []byte
	}{
		{"AppendQueryEscape", func(dst []byte) []byte { return AppendQueryEscape(dst, plain) }},
		{"AppendPathEscape", func(dst []byte) []byte { return AppendPathEscape(dst, plain) }},
		{"AppendPathSegmentsEscape", func(dst []byte) []byte { return AppendPathSegmentsEscape(dst, plain) }},
		{"AppendQueryUnescape", func(dst []byte) []byte { out, _ := AppendQueryUnescape(dst, escaped); return out }}, //nolint:errcheck // escaped is well formed
		{"AppendPathUnescape", func(dst []byte) []byte { out, _ := AppendPathUnescape(dst, escaped); return out }},   //nolint:errcheck // escaped is well formed
	} {
		require.InDelta(t, once, testing.AllocsPerRun(20, func() { sink = tc.fn(nil) }), 0, tc.name)
		require.InDelta(t, once, testing.AllocsPerRun(20, func() { sink = tc.fn(small) }), 0, tc.name)
		require.InDelta(t, 0, testing.AllocsPerRun(20, func() { sink = tc.fn(roomy) }), 0, tc.name)
		// Room for the result and no more is room enough.
		fit := make([]byte, 0, len(tc.fn(nil)))
		require.InDelta(t, 0, testing.AllocsPerRun(20, func() { sink = tc.fn(fit) }), 0, tc.name)
		// Grown to its size, not to the next power of two: the size classes are an eighth apart.
		require.Less(t, cap(tc.fn(nil)), len(tc.fn(nil))*5/4, tc.name)
	}

	// A space is one byte in a query and a %XX one byte decoded: neither is sized as three.
	require.Less(t, cap(AppendQueryEscape(nil, strings.Repeat(" ", 200))), 250)
	decoded, err := AppendPathUnescape(nil, strings.Repeat("%41", 200))
	require.NoError(t, err)
	require.Less(t, cap(decoded), 250)
	// 65 bytes are one more than a size class holds, so a size short by one would grow twice.
	edge := strings.Repeat("a", 61) + "/a"
	require.Len(t, AppendQueryEscape(nil, edge), 65)
	require.InDelta(t, once, testing.AllocsPerRun(20, func() { sink = AppendQueryEscape(nil, edge) }), 0)
	edge = strings.Repeat("a", 64) + "%41"
	require.InDelta(t, once, testing.AllocsPerRun(20, func() { sink, err = AppendQueryUnescape(nil, edge) }), 0)
	require.NoError(t, err)
	require.Len(t, sink, 65)
}

func Benchmark_AppendQueryEscape(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"clean-64B", strings.Repeat("token-64", 8)},
		{"mixed-64B", strings.Repeat("a b&c=d ", 8)},
		{"words-64B", strings.Repeat("some words here ", 4)},
	}
	for _, input := range inputs {
		// A dst that just holds the result, as a reused buffer comes to be: nothing to grow.
		dst := make([]byte, 0, len(url.QueryEscape(input.value)))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendQueryEscape(dst, input.value)
			}
		})
		// A dst without the room: what growing it costs.
		b.Run(input.name+"/fiber-grow", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendQueryEscape(nil, input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				url.QueryEscape(input.value)
			}
		})
	}
}

func Benchmark_AppendPathSegmentsEscape(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"1-segment", strings.Repeat("token-64", 3)},
		{"3-segments", "docs/guide/indexhtmlabc"},
		{"6-segments", "do/cs/gu/ide/index/html"},
		{"escaping", "do cs/gü ide/100%/html"},
	}
	for _, input := range inputs {
		dst := make([]byte, 0, 3*len(input.value))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendPathSegmentsEscape(dst, input.value)
			}
		})
		// What a caller has to write without this helper.
		b.Run(input.name+"/split-loop", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, start := dst, 0
				for i := range len(input.value) {
					if input.value[i] != '/' {
						continue
					}
					out = AppendPathEscape(out, input.value[start:i])
					out = append(out, '/')
					start = i + 1
				}
				_ = AppendPathEscape(out, input.value[start:])
			}
		})
	}
}

func Benchmark_AppendQueryUnescape(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"plain-64B", strings.Repeat("noescape", 8)},
		{"escaped-64B", strings.Repeat("a%20b+cd", 8)},
		{"utf8-62B", strings.Repeat("h%C3%A9llo+w%C3%B6rld%2C+", 2) + "h%C3%A9llo+w"},
	}
	for _, input := range inputs {
		dst := make([]byte, 0, len(input.value))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := AppendQueryUnescape(dst, input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
		// A dst without the room: what growing it costs.
		b.Run(input.name+"/fiber-grow", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := AppendQueryUnescape(nil, input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := url.QueryUnescape(input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func Benchmark_AppendPathUnescape(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"plain-25B", "/products/42/reviews/new!"},
		{"escaped-31B", "/products/caf%C3%A9%20bar/9%2F9"},
	}
	for _, input := range inputs {
		dst := make([]byte, 0, len(input.value))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := AppendPathUnescape(dst, input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := url.PathUnescape(input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
