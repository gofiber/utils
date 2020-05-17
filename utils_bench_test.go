// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import "testing"

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
