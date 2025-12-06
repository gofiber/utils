package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
			// bytes.Equal treats nil and empty slices as equal
			require.True(t, bytes.Equal(stdResultBytes, resultBytes), "TrimSpace should match bytes.TrimSpace for %s", tc.name)
		})
	}
}

func Benchmark_TrimSpace(b *testing.B) {
	for _, tc := range trimSpaceBenchmarkCases {
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
	for _, tc := range trimSpaceBenchmarkCases {
		tc := tc // capture range variable
		input := []byte(tc.input)
		b.Run("fiber/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(input)))
			b.ResetTimer()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = TrimSpace(input)
			}
			_ = res
		})
		b.Run("default/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(input)))
			b.ResetTimer()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = bytes.TrimSpace(input)
			}
			_ = res
		})
	}
}

// Shared test cases for AddTrailingSlash benchmarks
var addTrailingSlashBenchmarkCases = []struct {
	name  string
	input string
}{
	{name: "empty", input: ""},
	{name: "slash-only", input: "/"},
	{name: "short-no-slash", input: "abc"},
	{name: "short-with-slash", input: "abc/"},
	{name: "path-no-slash", input: "/api/v1/users"},
	{name: "path-with-slash", input: "/api/v1/users/"},
	{name: "long-no-slash", input: "/api/v1/users/profile/settings/notifications"},
	{name: "long-with-slash", input: "/api/v1/users/profile/settings/notifications/"},
}

func Test_AddTrailingSlash(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "/"},
		{name: "slash-only", in: "/", want: "/"},
		{name: "short-no-slash", in: "abc", want: "abc/"},
		{name: "short-with-slash", in: "abc/", want: "abc/"},
		{name: "path-no-slash", in: "/api/v1/users", want: "/api/v1/users/"},
		{name: "path-with-slash", in: "/api/v1/users/", want: "/api/v1/users/"},
		{name: "double-slash", in: "abc//", want: "abc//"},
		{name: "root-path", in: "/", want: "/"},
		{name: "spaces", in: "  ", want: "  /"},
		{name: "unicode", in: "/日本語", want: "/日本語/"},
		{name: "unicode-with-slash", in: "/日本語/", want: "/日本語/"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name+"/string", func(t *testing.T) {
			t.Parallel()
			result := AddTrailingSlash(tc.in)
			require.Equal(t, tc.want, result)
		})
		t.Run(tc.name+"/bytes", func(t *testing.T) {
			t.Parallel()
			result := AddTrailingSlash([]byte(tc.in))
			require.Equal(t, []byte(tc.want), result)
		})
	}
}

func Test_AddTrailingSlash_NoMutation(t *testing.T) {
	t.Parallel()

	// Ensure original byte slice is not mutated
	original := []byte("test")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	_ = AddTrailingSlash(original)

	require.Equal(t, originalCopy, original, "original slice should not be mutated")
}

func Test_AddTrailingSlash_AlreadyHasSlash_ReturnsSame(t *testing.T) {
	t.Parallel()

	// For byte slices with trailing slash, should return the same slice instance
	input := []byte("test/")
	result := AddTrailingSlash(input)
	require.Equal(t, input, result)
	require.Same(t, &input[0], &result[0], "should return same slice instance")

	// For strings with trailing slash, should return the same string
	inputStr := "test/"
	resultStr := AddTrailingSlash(inputStr)
	require.Equal(t, inputStr, resultStr)
}

func Test_AddTrailingSlash_NamedTypes(t *testing.T) {
	t.Parallel()

	type MyString string
	type MyBytes []byte

	t.Run("named string without slash", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyString("abc"))
		require.Equal(t, MyString("abc/"), result)
	})

	t.Run("named string with slash", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyString("abc/"))
		require.Equal(t, MyString("abc/"), result)
	})

	t.Run("named bytes without slash", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyBytes("abc"))
		require.Equal(t, MyBytes("abc/"), result)
	})

	t.Run("named bytes with slash", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyBytes("abc/"))
		require.Equal(t, MyBytes("abc/"), result)
	})

	t.Run("empty named string", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyString(""))
		require.Equal(t, MyString("/"), result)
	})

	t.Run("empty named bytes", func(t *testing.T) {
		t.Parallel()
		result := AddTrailingSlash(MyBytes(""))
		require.Equal(t, MyBytes("/"), result)
	})
}

func Benchmark_AddTrailingSlash(b *testing.B) {
	for _, tc := range addTrailingSlashBenchmarkCases {
		tc := tc
		b.Run("string/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = AddTrailingSlash(tc.input)
			}
			_ = res
		})
	}
}

func Benchmark_AddTrailingSlashBytes(b *testing.B) {
	for _, tc := range addTrailingSlashBenchmarkCases {
		tc := tc
		input := []byte(tc.input)
		b.Run("bytes/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(input)))
			b.ResetTimer()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = AddTrailingSlash(input)
			}
			_ = res
		})
	}
}
