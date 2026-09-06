package utils

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FormatUint(t *testing.T) {
	t.Parallel()
	tests := []uint64{
		0, 1, 9, 10, 11, 99, 100, 101, 999, 1000,
		12345, 123456789, 9999999999,
		math.MaxUint32, math.MaxUint64,
	}
	for _, tt := range tests {
		expected := strconv.FormatUint(tt, 10)
		result := FormatUint(tt)
		require.Equal(t, expected, result, "FormatUint(%d)", tt)
	}
}

func Test_FormatInt(t *testing.T) {
	t.Parallel()
	tests := []int64{
		0, 1, -1, 9, -9, 10, -10, 99, -99, 100, -100,
		12345, -12345, 123456789, -123456789,
		math.MaxInt32, math.MinInt32,
		math.MaxInt64, math.MinInt64,
	}
	for _, tt := range tests {
		expected := strconv.FormatInt(tt, 10)
		result := FormatInt(tt)
		require.Equal(t, expected, result, "FormatInt(%d)", tt)
	}
}

func Test_FormatUint32(t *testing.T) {
	t.Parallel()
	tests := []uint32{0, 1, 99, 100, 12345, math.MaxUint32}
	for _, tt := range tests {
		expected := strconv.FormatUint(uint64(tt), 10)
		result := FormatUint32(tt)
		require.Equal(t, expected, result, "FormatUint32(%d)", tt)
	}
}

func Test_FormatInt32(t *testing.T) {
	t.Parallel()
	tests := []int32{0, 1, -1, 99, -99, math.MaxInt32, math.MinInt32}
	for _, tt := range tests {
		expected := strconv.FormatInt(int64(tt), 10)
		result := FormatInt32(tt)
		require.Equal(t, expected, result, "FormatInt32(%d)", tt)
	}
}

func Test_FormatUint16(t *testing.T) {
	t.Parallel()
	tests := []uint16{0, 1, 99, 100, 12345, math.MaxUint16}
	for _, tt := range tests {
		expected := strconv.FormatUint(uint64(tt), 10)
		result := FormatUint16(tt)
		require.Equal(t, expected, result, "FormatUint16(%d)", tt)
	}
}

func Test_FormatInt16(t *testing.T) {
	t.Parallel()
	tests := []int16{0, 1, -1, 99, -99, math.MaxInt16, math.MinInt16}
	for _, tt := range tests {
		expected := strconv.FormatInt(int64(tt), 10)
		result := FormatInt16(tt)
		require.Equal(t, expected, result, "FormatInt16(%d)", tt)
	}
}

func Test_FormatUint8(t *testing.T) {
	t.Parallel()
	tests := []uint8{0, 1, 99, 100, math.MaxUint8}
	for _, tt := range tests {
		expected := strconv.FormatUint(uint64(tt), 10)
		result := FormatUint8(tt)
		require.Equal(t, expected, result, "FormatUint8(%d)", tt)
	}
}

func Test_FormatInt8(t *testing.T) {
	t.Parallel()
	tests := []int8{0, 1, -1, 99, -99, math.MaxInt8, math.MinInt8}
	for _, tt := range tests {
		expected := strconv.FormatInt(int64(tt), 10)
		result := FormatInt8(tt)
		require.Equal(t, expected, result, "FormatInt8(%d)", tt)
	}
}

func Test_AppendUint(t *testing.T) {
	t.Parallel()
	tests := []uint64{0, 1, 99, 100, 12345, 123456789, math.MaxUint64}
	for _, tt := range tests {
		expected := strconv.AppendUint([]byte("prefix"), tt, 10)
		result := AppendUint([]byte("prefix"), tt)
		require.Equal(t, expected, result, "AppendUint(%d)", tt)
	}
}

func Test_AppendInt(t *testing.T) {
	t.Parallel()
	tests := []int64{0, 1, -1, 99, -99, 100, -100, 12345, -12345, math.MaxInt64, math.MinInt64}
	for _, tt := range tests {
		expected := strconv.AppendInt([]byte("prefix"), tt, 10)
		result := AppendInt([]byte("prefix"), tt)
		require.Equal(t, expected, result, "AppendInt(%d)", tt)
	}
}

func Test_AppendInt_SmallNegativeCache(t *testing.T) {
	t.Parallel()
	// Test all small negative integers that should use the cache (-1 to -99)
	for i := int64(-1); i >= -99; i-- {
		t.Run(strconv.FormatInt(i, 10), func(t *testing.T) {
			t.Parallel()
			expected := strconv.AppendInt([]byte("prefix"), i, 10)
			result := AppendInt([]byte("prefix"), i)
			require.Equal(t, expected, result)
		})
	}
	// Verify boundary: -100 should NOT use the cache
	t.Run("boundary/-100", func(t *testing.T) {
		t.Parallel()
		expected := strconv.AppendInt([]byte("prefix"), -100, 10)
		result := AppendInt([]byte("prefix"), -100)
		require.Equal(t, expected, result)
	})
}

// Benchmarks

func Benchmark_FormatUint(b *testing.B) {
	inputs := []struct {
		name  string
		value uint64
	}{
		{"small", 42},
		{"medium", 123456789},
		{"large", math.MaxUint64},
	}

	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				FormatUint(input.value)
			}
		})
		b.Run(input.name+"/strconv", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				strconv.FormatUint(input.value, 10)
			}
		})
	}
}

func Benchmark_FormatInt(b *testing.B) {
	inputs := []struct {
		name  string
		value int64
	}{
		{"small_pos", 42},
		{"small_neg", -42},
		{"medium_pos", 123456789},
		{"medium_neg", -123456789},
		{"large_pos", math.MaxInt64},
		{"large_neg", math.MinInt64},
	}

	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				FormatInt(input.value)
			}
		})
		b.Run(input.name+"/strconv", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				strconv.FormatInt(input.value, 10)
			}
		})
	}
}

func Benchmark_FormatUint32(b *testing.B) {
	input := uint32(123456789)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatUint32(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt32(b *testing.B) {
	input := int32(-123456789)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatInt32(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_FormatUint16(b *testing.B) {
	input := uint16(12345)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatUint16(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt16(b *testing.B) {
	input := int16(-12345)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatInt16(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_FormatUint8(b *testing.B) {
	input := uint8(255)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatUint8(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt8(b *testing.B) {
	input := int8(-128)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FormatInt8(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_AppendUint(b *testing.B) {
	input := uint64(123456789)
	dst := make([]byte, 0, 32)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			AppendUint(dst, input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			strconv.AppendUint(dst, input, 10)
		}
	})
}

func Benchmark_AppendInt(b *testing.B) {
	inputs := []struct {
		name  string
		value int64
	}{
		{"small_neg", -42},
		{"medium_neg", -123456789},
	}
	dst := make([]byte, 0, 32)

	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendInt(dst, input.value)
			}
		})
		b.Run(input.name+"/strconv", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				strconv.AppendInt(dst, input.value, 10)
			}
		})
	}
}

func Test_UintDigits_IntDigits(t *testing.T) {
	t.Parallel()
	checkU := func(n uint64) {
		t.Helper()
		require.Equal(t, len(strconv.FormatUint(n, 10)), uintDigits(n), "n=%d", n)
	}
	for n := range uint64(100000) {
		checkU(n)
	}
	// Every power-of-ten boundary up to the uint64 limit.
	p := uint64(1)
	for range 20 {
		checkU(p - 1)
		checkU(p)
		checkU(p + 1)
		if p > math.MaxUint64/10 {
			break
		}
		p *= 10
	}
	checkU(math.MaxUint64)

	for n := int64(-100000); n < 100000; n++ {
		require.Equal(t, len(strconv.FormatInt(n, 10)), intDigits(n), "n=%d", n)
	}
	require.Equal(t, 20, intDigits(math.MinInt64))
	require.Equal(t, 19, intDigits(math.MaxInt64))
}

// Test_Format_Sweep pins the lane-group formatting to strconv across every
// value up to two lane groups' worth of shapes: a dense sweep of the first
// two million values, every decimal boundary (10^k-1, 10^k, 10^k+1) for
// every width, the group boundaries the 8-digit kernel splits at, the
// extremes of every width, and a deterministic pseudo-random spread over
// all bit lengths.
func Test_Format_Sweep(t *testing.T) {
	t.Parallel()
	check := func(n uint64) {
		t.Helper()
		want := strconv.FormatUint(n, 10)
		require.Equal(t, want, FormatUint(n), "FormatUint(%d)", n)
		require.Equal(t, want, string(AppendUint([]byte("x"), n))[1:], "AppendUint(%d)", n)
		if n <= math.MaxUint32 {
			require.Equal(t, want, FormatUint32(uint32(n)), "FormatUint32(%d)", n)
		}
		if n <= math.MaxUint16 {
			require.Equal(t, want, FormatUint16(uint16(n)), "FormatUint16(%d)", n)
		}
		if n <= math.MaxInt64 {
			i := int64(n)
			require.Equal(t, strconv.FormatInt(i, 10), FormatInt(i), "FormatInt(%d)", i)
			require.Equal(t, strconv.FormatInt(-i, 10), FormatInt(-i), "FormatInt(%d)", -i)
			require.Equal(t, strconv.FormatInt(-i, 10), string(AppendInt([]byte("x"), -i))[1:], "AppendInt(%d)", -i)
		}
		if n <= math.MaxInt32 {
			i := int32(n)
			require.Equal(t, strconv.FormatInt(int64(i), 10), FormatInt32(i), "FormatInt32(%d)", i)
			require.Equal(t, strconv.FormatInt(int64(-i), 10), FormatInt32(-i), "FormatInt32(%d)", -i)
		}
		if n <= math.MaxInt16 {
			i := int16(n)
			require.Equal(t, strconv.FormatInt(int64(i), 10), FormatInt16(i), "FormatInt16(%d)", i)
			require.Equal(t, strconv.FormatInt(int64(-i), 10), FormatInt16(-i), "FormatInt16(%d)", -i)
		}
	}

	for n := range uint64(300_000) {
		check(n)
	}
	for _, p := range pow10 {
		check(p - 1)
		check(p)
		check(p + 1)
	}
	for _, n := range []uint64{
		99999999, 100000000, 100000001, // one lane group
		9999999999999999, 10000000000000000, 10000000000000001, // two lane groups
		math.MaxUint16, math.MaxInt16, math.MaxUint32, math.MaxInt32,
		math.MaxInt64, math.MaxInt64 + 1, math.MaxUint64,
		1844674407370955161, 1844674407370955162, // largest top groups
	} {
		check(n)
	}
	require.Equal(t, "-9223372036854775808", FormatInt(math.MinInt64))
	require.Equal(t, "-9223372036854775808", string(AppendInt(nil, math.MinInt64)))
	require.Equal(t, "-2147483648", FormatInt32(math.MinInt32))
	require.Equal(t, "-32768", FormatInt16(math.MinInt16))

	// xorshift spread: a few values at every bit length.
	x := uint64(0x9E3779B97F4A7C15)
	for range 20000 {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		check(x)
		check(x >> (x % 64))
	}
}
