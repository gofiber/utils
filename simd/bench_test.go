package simd

import (
	"bytes"
	"strconv"
	"testing"
)

// The benchmarks pair each function with the closest idiomatic alternative
// (stdlib where one exists, a scalar loop otherwise) using the
// fiber-vs-default sub-run naming shared with the other packages. Inputs are
// worst-case misses so every variant scans the full haystack.

var benchSizes = []int{8, 32, 64, 512, 4096}

func benchName(n int) string {
	return strconv.Itoa(n) + "B"
}

// benchData is 'a'..'w' repeating: pure ASCII, no digits, no match for the
// searched needles below.
func benchData(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a' + byte(i%23)
	}
	return b
}

func Benchmark_Memchr2(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if Memchr2(data, 'z', 'y') != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if bytes.ContainsAny(data, "zy") {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_Memchr3(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if Memchr3(data, 'z', 'y', 'x') != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if bytes.ContainsAny(data, "zyx") {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_MemchrPair(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if MemchrPair(data, 'a', 'c', 1) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if refMemchrPair(data, 'a', 'c', 1) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_MemchrDigit(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if MemchrDigit(data) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchrDigitGeneric(data) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_MemchrNotWord(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if MemchrNotWord(data) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchrNotWordGeneric(data) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_Memmem(b *testing.B) {
	needle := []byte("q7z")
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if Memmem(data, needle) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if bytes.Contains(data, needle) {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

func Benchmark_IsASCII(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if !IsASCII(data) {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
		b.Run(benchName(n)+"/swar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if !isASCIIGeneric(data) {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
	}
}
