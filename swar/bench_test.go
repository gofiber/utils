package swar

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The benchmarks below compare the SWAR primitives against their stdlib
// counterparts, mirroring the swar-vs-stdlib benchmarks in the root package.
// The primitives are exported as building blocks, so each benchmark measures
// the canonical loop composed from them (documented in the examples), not a
// bare primitive call: that is the unit of work a stdlib function performs.
// The composed loops live below and are pinned to the stdlib results by
// Test_BenchLoops_MatchStdlib.

// swarIndexByte is ExampleZeroLanes' canonical first-match scan: broadcast
// once, scan full words, and finish 8+ byte inputs with one overlapping word.
func swarIndexByte(s string, c byte) int {
	n := len(s)
	if n < WordLen {
		for i := range n {
			if s[i] == c {
				return i
			}
		}
		return -1
	}
	needle := Broadcast(c)
	i := 0
	for ; i+WordLen <= n; i += WordLen {
		if m := ZeroLanes(Load8(s, i) ^ needle); m != 0 {
			return i + FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := ZeroLanes(Load8(s, n-WordLen) ^ needle); m != 0 {
		return n - WordLen + FirstLane(m)
	}
	return -1
}

// swarLastIndexByte is the reverse scan: MatchByteMask's exact per-lane
// result feeds LastLane, and the sub-word head falls back to scalar code.
func swarLastIndexByte(b []byte, c byte) int {
	i := len(b)
	for ; i >= WordLen; i -= WordLen {
		if m := MatchByteMask(Load8(b, i-WordLen), c); m != 0 {
			return i - WordLen + LastLane(m)
		}
	}
	for i--; i >= 0; i-- {
		if b[i] == c {
			return i
		}
	}
	return -1
}

// swarIndexDigit finds the first ASCII digit via MatchRangeMask, with a
// scalar tail for the trailing sub-word bytes.
func swarIndexDigit(s string) int {
	i := 0
	for ; i+WordLen <= len(s); i += WordLen {
		if m := MatchRangeMask(Load8(s, i), '0', '9'); m != 0 {
			return i + FirstLane(m)
		}
	}
	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return i
		}
	}
	return -1
}

// swarToLower lower-cases ASCII in place, one word at a time, with a scalar
// tail — the loop internal/caseconv builds from the same primitives.
func swarToLower(b []byte) {
	i := 0
	for ; i+WordLen <= len(b); i += WordLen {
		Store8(b, i, ToLowerWord(Load8(b, i)))
	}
	for ; i < len(b); i++ {
		if c := b[i]; c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
}

// Test_BenchLoops_MatchStdlib pins the composed benchmark loops to the
// stdlib functions they are measured against, across sub-word, exact-word,
// and multi-word lengths with the needle at every position and absent.
func Test_BenchLoops_MatchStdlib(t *testing.T) {
	t.Parallel()
	sizes := []int{0, 1, 7, 8, 9, 16, 31, 64, 512}
	for _, n := range sizes {
		src := make([]byte, n)
		for i := range src {
			src[i] = 'A' + byte(i%23) // uppercase, no digits, no needle bytes
		}
		for at := -1; at < n; at++ {
			b := append([]byte(nil), src...)
			if at >= 0 {
				b[at] = ','
			}
			s := string(b)
			require.Equal(t, strings.IndexByte(s, ','), swarIndexByte(s, ','), "IndexByte size %d at %d", n, at)
			require.Equal(t, bytes.LastIndexByte(b, ','), swarLastIndexByte(b, ','), "LastIndexByte size %d at %d", n, at)

			d := append([]byte(nil), src...)
			if at >= 0 {
				d[at] = '7'
			}
			require.Equal(t, bytes.IndexAny(d, "0123456789"), swarIndexDigit(string(d)), "IndexDigit size %d at %d", n, at)
		}
		lower := append([]byte(nil), src...)
		swarToLower(lower)
		require.Equal(t, string(bytes.ToLower(src)), string(lower), "ToLower size %d", n)
	}
}

var benchSizes = []int{7, 8, 16, 32, 64, 512}

func benchInput(size, needleAt int, needle byte) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = 'a' + byte(i%23)
	}
	if needleAt >= 0 && needleAt < size {
		b[needleAt] = needle
	}
	return b
}

func Benchmark_Load8(b *testing.B) {
	buf := benchInput(4096, -1, 0)
	b.Run("swar", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		var acc uint64
		for b.Loop() {
			for i := 0; i+WordLen <= len(buf); i += WordLen {
				acc += Load8(buf, i)
			}
		}
		_ = acc
	})
	b.Run("stdlib-binary", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		var acc uint64
		for b.Loop() {
			for i := 0; i+WordLen <= len(buf); i += WordLen {
				acc += binary.LittleEndian.Uint64(buf[i:])
			}
		}
		_ = acc
	})
}

func Benchmark_Store8(b *testing.B) {
	buf := make([]byte, 4096)
	w := Broadcast('x')
	b.Run("swar", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		for b.Loop() {
			for i := 0; i+WordLen <= len(buf); i += WordLen {
				Store8(buf, i, w)
			}
		}
	})
	b.Run("stdlib-binary", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		for b.Loop() {
			for i := 0; i+WordLen <= len(buf); i += WordLen {
				binary.LittleEndian.PutUint64(buf[i:], w)
			}
		}
	})
}

func Benchmark_IndexByte(b *testing.B) {
	for _, size := range benchSizes {
		s := string(benchInput(size, size-1, ',')) // worst case: match at the very end
		b.Run(strconv.Itoa(size)+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = swarIndexByte(s, ',')
			}
			_ = r
		})
		b.Run(strconv.Itoa(size)+"B/stdlib-strings", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = strings.IndexByte(s, ',')
			}
			_ = r
		})
	}
}

func Benchmark_LastIndexByte(b *testing.B) {
	for _, size := range benchSizes {
		buf := benchInput(size, 0, ',') // worst case for a reverse scan: match at the start
		b.Run(strconv.Itoa(size)+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = swarLastIndexByte(buf, ',')
			}
			_ = r
		})
		b.Run(strconv.Itoa(size)+"B/stdlib-bytes", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = bytes.LastIndexByte(buf, ',')
			}
			_ = r
		})
	}
}

func Benchmark_IndexDigit(b *testing.B) {
	for _, size := range benchSizes {
		s := string(benchInput(size, size-1, '7'))
		b.Run(strconv.Itoa(size)+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = swarIndexDigit(s)
			}
			_ = r
		})
		b.Run(strconv.Itoa(size)+"B/stdlib-indexany", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = strings.IndexAny(s, "0123456789")
			}
			_ = r
		})
	}
}

func Benchmark_ToLowerWord(b *testing.B) {
	for _, size := range benchSizes {
		template := bytes.ToUpper(benchInput(size, -1, 0))
		b.Run(strconv.Itoa(size)+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			work := make([]byte, len(template))
			for b.Loop() {
				copy(work, template)
				swarToLower(work)
			}
		})
		b.Run(strconv.Itoa(size)+"B/stdlib-bytes", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r []byte
			for b.Loop() {
				r = bytes.ToLower(template)
			}
			_ = r
		})
	}
}
