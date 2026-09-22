package utils

import (
	"net/url"
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

func Benchmark_AppendQueryEscape(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"clean-64B", strings.Repeat("token-64", 8)},
		{"mixed-64B", strings.Repeat("a b&c=d ", 8)},
	}
	for _, input := range inputs {
		dst := make([]byte, 0, 3*len(input.value))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendQueryEscape(dst, input.value)
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
		b.Run(input.name+"/segments", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendPathSegmentsEscape(dst, input.value)
			}
		})
		// What a caller has to write without this helper.
		b.Run(input.name+"/per-segment", func(b *testing.B) {
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
