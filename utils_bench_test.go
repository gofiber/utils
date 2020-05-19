// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"bytes"
	"strings"
	"testing"
)

// go test -v -run=^$ -bench=Benchmark_ToLower -benchmem -count=4
func Benchmark_ToLower(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg"
	var res string
	for n := 0; n < b.N; n++ {
		res = ToLower(str)
	}
	if res != "long_demo_string" {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_ToLower_Default(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg"
	var res string
	for n := 0; n < b.N; n++ {
		res = strings.ToLower(str)
	}
	if res != "long_demo_string" {
		b.Fatalf("Invalid result: %s", res)
	}
}

// go test -v -run=^$ -bench=Benchmark_ToLowerBytes -benchmem -count=4
func Benchmark_ToLowerBytes(b *testing.B) {
	var str = []byte("LoNg_DeMo_StRiNg")
	var res []byte
	for n := 0; n < b.N; n++ {
		res = ToLowerBytes(str)
	}
	if !bytes.Equal([]byte("long_demo_string"), res) {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_ToLowerBytes_Default(b *testing.B) {
	var str = []byte("LoNg_DeMo_StRiNg")
	var res []byte
	for n := 0; n < b.N; n++ {
		res = bytes.ToLower(str)
	}
	if !bytes.Equal([]byte("long_demo_string"), res) {
		b.Fatalf("Invalid result: %s", res)
	}
}

// go test -v -run=^$ -bench=Benchmark_ToUpper -benchmem -count=4
func Benchmark_ToUpper(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg"
	var res string
	for n := 0; n < b.N; n++ {
		res = ToUpper(str)
	}
	if res != "LONG_DEMO_STRING" {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_ToUpper_Default(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg"
	var res string
	for n := 0; n < b.N; n++ {
		res = strings.ToUpper(str)
	}
	if res != "LONG_DEMO_STRING" {
		b.Fatalf("Invalid result: %s", res)
	}
}

// go test -v -run=^$ -bench=Benchmark_ToUpperBytes -benchmem -count=4
func Benchmark_ToUpperBytes(b *testing.B) {
	var str = []byte("LoNg_DeMo_StRiNg")
	var res []byte
	for n := 0; n < b.N; n++ {
		res = ToUpperBytes(str)
	}
	if !bytes.Equal([]byte("LONG_DEMO_STRING"), res) {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_ToUpperBytes_Default(b *testing.B) {
	var str = []byte("LoNg_DeMo_StRiNg")
	var res []byte
	for n := 0; n < b.N; n++ {
		res = bytes.ToUpper(str)
	}
	if !bytes.Equal([]byte("LONG_DEMO_STRING"), res) {
		b.Fatalf("Invalid result: %s", res)
	}
}

// go test -v -run=^$ -bench=Benchmark_TrimRight -benchmem -count=4
func Benchmark_TrimRight(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg/////"
	var res string
	for n := 0; n < b.N; n++ {
		res = TrimRight(str, '/')
	}
	if res != "LoNg_DeMo_StRiNg" {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_TrimRight_Default(b *testing.B) {
	var str = "LoNg_DeMo_StRiNg/////"
	var res string
	for n := 0; n < b.N; n++ {
		res = strings.TrimRight(str, "/")
	}
	if res != "LoNg_DeMo_StRiNg" {
		b.Fatalf("Invalid result: %s", res)
	}
}

// go test -v -run=^$ -bench=Benchmark_Trim -benchmem -count=4
func Benchmark_Trim(b *testing.B) {
	var str = "/////LoNg_DeMo_StRiNg/////"
	var res string
	for n := 0; n < b.N; n++ {
		res = Trim(str, '/')
	}
	if res != "LoNg_DeMo_StRiNg" {
		b.Fatalf("Invalid result: %s", res)
	}
}
func Benchmark_Trim_Default(b *testing.B) {
	var str = "/////LoNg_DeMo_StRiNg/////"
	var res string
	for n := 0; n < b.N; n++ {
		res = strings.Trim(str, "/")
	}
	if res != "LoNg_DeMo_StRiNg" {
		b.Fatalf("Invalid result: %s", res)
	}
}
