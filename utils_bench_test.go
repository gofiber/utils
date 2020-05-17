// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"bytes"
	"strings"
	"testing"
)

func Benchmark_Utils_getMIME(b *testing.B) {
	var res string
	for n := 0; n < b.N; n++ {
		res = GetMIME(".json")
		res = GetMIME(".xml")
		res = GetMIME("xml")
		res = GetMIME("json")
	}
	AssertEqual(b, "application/json", res)
}

func Benchmark_Utils_ToLower(b *testing.B) {
	var path = "/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts"
	var res string

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLower(path)
		}
		AssertEqual(b, "/repos/gofiber/fiber/issues/187643/comments", res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower(path)
		}
		AssertEqual(b, "/repos/gofiber/fiber/issues/187643/comments", res)
	})
}

func Benchmark_Utils_ToLowerBytes(b *testing.B) {
	var path = []byte("/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts")
	var res []byte

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLowerBytes(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/repos/gofiber/fiber/issues/187643/comments"), res), true)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.ToLower(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/repos/gofiber/fiber/issues/187643/comments"), res), true)
	})
}

func Benchmark_Utils_ToUpper(b *testing.B) {
	var path = "/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts"
	var res string

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToUpper(path)
		}
		AssertEqual(b, "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS", res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToUpper(path)
		}
		AssertEqual(b, "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS", res)
	})
}

func Benchmark_Utils_EqualFolds(b *testing.B) {
	var left = []byte("/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts")
	var right = []byte("/RePos/goFiber/Fiber/issues/187643/COMMENTS")
	var res bool

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = EqualsFold(left, right)
		}
		AssertEqual(b, true, res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.EqualFold(left, right)
		}
		AssertEqual(b, true, res)
	})
}
func Benchmark_Utils_Trim(b *testing.B) {
	var res string

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = Trim("  foobar   ", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.Trim("  foobar   ", " ")
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("original 2", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimSpace("  foobar   ")
		}
		AssertEqual(b, "foobar", res)
	})
}
func Benchmark_Utils_TrimLeft(b *testing.B) {
	var res string

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimLeft("  foobar", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimLeft("  foobar", " ")
		}
		AssertEqual(b, "foobar", res)
	})
}
func Benchmark_Utils_TrimRight(b *testing.B) {
	var res string

	b.Run("optimized", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimRight("foobar  ", ' ')
		}
		AssertEqual(b, "foobar", res)
	})
	b.Run("original", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.TrimRight("foobar  ", " ")
		}
		AssertEqual(b, "foobar", res)
	})
}
