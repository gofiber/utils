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
		name     string
		input    string
		expected string
	}{
		// Edge cases
		{name: "empty", input: "", expected: ""},
		{name: "no-whitespace", input: "hello", expected: "hello"},
		{name: "all-spaces", input: "     ", expected: ""},
		{name: "all-tabs", input: "\t\t\t", expected: ""},
		{name: "all-newlines", input: "\n\n\n", expected: ""},
		{name: "mixed-whitespace-only", input: " \t\n\r\v\f ", expected: ""},

		// Leading whitespace
		{name: "leading-space", input: " hello", expected: "hello"},
		{name: "leading-spaces", input: "   hello", expected: "hello"},
		{name: "leading-tab", input: "\thello", expected: "hello"},
		{name: "leading-newline", input: "\nhello", expected: "hello"},
		{name: "leading-mixed", input: " \t\n hello", expected: "hello"},

		// Trailing whitespace
		{name: "trailing-space", input: "hello ", expected: "hello"},
		{name: "trailing-spaces", input: "hello   ", expected: "hello"},
		{name: "trailing-tab", input: "hello\t", expected: "hello"},
		{name: "trailing-newline", input: "hello\n", expected: "hello"},
		{name: "trailing-mixed", input: "hello \t\n ", expected: "hello"},

		// Both leading and trailing
		{name: "both-space", input: " hello ", expected: "hello"},
		{name: "both-mixed", input: " \t\nhello\n\t ", expected: "hello"},
		{name: "both-many", input: "    hello world    ", expected: "hello world"},

		// Whitespace in the middle (should be preserved)
		{name: "middle-space", input: "hello world", expected: "hello world"},
		{name: "middle-tab", input: "hello\tworld", expected: "hello\tworld"},
		{name: "middle-newline", input: "hello\nworld", expected: "hello\nworld"},
		{name: "middle-and-edges", input: "  hello\t\nworld  ", expected: "hello\t\nworld"},

		// Single character
		{name: "single-char", input: "a", expected: "a"},
		{name: "single-space", input: " ", expected: ""},

		// All ASCII whitespace characters
		{name: "all-whitespace-types", input: " \t\n\r\v\fhello\f\v\r\n\t ", expected: "hello"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test string variant
			result := TrimSpace(tc.input)
			require.Equal(t, tc.expected, result, "TrimSpace failed for string %s", tc.name)

			// Verify it matches strings.TrimSpace behavior
			stdResult := strings.TrimSpace(tc.input)
			require.Equal(t, stdResult, result, "TrimSpace should match strings.TrimSpace for %s", tc.name)

			// Test []byte variant
			resultBytes := TrimSpace([]byte(tc.input))
			require.Equal(t, []byte(tc.expected), resultBytes, "TrimSpace failed for bytes %s", tc.name)

			// Verify it matches bytes.TrimSpace behavior (nil and empty slices are treated as equal)
			stdResultBytes := bytes.TrimSpace([]byte(tc.input))
			// Only compare if not both empty (nil vs [] are both valid empty slices)
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
		{name: "no-trim", input: "hello"},
		{name: "no-trim-long", input: "hello world this is a longer string"},
		{name: "leading-space", input: " hello"},
		{name: "trailing-space", input: "hello "},
		{name: "both-spaces", input: " hello "},
		{name: "both-many-spaces", input: "    hello world    "},
		{name: "leading-mixed", input: " \t\n hello"},
		{name: "trailing-mixed", input: "hello \t\n "},
		{name: "both-mixed", input: " \t\nhello world\n\t "},
		{name: "all-spaces", input: "          "},
		{name: "medium-no-trim", input: strings.Repeat("a", 64)},
		{name: "medium-with-trim", input: "  " + strings.Repeat("a", 64) + "  "},
		{name: "large-no-trim", input: strings.Repeat("hello", 50)},
		{name: "large-with-trim", input: "    " + strings.Repeat("hello", 50) + "    "},
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
		{name: "no-trim", input: []byte("hello")},
		{name: "no-trim-long", input: []byte("hello world this is a longer string")},
		{name: "leading-space", input: []byte(" hello")},
		{name: "trailing-space", input: []byte("hello ")},
		{name: "both-spaces", input: []byte(" hello ")},
		{name: "both-many-spaces", input: []byte("    hello world    ")},
		{name: "leading-mixed", input: []byte(" \t\n hello")},
		{name: "trailing-mixed", input: []byte("hello \t\n ")},
		{name: "both-mixed", input: []byte(" \t\nhello world\n\t ")},
		{name: "all-spaces", input: []byte("          ")},
		{name: "medium-no-trim", input: bytes.Repeat([]byte("a"), 64)},
		{name: "medium-with-trim", input: append(append([]byte("  "), bytes.Repeat([]byte("a"), 64)...), []byte("  ")...)},
		{name: "large-no-trim", input: bytes.Repeat([]byte("hello"), 50)},
		{name: "large-with-trim", input: append(append([]byte("    "), bytes.Repeat([]byte("hello"), 50)...), []byte("    ")...)},
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
