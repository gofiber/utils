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

// memmemBenchSizes brackets the 128-byte crossover that memmemMinHaystack
// encodes. Benchmark_Memmem measures the exported routing end to end — its
// sub-128B simd rows run bytes.Index plus dispatch, so they show routing
// overhead only — while Benchmark_Memmem_Prefilter forces the prefilter
// helpers at every size, making the crossover itself reproducible from
// this file.
var memmemBenchSizes = []int{8, 32, 64, 96, 128, 192, 512, 4096}

func Benchmark_Memmem(b *testing.B) {
	needle := []byte("q7z")
	for _, n := range memmemBenchSizes {
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

// Benchmark_Memmem_Prefilter drives the prefilter helpers directly,
// bypassing Memmem's 128-byte routing, so the memmemMinHaystack tradeoff
// is measurable instead of asserted. Each helper gets the needle shape
// Memmem actually routes to it: the paired scan (short needle, two
// distinct rare bytes) breaks even with bytes.Index around the 64B point
// and wins from 96B up, so 128 carries a conservative margin that also
// covers the per-call SelectRareBytes cost these direct calls exclude.
// The single scan (long needle) anchors on 'q', which benchData's
// repeating 'a'..'w' cycle yields every 23 bytes, so from a few hundred
// bytes up it burns its miss budget and lands on the fallback, tracking
// bytes.Index within a constant — the documented bounded worst case; its
// wins come on data where the anchor byte is genuinely rare.
func Benchmark_Memmem_Prefilter(b *testing.B) {
	pairNeedle := []byte("q7z")
	pairInfo := SelectRareBytes(pairNeedle)
	singleNeedle := []byte("the quick brown fox")
	singleInfo := SelectRareBytes(singleNeedle)
	for _, n := range memmemBenchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/paired", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memmemPaired(data, pairNeedle, pairInfo) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/paired-stdlib", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if bytes.Contains(data, pairNeedle) {
					b.Fatal("unexpected match")
				}
			}
		})
		if n < len(singleNeedle) {
			continue
		}
		b.Run(benchName(n)+"/single", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memmemSingle(data, singleNeedle, singleInfo.Byte1, singleInfo.Index1) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/single-stdlib", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if bytes.Contains(data, singleNeedle) {
					b.Fatal("unexpected match")
				}
			}
		})
	}
}

// Benchmark_Memmem_Adversarial pins the miss-budget fallback: every
// haystack position is a rare-byte anchor and the needle mismatches only at
// its last byte, the shape that used to degrade to O(n*m). With the budget
// it must track bytes.Index within a constant.
func Benchmark_Memmem_Adversarial(b *testing.B) {
	// The mismatch byte must rank as more common than 'z' so the rare-byte
	// selection anchors on 'z' and every haystack position is a candidate.
	haystack := bytes.Repeat([]byte{'z'}, 1<<20)
	needle := append(bytes.Repeat([]byte{'z'}, 999), 'e')
	b.Run("simd", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			if Memmem(haystack, needle) != -1 {
				b.Fatal("unexpected match")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			if bytes.Contains(haystack, needle) {
				b.Fatal("unexpected match")
			}
		}
	})
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
