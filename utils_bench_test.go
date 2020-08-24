// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"testing"
)

// go test -v -run=^$ -bench=Benchmark_UUID -benchmem -count=2

func Benchmark_UUID(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = UUID()
		}
		AssertEqual(b, 36, len(res))
	})
	b.Run("default", func(b *testing.B) {
		rnd := make([]byte, 16)
		_, _ = rand.Read(rnd)
		for n := 0; n < b.N; n++ {
			res = fmt.Sprintf("%x-%x-%x-%x-%x", rnd[0:4], rnd[4:6], rnd[6:8], rnd[8:10], rnd[10:])
		}
		AssertEqual(b, 36, len(res))
	})
}

// go test -v -run=^$ -bench=GetString -benchmem -count=2

func Benchmark_GetString(b *testing.B) {
	var hello = []byte("Hello, World!")
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = GetString(hello)
		}
		AssertEqual(b, "Hello, World!", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = string(hello)
		}
		AssertEqual(b, "Hello, World!", res)
	})
}

// go test -v -cpu 2 -run=^$ -bench=GetBytes -benchmem -count=4

func Benchmark_GetBytes(b *testing.B) {
	var hello = "Hello, World!"
	var res []byte
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = GetBytes(hello)
		}
		AssertEqual(b, []byte("Hello, World!"), res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = []byte(hello)
		}
		AssertEqual(b, []byte("Hello, World!"), res)
	})
}

// go test -v -run=^$ -bench=Benchmark_GetMIME -benchmem -count=2
func Benchmark_GetMIME(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = GetMIME(".xml")
			res = GetMIME(".txt")
			res = GetMIME(".png")
			res = GetMIME(".exe")
			res = GetMIME(".json")
		}
		AssertEqual(b, "application/json", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = mime.TypeByExtension(".xml")
			res = mime.TypeByExtension(".txt")
			res = mime.TypeByExtension(".png")
			res = mime.TypeByExtension(".exe")
			res = mime.TypeByExtension(".json")
		}
		AssertEqual(b, "application/json", res)
	})
}

// go test -v -run=^$ -bench=Benchmark_StatusMessage -benchmem -count=4
func Benchmark_StatusMessage(b *testing.B) {
	var res string
	for n := 0; n < b.N; n++ {
		res = StatusMessage(http.StatusNotExtended)
	}
	AssertEqual(b, "Not Extended", res)
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
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower(path)
		}
		AssertEqual(b, "/repos/gofiber/fiber/issues/187643/comments", res)
	})
}

func Benchmark_ToLowerBytes(b *testing.B) {
	var path = []byte("/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts")
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLowerBytes(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/repos/gofiber/fiber/issues/187643/comments"), res), true)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.ToLower(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/repos/gofiber/fiber/issues/187643/comments"), res), true)
	})
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
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToUpper(path)
		}
		AssertEqual(b, "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS", res)
	})
}

func Benchmark_ToUpperBytes(b *testing.B) {
	var path = []byte("/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts")
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToUpperBytes(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS"), res), true)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.ToUpper(path)
		}
		AssertEqual(b, bytes.EqualFold(GetBytes("/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS"), res), true)
	})
}

func Benchmark_EqualFolds(b *testing.B) {
	var left = []byte("/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts")
	var right = []byte("/RePos/goFiber/Fiber/issues/187643/COMMENTS")
	var res bool

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = EqualsFold(left, right)
		}
		AssertEqual(b, true, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.EqualFold(left, right)
		}
		AssertEqual(b, true, res)
	})
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

func Benchmark_TrimBytes(b *testing.B) {
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimBytes([]byte("  foobar   "), ' ')
		}
		AssertEqual(b, []byte("foobar"), res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.Trim([]byte("  foobar   "), " ")
		}
		AssertEqual(b, []byte("foobar"), res)
	})
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

func Benchmark_TrimLeftBytes(b *testing.B) {
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimLeftBytes([]byte("  foobar"), ' ')
		}
		AssertEqual(b, []byte("foobar"), res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.TrimLeft([]byte("  foobar"), " ")
		}
		AssertEqual(b, []byte("foobar"), res)
	})
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

func Benchmark_TrimRightBytes(b *testing.B) {
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = TrimRightBytes([]byte("foobar  "), ' ')
		}
		AssertEqual(b, []byte("foobar"), res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.TrimRight([]byte("foobar  "), " ")
		}
		AssertEqual(b, []byte("foobar"), res)
	})
}
