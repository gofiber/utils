package simd

import (
	"bytes"
	"strconv"
	"testing"
)

// Each function is paired with its closest alternative, on worst-case misses.

var benchSizes = []int{8, 32, 64, 512, 4096}

func benchName(n int) string {
	return strconv.Itoa(n) + "B"
}

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
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchr2Generic(data, 'z', 'y') != -1 {
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
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchr3Generic(data, 'z', 'y', 'x') != -1 {
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
		b.Run(benchName(n)+"/swar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchrPairGeneric(data, 'a', 'c', 1) != -1 {
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

// benchNonWordData holds no \w byte (benchData starts with one), so MemchrWord scans it all.
func benchNonWordData(n int) []byte {
	const punct = "!\"#$%&'()*+,-./:;<=>?@[\\]^`{|}~ "
	b := make([]byte, n)
	for i := range b {
		b[i] = punct[i%len(punct)]
	}
	return b
}

func Benchmark_MemchrWord(b *testing.B) {
	for _, n := range benchSizes {
		data := benchNonWordData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if MemchrWord(data) != -1 {
					b.Fatal("unexpected match")
				}
			}
		})
		b.Run(benchName(n)+"/scalar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if memchrWordGeneric(data) != -1 {
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

// memmemBenchSizes brackets the routing points (65, 128); Benchmark_Memmem's simd rows up to
// 64B show only routing overhead.
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

// Benchmark_Memmem_Prefilter runs the prefilter helpers directly to measure the memmemMin*Haystack
// crossovers. The paired scan wins from 65B, where bytes.Index leaves its brute-force kernel. The
// single scan anchors on 'q', which benchData yields every 23 bytes: from a few hundred bytes it
// exhausts its miss budget and tracks bytes.Index, the bounded worst case.
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

// Benchmark_Memmem_Adversarial pins the miss budget: every position is an anchor and the needle
// mismatches only at its end, which used to be O(n*m).
func Benchmark_Memmem_Adversarial(b *testing.B) {
	// The mismatch byte must be more common than 'z', so 'z' anchors.
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

func Benchmark_FirstNonASCII(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if FirstNonASCII(data) != -1 {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
		b.Run(benchName(n)+"/swar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if firstNonASCIIGeneric(data) != -1 {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
	}
}

func Benchmark_CountNonASCII(b *testing.B) {
	for _, n := range benchSizes {
		data := benchData(n)
		b.Run(benchName(n)+"/simd", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if CountNonASCII(data) != 0 {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
		b.Run(benchName(n)+"/swar", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if countNonASCIIGeneric(data) != 0 {
					b.Fatal("unexpected non-ASCII")
				}
			}
		})
	}
}
