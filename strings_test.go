// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"strings"
	"testing"
)

func Test_Utils_ToUpper(t *testing.T) {
	t.Parallel()
	res := ToUpper("/my/name/is/:param/*")
	AssertEqual(t, "/MY/NAME/IS/:PARAM/*", res)
}

func Benchmark_ToUpper(b *testing.B) {
	var path = "/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts"
	var res string

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToUpper(path)
		}
		AssertEqual(b, "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS", res)
	})
	b.Run("IfToUpper-Upper", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToUpper(upperStr)
		}
		require.Equal(b, upperStr, res)
	})
	b.Run("IfToUpper-Mixed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToUpper(largeStr)
		}
		require.Equal(b, upperStr, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToUpper(path)
		}
		AssertEqual(b, "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS", res)
	})
}

func Test_Utils_ToLower(t *testing.T) {
	t.Parallel()
	res := ToLower("/MY/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my/name/is/:param/*", res)
	res = ToLower("/MY1/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my1/name/is/:param/*", res)
	res = ToLower("/MY2/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my2/name/is/:param/*", res)
	res = ToLower("/MY3/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my3/name/is/:param/*", res)
	res = ToLower("/MY4/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my4/name/is/:param/*", res)
}

func Benchmark_ToLower(b *testing.B) {
	var path = "/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts"
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLower(path)
		}
		AssertEqual(b, "/repos/gofiber/fiber/issues/187643/comments", res)
	})
	b.Run("IfToLower-Lower", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower(lowerStr)
		}
		require.Equal(b, lowerStr, res)
	})
	b.Run("IfToLower-Mixed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower(largeStr)
		}
		require.Equal(b, lowerStr, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower(path)
		}
		AssertEqual(b, "/repos/gofiber/fiber/issues/187643/comments", res)
	})
}

func Test_Utils_TrimRight(t *testing.T) {
	t.Parallel()
	res := TrimRight("/test//////", '/')
	AssertEqual(t, "/test", res)

	res = TrimRight("/test", '/')
	AssertEqual(t, "/test", res)
}
func Benchmark_TrimRight(b *testing.B) {
	var res string

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimRight("foobar  ", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimRight("foobar  ", " ")
		}
		AssertEqual(b, "foobar", res)
	})
}

func Test_Utils_TrimLeft(t *testing.T) {
	t.Parallel()
	res := TrimLeft("////test/", '/')
	AssertEqual(t, "test/", res)

	res = TrimLeft("test/", '/')
	AssertEqual(t, "test/", res)
}
func Benchmark_TrimLeft(b *testing.B) {
	var res string

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimLeft("  foobar", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimLeft("  foobar", " ")
		}
		AssertEqual(b, "foobar", res)
	})
}
func Test_Utils_Trim(t *testing.T) {
	t.Parallel()
	res := Trim("   test  ", ' ')
	AssertEqual(t, "test", res)

	res = Trim("test", ' ')
	AssertEqual(t, "test", res)

	res = Trim(".test", '.')
	AssertEqual(t, "test", res)
}

func Benchmark_Trim(b *testing.B) {
	var res string

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = Trim("  foobar   ", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.Trim("  foobar   ", " ")
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("default.trimspace", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimSpace("  foobar   ")
		}
		AssertEqual(b, "foobar", res)
	})
}

func Test_IfToUpper(t *testing.T) {
	t.Parallel()
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("MYNAMEISPARAM")) // already uppercase
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("mynameisparam")) // lowercase to uppercase
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("MyNameIsParam")) // mixed case
}

func Test_IfToLower(t *testing.T) {
	t.Parallel()
	require.Equal(t, "mynameisparam", IfToLower("mynameisparam"))           // already lowercase
	require.Equal(t, "mynameisparam", IfToLower("myNameIsParam"))           // mixed case
	require.Equal(t, "https://gofiber.io", IfToLower("https://gofiber.io")) // Origin Header Type URL
	require.Equal(t, "mynameisparam", IfToLower("MYNAMEISPARAM"))           // uppercase to lowercase
}

// Benchmark_IfToLower_HeadersOrigin benchmarks the IfToLower function with an origin header type URL.
// These headers are typically lowercase, so the function should return the input string without modification.
func Benchmark_IfToToLower_HeadersOrigin(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLower("https://gofiber.io")
		}
		require.Equal(b, "https://gofiber.io", res)
	})
	b.Run("IfToLower-Lower", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower("https://gofiber.io")
		}
		require.Equal(b, "https://gofiber.io", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower("https://gofiber.io")
		}
		require.Equal(b, "https://gofiber.io", res)
	})
}
