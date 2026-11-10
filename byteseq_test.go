package utils

import (
	"bytes"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

const (
	upperStr = "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS/REPOS/GOFIBER/FIBER/ISSUES/COMMENTS"
	lowerStr = "/repos/gofiber/fiber/issues/187643/comments/repos/gofiber/fiber/issues/comments"
)

// Shared test cases for TrimSpace benchmarks
var trimSpaceBenchmarkCases = []struct {
	name  string
	input string
}{
	{name: "empty", input: ""},
	{name: "spaces", input: "   "},
	{name: "ascii-word", input: "  fiber  "},
	{name: "auth-header-bearer", input: "   Bearer abc.def.ghi   "},
	{name: "auth-header-basic", input: "\tBasic QWxhZGRpbjpvcGVuIHNlc2FtZQ==   "},
	{name: "accept-encoding", input: "  gzip, deflate, br  "},
	{name: "content-type", input: "  application/json  "},
	{name: "x-forwarded-for", input: " 203.0.113.5, 198.51.100.7  "},
	{name: "query-params", input: "  user=alice&id=42  "},
	{name: "ascii-long", input: "   " + strings.Repeat("fiber utils ", 8) + "   "},
	{name: "no-trim", input: "fiber"},
	{name: "mixed-whitespace", input: "\n\t fiber utils \r\n"},
}

func Test_EqualFold(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Expected bool
		S1       string
		S2       string
	}{
		{Expected: true, S1: "/MY/NAME/IS/:PARAM/*", S2: "/my/name/is/:param/*"},
		{Expected: true, S1: "/MY/NAME/IS/:PARAM/*", S2: "/my/name/is/:param/*"},
		{Expected: true, S1: "/MY1/NAME/IS/:PARAM/*", S2: "/MY1/NAME/IS/:PARAM/*"},
		{Expected: false, S1: "/my2/name/is/:param/*", S2: "/my2/name"},
		{Expected: false, S1: "/dddddd", S2: "eeeeee"},
		{Expected: false, S1: "\na", S2: "*A"},
		{Expected: true, S1: "/MY3/NAME/IS/:PARAM/*", S2: "/my3/name/is/:param/*"},
		{Expected: true, S1: "/MY4/NAME/IS/:PARAM/*", S2: "/my4/nAME/IS/:param/*"},
	}

	for _, tc := range testCases {
		res := EqualFold(tc.S1, tc.S2)
		require.Equal(t, tc.Expected, res, "string")

		res = EqualFold([]byte(tc.S1), []byte(tc.S2))
		require.Equal(t, tc.Expected, res, "bytes")
	}
}

func Benchmark_EqualFoldBytes(b *testing.B) {
	left := []byte(upperStr)
	right := []byte(lowerStr)
	var res bool
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = EqualFold(left, right)
		}
		require.True(b, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = bytes.EqualFold(left, right)
		}
		require.True(b, res)
	})
}

// go test -v -run=^$ -bench=Benchmark_EqualFold -benchmem -count=4 ./utils/
func Benchmark_EqualFold(b *testing.B) {
	var res bool
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = EqualFold(upperStr, lowerStr)
		}
		require.True(b, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = strings.EqualFold(upperStr, lowerStr)
		}
		require.True(b, res)
	})
}

func Test_TrimRight(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		S1     string
		S2     string
		Cutset byte
	}{
		{S1: "/test//////", S2: "/test", Cutset: '/'},
		{S1: "/test", S2: "/test", Cutset: '/'},
		{S1: " ", S2: "", Cutset: ' '},
		{S1: "  ", S2: "", Cutset: ' '},
		{S1: "", S2: "", Cutset: ' '},
	}

	for _, tc := range testCases {
		res := TrimRight(tc.S1, tc.Cutset)
		require.Equal(t, tc.S2, res, "string")

		resB := TrimRight([]byte(tc.S1), tc.Cutset)
		require.Equal(t, []byte(tc.S2), resB, "bytes")
	}
}

func Benchmark_TrimRight(b *testing.B) {
	var res string
	word := "foobar  "
	expected := "foobar"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = TrimRight(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = strings.TrimRight(word, " ")
		}
		require.Equal(b, expected, res)
	})
}

func Benchmark_TrimRightBytes(b *testing.B) {
	var res []byte
	word := []byte("foobar  ")
	expected := []byte("foobar")

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = TrimRight(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = bytes.TrimRight(word, " ")
		}
		require.Equal(b, expected, res)
	})
}

func Test_TrimLeft(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		S1     string
		S2     string
		Cutset byte
	}{
		{S1: "////test/", S2: "test/", Cutset: '/'},
		{S1: "test/", S2: "test/", Cutset: '/'},
		{S1: " ", S2: "", Cutset: ' '},
		{S1: "  ", S2: "", Cutset: ' '},
		{S1: "", S2: "", Cutset: ' '},
	}

	for _, tc := range testCases {
		res := TrimLeft(tc.S1, tc.Cutset)
		require.Equal(t, tc.S2, res, "string")

		resB := TrimLeft([]byte(tc.S1), tc.Cutset)
		require.Equal(t, []byte(tc.S2), resB, "bytes")
	}
}

func Benchmark_TrimLeft(b *testing.B) {
	var res string
	word := "  foobar"
	expected := "foobar"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = TrimLeft(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = strings.TrimLeft(word, " ")
		}
		require.Equal(b, expected, res)
	})
}

func Benchmark_TrimLeftBytes(b *testing.B) {
	var res []byte
	word := []byte("  foobar")
	expected := []byte("foobar")

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = TrimLeft(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = bytes.TrimLeft(word, " ")
		}
		require.Equal(b, expected, res)
	})
}

func Test_Trim(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		S1     string
		S2     string
		Cutset byte
	}{
		{S1: "   test  ", S2: "test", Cutset: ' '},
		{S1: "test", S2: "test", Cutset: ' '},
		{S1: ".test", S2: "test", Cutset: '.'},
		{S1: " ", S2: "", Cutset: ' '},
		{S1: "  ", S2: "", Cutset: ' '},
		{S1: "", S2: "", Cutset: ' '},
	}

	for _, tc := range testCases {
		res := Trim(tc.S1, tc.Cutset)
		require.Equal(t, tc.S2, res, "string")

		resB := Trim([]byte(tc.S1), tc.Cutset)
		require.Equal(t, []byte(tc.S2), resB, "bytes")
	}
}

func Benchmark_Trim(b *testing.B) {
	var res string
	word := "  foobar   "
	expected := "foobar"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = Trim(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = strings.Trim(word, " ")
		}
		require.Equal(b, expected, res)
	})
	b.Run("default.trimspace", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = strings.TrimSpace(word)
		}
		require.Equal(b, expected, res)
	})
}

func Benchmark_TrimBytes(b *testing.B) {
	var res []byte
	word := []byte("  foobar   ")
	expected := []byte("foobar")

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = Trim(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = bytes.Trim(word, " ")
		}
		require.Equal(b, expected, res)
	})
	b.Run("default.trimspace", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = bytes.TrimSpace(word)
		}
		require.Equal(b, expected, res)
	})
}

func Test_Trim_Edge(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		cut   byte
		exp   string
	}{
		{"foobar", 'x', "foobar"},
		{"", ' ', ""},
		{"xxfoobarxx", 'x', "foobar"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, c.exp, Trim(c.input, c.cut))
			require.Equal(t, []byte(c.exp), Trim([]byte(c.input), c.cut))
		})
	}
}

func Test_TrimSpace(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "spaces", input: "   "},
		{name: "tabs", input: "\t\t"},
		{name: "ascii-word", input: "  fiber  "},
		{name: "auth-header-bearer", input: "   Bearer abc.def.ghi   "},
		{name: "auth-header-basic", input: "\tBasic QWxhZGRpbjpvcGVuIHNlc2FtZQ==   "},
		{name: "accept-encoding", input: "  gzip, deflate, br  "},
		{name: "content-type-json", input: "  application/json  "},
		{name: "x-forwarded-for", input: " 203.0.113.5, 198.51.100.7  "},
		{name: "query-params", input: "  user=alice&id=42  "},
		{name: "ascii-long", input: "   " + strings.Repeat("fiber utils ", 8) + "   "},
		{name: "mixed-whitespace", input: "\n\t fiber utils \r\n"},
		{name: "no-trim", input: "fiber"},
		{name: "utf8-valid", input: "  Hello, 世界!  "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test string variant - must match strings.TrimSpace
			result := TrimSpace(tc.input)
			stdResult := strings.TrimSpace(tc.input)
			require.Equal(t, stdResult, result, "TrimSpace should match strings.TrimSpace for %s", tc.name)

			// Test []byte variant - must match bytes.TrimSpace
			resultBytes := TrimSpace([]byte(tc.input))
			stdResultBytes := bytes.TrimSpace([]byte(tc.input))
			// bytes.Equal treats nil and empty slices as equal
			require.True(t, bytes.Equal(stdResultBytes, resultBytes), "TrimSpace should match bytes.TrimSpace for %s", tc.name)
		})
	}
}

func Benchmark_TrimSpace(b *testing.B) {
	for _, tc := range trimSpaceBenchmarkCases {
		b.Run("fiber/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for i := 0; i < b.N; i++ {
				res = TrimSpace(tc.input)
			}
			_ = res
		})
		b.Run("default/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for i := 0; i < b.N; i++ {
				res = strings.TrimSpace(tc.input)
			}
			_ = res
		})
	}
}

func Benchmark_TrimSpaceBytes(b *testing.B) {
	for _, tc := range trimSpaceBenchmarkCases {
		input := []byte(tc.input)
		b.Run("fiber/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(input)))
			b.ResetTimer()
			var res []byte
			for i := 0; i < b.N; i++ {
				res = TrimSpace(input)
			}
			_ = res
		})
		b.Run("default/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(input)))
			b.ResetTimer()
			var res []byte
			for i := 0; i < b.N; i++ {
				res = bytes.TrimSpace(input)
			}
			_ = res
		})
	}
}

// Benchmark_EqualFold_Short measures token-sized inputs that are equal up to case, so every byte is inspected.
func Benchmark_EqualFold_Short(b *testing.B) {
	inputs := []struct {
		name string
		x, y string
	}{
		{"3B", "GET", "get"},
		{"4B", "POST", "post"},
		{"5B", "close", "Close"},
		{"7B", "chunked", "Chunked"},
		{"8B", "identity", "Identity"},
		{"10B", "keep-alive", "Keep-Alive"},
		{"16B", "application/json", "Application/JSON"},
	}
	var res bool
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res = EqualFold(input.x, input.y)
			}
			require.True(b, res)
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res = strings.EqualFold(input.x, input.y)
			}
			require.True(b, res)
		})
	}
}

// Test_EqualFold_Lengths covers every scan shape with a case-flipped twin and a changed byte at every position.
func Test_EqualFold_Lengths(t *testing.T) {
	t.Parallel()
	base := "Content-Type: Multipart"
	for n := 0; n <= len(base); n++ {
		a := base[:n]
		flipped := make([]byte, n)
		for i := range n {
			c := a[i]
			switch {
			case c >= 'a' && c <= 'z':
				c -= 0x20
			case c >= 'A' && c <= 'Z':
				c += 0x20
			}
			flipped[i] = c
		}
		require.True(t, EqualFold(a, string(flipped)), "len %d", n)
		require.True(t, EqualFold([]byte(a), flipped), "len %d", n)
		for i := range n {
			for _, c := range []byte{'!', '[', '{', '@', '`', 0x00, 0x80, 0xFF, a[i] ^ 0x01} {
				if c == a[i] {
					continue
				}
				mutated := bytes.Clone(flipped)
				mutated[i] = c
				want := asciiFoldUpper(a) == asciiFoldUpper(string(mutated))
				require.Equal(t, want, EqualFold(a, string(mutated)), "len %d pos %d byte %#x", n, i, c)
				require.Equal(t, want, EqualFold([]byte(a), mutated), "len %d pos %d byte %#x", n, i, c)
			}
		}
	}
	// '[' and '{' differ only in bit 5, like a letter pair, and must not fold.
	for n := 1; n <= 16; n++ {
		require.False(t, EqualFold(strings.Repeat("[", n), strings.Repeat("{", n)), "len %d", n)
		require.False(t, EqualFold(strings.Repeat("@", n), strings.Repeat("`", n)), "len %d", n)
	}
}

func forEachString(alphabet string, maxLen int, fn func(s string)) {
	buf := make([]byte, 0, maxLen)
	var walk func(depth int)
	walk = func(depth int) {
		fn(string(buf))
		if depth == 0 {
			return
		}
		for i := range len(alphabet) {
			buf = append(buf, alphabet[i])
			walk(depth - 1)
			buf = buf[:len(buf)-1]
		}
	}
	walk(maxLen)
}

// requireSameView fails unless got is want as the stdlib returns it: same bytes, nil-ness and window.
func requireSameView(tb testing.TB, what string, in, got, want []byte) {
	tb.Helper()
	switch {
	case (got == nil) != (want == nil):
		tb.Fatalf("%s(%q): result nil = %v, stdlib nil = %v", what, in, got == nil, want == nil)
	case !bytes.Equal(got, want):
		tb.Fatalf("%s(%q) = %q, stdlib gives %q", what, in, got, want)
	case cap(got) != cap(want):
		tb.Fatalf("%s(%q): capacity %d, stdlib's %d", what, in, cap(got), cap(want))
	case len(got) > 0 && &got[0] != &want[0]:
		tb.Fatalf("%s(%q): not the stdlib's window of the input", what, in)
	}
}

// requireTrimmed is requireSameView except for an empty result: bytes.Trim, TrimLeft and TrimSpace
// return nil there (undocumented), these return an empty slice of the input, keeping TrimSpace inlinable.
func requireTrimmed(tb testing.TB, what string, in, got, want []byte) {
	tb.Helper()
	if len(want) > 0 {
		requireSameView(tb, what, in, got, want)
		return
	}
	if len(got) != 0 {
		tb.Fatalf("%s(%q) = %q, stdlib gives nothing", what, in, got)
	}
	if (got == nil) != (in == nil) {
		tb.Fatalf("%s(%q): an empty result must be a slice of the input, so nil only for nil; got nil = %v for input nil = %v", what, in, got == nil, in == nil)
	}
}

func refTrimLeft(s string, c byte) string {
	for len(s) > 0 && s[0] == c {
		s = s[1:]
	}
	return s
}

func refTrimRight(s string, c byte) string {
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

func Test_Trim_StdlibParity(t *testing.T) {
	t.Parallel()
	for cut := range 256 {
		c := byte(cut)
		cutset := string([]byte{c})
		if c < utf8.RuneSelf {
			// Empty and nil inputs: bytes.TrimRight returns its input, the others nil.
			for _, in := range [][]byte{{}, nil} {
				requireTrimmed(t, "TrimLeft", in, TrimLeft(in, c), bytes.TrimLeft(in, cutset))
				requireSameView(t, "TrimRight", in, TrimRight(in, c), bytes.TrimRight(in, cutset))
				requireTrimmed(t, "Trim", in, Trim(in, c), bytes.Trim(in, cutset))
			}
		}
		forEachString(string([]byte{c, 'a', ' ', 0xe9}), 5, func(s string) {
			wantLeft, wantRight, wantBoth := refTrimLeft(s, c), refTrimRight(s, c), refTrimRight(refTrimLeft(s, c), c)
			if c < utf8.RuneSelf {
				// Below 0x80 the stdlib's cutset is the byte; above it, malformed UTF-8 matching every malformed byte.
				if strings.TrimLeft(s, cutset) != wantLeft || strings.TrimRight(s, cutset) != wantRight || strings.Trim(s, cutset) != wantBoth {
					t.Fatalf("reference trims disagree with the stdlib for %q cutset %q", s, cutset)
				}
				b := []byte(s)
				requireTrimmed(t, "TrimLeft", b, TrimLeft(b, c), bytes.TrimLeft(b, cutset))
				requireSameView(t, "TrimRight", b, TrimRight(b, c), bytes.TrimRight(b, cutset))
				requireTrimmed(t, "Trim", b, Trim(b, c), bytes.Trim(b, cutset))
			}
			if got := TrimLeft(s, c); got != wantLeft {
				t.Fatalf("TrimLeft(%q, %#x) = %q, want %q", s, c, got, wantLeft)
			}
			if got := TrimRight(s, c); got != wantRight {
				t.Fatalf("TrimRight(%q, %#x) = %q, want %q", s, c, got, wantRight)
			}
			if got := Trim(s, c); got != wantBoth {
				t.Fatalf("Trim(%q, %#x) = %q, want %q", s, c, got, wantBoth)
			}
			b := []byte(s)
			if got := TrimLeft(b, c); string(got) != wantLeft || (len(got) > 0 && &got[0] != &b[len(b)-len(wantLeft)]) {
				t.Fatalf("TrimLeft(bytes %q, %#x) = %q, want the tail %q", s, c, got, wantLeft)
			}
			if got := TrimRight(b, c); string(got) != wantRight || cap(got) != cap(b) {
				t.Fatalf("TrimRight(bytes %q, %#x) = %q (cap %d), want the head %q", s, c, got, cap(got), wantRight)
			}
			if got := Trim(b, c); string(got) != wantBoth {
				t.Fatalf("Trim(bytes %q, %#x) = %q, want %q", s, c, got, wantBoth)
			}
		})
	}
}

func Test_TrimSpace_StdlibParity(t *testing.T) {
	t.Parallel()
	// Whitespace and its neighbors in a table or range test.
	forEachString("\t\n\v\f\r \x00\x08\x0e\x1f\x7f!a", 5, func(s string) {
		if got, want := TrimSpace(s), strings.TrimSpace(s); got != want {
			t.Fatalf("TrimSpace(%q) = %q, stdlib gives %q", s, got, want)
		}
		b := []byte(s)
		requireTrimmed(t, "TrimSpace", b, TrimSpace(b), bytes.TrimSpace(b))
	})
	requireTrimmed(t, "TrimSpace(empty)", []byte{}, TrimSpace([]byte{}), bytes.TrimSpace([]byte{}))
	requireTrimmed(t, "TrimSpace(nil)", nil, TrimSpace([]byte(nil)), bytes.TrimSpace(nil))

	type word string
	type blob []byte
	if got := TrimSpace(word(" \tx\n")); got != "x" {
		t.Fatalf("TrimSpace(word) = %q", got)
	}
	if got := TrimSpace(blob(" \tx\n")); string(got) != "x" {
		t.Fatalf("TrimSpace(blob) = %q", got)
	}
}

// Test_TrimSpace_OnlyASCIIWhitespace pins the deliberate difference: only the six ASCII whitespace
// bytes are trimmed, not Unicode spaces such as U+00A0.
func Test_TrimSpace_OnlyASCIIWhitespace(t *testing.T) {
	t.Parallel()
	for r := rune(utf8.RuneSelf); r <= unicode.MaxRune; r++ {
		if !unicode.IsSpace(r) {
			continue
		}
		in := " " + string(r) + "x" + string(r) + " "
		want := string(r) + "x" + string(r)
		if got := TrimSpace(in); got != want {
			t.Fatalf("TrimSpace(%q) = %q, want %q", in, got, want)
		}
		if got := TrimSpace([]byte(in)); string(got) != want {
			t.Fatalf("TrimSpace(bytes %q) = %q, want %q", in, got, want)
		}
		if strings.TrimSpace(in) != "x" {
			t.Fatalf("the stdlib no longer trims U+%04X: the documented difference is gone", r)
		}
	}
	// A lone byte above 0x7F is no space to either (0x85, 0xA0 too): the stdlib reads it as malformed UTF-8.
	for b := utf8.RuneSelf; b < 256; b++ {
		for _, in := range []string{string([]byte{byte(b)}), " " + string([]byte{byte(b)}) + " ", "x" + string([]byte{byte(b)})} {
			if got, want := TrimSpace(in), strings.TrimSpace(in); got != want {
				t.Fatalf("TrimSpace(%q) = %q, stdlib gives %q", in, got, want)
			}
		}
	}
	if got := TrimSpace("\u00a0 x \u00a0"); got != "\u00a0 x \u00a0" {
		t.Fatalf("TrimSpace left %q", got)
	}
}

func Test_EqualFold_StdlibParity(t *testing.T) {
	t.Parallel()
	const edges = "aAzZ@[`{1"
	forEachString(edges, 3, func(a string) {
		forEachString(edges, 3, func(b string) {
			want := strings.EqualFold(a, b)
			if got := EqualFold(a, b); got != want {
				t.Fatalf("EqualFold(%q, %q) = %v, stdlib gives %v", a, b, got, want)
			}
			if got := EqualFold([]byte(a), []byte(b)); got != want {
				t.Fatalf("EqualFold(bytes %q, %q) = %v, stdlib gives %v", a, b, got, want)
			}
		})
	})

	for n := 1; n <= 40; n++ {
		base := []byte(strings.Repeat("aBcDeFgH", 6)[:n])
		for pos := range n {
			for v := range 256 {
				other := bytes.Clone(base)
				other[pos] = byte(v)
				want := asciiFoldLower(string(base)) == asciiFoldLower(string(other))
				if v < utf8.RuneSelf {
					if std := bytes.EqualFold(base, other); std != want {
						t.Fatalf("the stdlib disagrees with the ASCII fold on %q and %q", base, other)
					}
				}
				if got := EqualFold(base, other); got != want {
					t.Fatalf("EqualFold(%q, %q) = %v, want %v", base, other, got, want)
				}
				if got := EqualFold(string(other), string(base)); got != want {
					t.Fatalf("EqualFold(%q, %q) = %v, want %v", other, base, got, want)
				}
			}
		}
	}
}

// Test_EqualFold_OnlyASCIIFolds pins the deliberate difference: the stdlib's Unicode folding relates letters of different bytes.
func Test_EqualFold_OnlyASCIIFolds(t *testing.T) {
	t.Parallel()
	for _, tc := range [][2]string{
		{"k", "\u212a"}, // KELVIN SIGN folds to k
		{"s", "\u017f"}, // LONG S folds to s
		{"\u00e9", "\u00c9"},
		{"\u03c3", "\u03c2"},
		{"\u01c5", "\u01c6"},
	} {
		if !strings.EqualFold(tc[0], tc[1]) {
			t.Fatalf("the stdlib no longer folds %q and %q: the documented difference is gone", tc[0], tc[1])
		}
		if EqualFold(tc[0], tc[1]) || EqualFold([]byte(tc[0]), []byte(tc[1])) {
			t.Fatalf("EqualFold(%q, %q) = true, want ASCII-only folding", tc[0], tc[1])
		}
		if !EqualFold(tc[0], tc[0]) {
			t.Fatalf("EqualFold(%q, %q) = false for identical bytes", tc[0], tc[0])
		}
	}
}
