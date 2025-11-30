package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		for n := 0; n < b.N; n++ {
			res = EqualFold(left, right)
		}
		require.True(b, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = EqualFold(upperStr, lowerStr)
		}
		require.True(b, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = TrimRight(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = TrimRight(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = TrimLeft(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = TrimLeft(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = Trim(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			res = strings.Trim(word, " ")
		}
		require.Equal(b, expected, res)
	})
	b.Run("default.trimspace", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		for n := 0; n < b.N; n++ {
			res = Trim(word, ' ')
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			res = bytes.Trim(word, " ")
		}
		require.Equal(b, expected, res)
	})
	b.Run("default.trimspace", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test string variant - must match strings.TrimSpace
			result := TrimSpace(tc.input)
			stdResult := strings.TrimSpace(tc.input)
			require.Equal(t, stdResult, result, "TrimSpace should match strings.TrimSpace for %s", tc.name)

			// Test []byte variant - must match bytes.TrimSpace
			resultBytes := TrimSpace([]byte(tc.input))
			stdResultBytes := bytes.TrimSpace([]byte(tc.input))
			// Both empty slices (nil vs []) are valid and equivalent
			if len(stdResultBytes) != 0 || len(resultBytes) != 0 {
				require.Equal(t, stdResultBytes, resultBytes, "TrimSpace should match bytes.TrimSpace for %s", tc.name)
			}
		})
	}
}

func Benchmark_TrimSpace(b *testing.B) {
	testCases := []struct {
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

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run("fiber/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = TrimSpace(tc.input)
			}
			_ = res
		})
		b.Run("default/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = strings.TrimSpace(tc.input)
			}
			_ = res
		})
	}
}

func Benchmark_TrimSpaceBytes(b *testing.B) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{name: "empty", input: []byte("")},
		{name: "spaces", input: []byte("   ")},
		{name: "ascii-word", input: []byte("  fiber  ")},
		{name: "auth-header-bearer", input: []byte("   Bearer abc.def.ghi   ")},
		{name: "auth-header-basic", input: []byte("\tBasic QWxhZGRpbjpvcGVuIHNlc2FtZQ==   ")},
		{name: "accept-encoding", input: []byte("  gzip, deflate, br  ")},
		{name: "content-type", input: []byte("  application/json  ")},
		{name: "x-forwarded-for", input: []byte(" 203.0.113.5, 198.51.100.7  ")},
		{name: "query-params", input: []byte("  user=alice&id=42  ")},
		{name: "ascii-long", input: []byte("   " + strings.Repeat("fiber utils ", 8) + "   ")},
		{name: "no-trim", input: []byte("fiber")},
		{name: "mixed-whitespace", input: []byte("\n\t fiber utils \r\n")},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run("fiber/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = TrimSpace(tc.input)
			}
			_ = res
		})
		b.Run("default/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = bytes.TrimSpace(tc.input)
			}
			_ = res
		})
	}
}
