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
	tests := []int64{0, 1, -1, 99, -99, 12345, -12345, math.MaxInt64, math.MinInt64}
	for _, tt := range tests {
		expected := strconv.AppendInt([]byte("prefix"), tt, 10)
		result := AppendInt([]byte("prefix"), tt)
		require.Equal(t, expected, result, "AppendInt(%d)", tt)
	}
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
			for i := 0; i < b.N; i++ {
				_ = FormatUint(input.value)
			}
		})
		b.Run(input.name+"/strconv", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = strconv.FormatUint(input.value, 10)
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
			for i := 0; i < b.N; i++ {
				_ = FormatInt(input.value)
			}
		})
		b.Run(input.name+"/strconv", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = strconv.FormatInt(input.value, 10)
			}
		})
	}
}

func Benchmark_FormatUint32(b *testing.B) {
	input := uint32(123456789)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatUint32(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt32(b *testing.B) {
	input := int32(-123456789)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatInt32(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_FormatUint16(b *testing.B) {
	input := uint16(12345)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatUint16(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt16(b *testing.B) {
	input := int16(-12345)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatInt16(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_FormatUint8(b *testing.B) {
	input := uint8(255)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatUint8(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatUint(uint64(input), 10)
		}
	})
}

func Benchmark_FormatInt8(b *testing.B) {
	input := int8(-128)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatInt8(input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatInt(int64(input), 10)
		}
	})
}

func Benchmark_AppendUint(b *testing.B) {
	input := uint64(123456789)
	dst := make([]byte, 0, 32)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = AppendUint(dst, input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.AppendUint(dst, input, 10)
		}
	})
}

func Benchmark_AppendInt(b *testing.B) {
	input := int64(-123456789)
	dst := make([]byte, 0, 32)

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = AppendInt(dst, input)
		}
	})
	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.AppendInt(dst, input, 10)
		}
	})
}
