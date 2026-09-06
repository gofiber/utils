package utils

import (
	"bytes"
	"strings"
	"testing"

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
