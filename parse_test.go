package utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	randv2 "math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/utils/v2/swar"
	"github.com/stretchr/testify/require"
)

func Test_ParseUint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     uint64
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"123456789", 123456789, true},
		{strconv.FormatUint(math.MaxUint64, 10), math.MaxUint64, true},
		{"18446744073709551616", 0, false},
		{"-1", 0, false},
		{"+123", 0, false},
		{"", 0, false},
		{"12a", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseUint(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseUint([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Test_ParseUint_Whitespace(t *testing.T) {
	t.Parallel()

	v, err := ParseUint(" 123")
	require.Error(t, err)
	require.Equal(t, uint64(0), v)
}

func Test_ParseNativeUint(t *testing.T) {
	t.Parallel()

	large := uint64(4294967296)
	tests := []struct {
		in       string
		val      uint
		success  bool
		rangeErr bool
	}{
		{"0", 0, true, false},
		{"42", 42, true, false},
		{strconv.FormatUint(uint64(math.MaxUint), 10), math.MaxUint, true, false},
		{"18446744073709551616", 0, false, true},
	}

	if strconv.IntSize == 32 {
		tests = append(tests, struct {
			in       string
			val      uint
			success  bool
			rangeErr bool
		}{strconv.FormatUint(large, 10), 0, false, true})
	} else {
		tests = append(tests, struct {
			in       string
			val      uint
			success  bool
			rangeErr bool
		}{strconv.FormatUint(large, 10), uint(large), true, false})
	}

	for _, tt := range tests {
		v, err := ParseNativeUint(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		} else if tt.rangeErr {
			require.ErrorIs(t, err, strconv.ErrRange)
		}

		b, err := ParseNativeUint([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		} else if tt.rangeErr {
			require.ErrorIs(t, err, strconv.ErrRange)
		}
	}
}

func Benchmark_ParseUint(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint(input)
			if err != nil {
				b.Fatal("failed to parse uint")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint([]byte(input))
			if err != nil {
				b.Fatal("failed to parse uint from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     int64
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"-42", -42, true},
		{"123456789", 123456789, true},
		{strconv.FormatInt(math.MaxInt64, 10), math.MaxInt64, true},
		{strconv.FormatInt(math.MinInt64, 10), math.MinInt64, true},
		{"9223372036854775808", 0, false},
		{"-9223372036854775809", 0, false},
		{"+10", 10, true},
		{"", 0, false},
		{"12a", 0, false},
		{"-", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseInt(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseInt([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Test_ParseInt_SignOnly(t *testing.T) {
	t.Parallel()

	tests := []string{"+", "-"}
	for _, in := range tests {
		v, err := ParseInt(in)
		require.Error(t, err)
		require.Equal(t, int64(0), v)
		b, err := ParseInt([]byte(in))
		require.Error(t, err)
		require.Equal(t, int64(0), b)
	}
}

func Test_ParseInt_Whitespace(t *testing.T) {
	t.Parallel()

	v, err := ParseInt(" 42")
	require.Error(t, err)
	require.Equal(t, int64(0), v)
}

func Test_ParseNativeInt(t *testing.T) {
	t.Parallel()

	large := int64(3000000000)
	tests := []struct {
		in       string
		val      int
		success  bool
		rangeErr bool
	}{
		{"0", 0, true, false},
		{"42", 42, true, false},
		{strconv.FormatInt(int64(math.MaxInt), 10), math.MaxInt, true, false},
		{strconv.FormatInt(int64(math.MinInt), 10), math.MinInt, true, false},
		{"9223372036854775808", 0, false, true},
	}

	if strconv.IntSize == 32 {
		tests = append(
			tests,
			struct {
				in       string
				val      int
				success  bool
				rangeErr bool
			}{strconv.FormatInt(large, 10), 0, false, true},
			struct {
				in       string
				val      int
				success  bool
				rangeErr bool
			}{strconv.FormatInt(-large, 10), 0, false, true},
		)
	} else {
		tests = append(
			tests,
			struct {
				in       string
				val      int
				success  bool
				rangeErr bool
			}{strconv.FormatInt(large, 10), int(large), true, false},
			struct {
				in       string
				val      int
				success  bool
				rangeErr bool
			}{strconv.FormatInt(-large, 10), int(-large), true, false},
		)
	}

	for _, tt := range tests {
		v, err := ParseNativeInt(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		} else if tt.rangeErr {
			require.ErrorIs(t, err, strconv.ErrRange)
		}

		b, err := ParseNativeInt([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		} else if tt.rangeErr {
			require.ErrorIs(t, err, strconv.ErrRange)
		}
	}
}

func Test_ParseUnsigned_SignOnly(t *testing.T) {
	t.Parallel()

	in := "+"
	v, err := ParseUint(in)
	require.Error(t, err)
	require.Equal(t, uint64(0), v)
	b, err := ParseUint([]byte(in))
	require.Error(t, err)
	require.Equal(t, uint64(0), b)
}

func Benchmark_ParseInt(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt(input)
			if err != nil {
				b.Fatal("failed to parse int")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt([]byte(input))
			if err != nil {
				b.Fatal("failed to parse int from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseInt32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     int32
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"2147483647", 2147483647, true},
		{"-2147483648", -2147483648, true},
		{"2147483648", 0, false},
		{"-2147483649", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseInt32(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseInt32([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseInt32(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt32(input)
			if err != nil {
				b.Fatal("failed to parse int32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt32([]byte(input))
			if err != nil {
				b.Fatal("failed to parse int32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseInt(input, 10, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseInt16(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     int16
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"32767", 32767, true},
		{"-32768", -32768, true},
		{"32768", 0, false},
		{"-32769", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseInt16(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		bts, err := ParseInt16([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, bts)
		}
	}
}

func Benchmark_ParseInt16(b *testing.B) {
	input := "12345"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt16(input)
			if err != nil {
				b.Fatal("failed to parse int16")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt16([]byte(input))
			if err != nil {
				b.Fatal("failed to parse int16 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseInt(input, 10, 16)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseInt8(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     int8
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"127", 127, true},
		{"-128", -128, true},
		{"128", 0, false},
		{"-129", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseInt8(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseInt8([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseInt8(b *testing.B) {
	input := "127"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt8(input)
			if err != nil {
				b.Fatal("failed to parse int8")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseInt8([]byte(input))
			if err != nil {
				b.Fatal("failed to parse int8 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseInt(input, 10, 8)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseUint32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     uint32
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"4294967295", 4294967295, true},
		{"4294967296", 0, false},
		{"-1", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseUint32(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseUint32([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseUint32(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint32(input)
			if err != nil {
				b.Fatal("failed to parse uint32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint32([]byte(input))
			if err != nil {
				b.Fatal("failed to parse uint32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseUint(input, 10, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseUint16(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     uint16
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"65535", 65535, true},
		{"65536", 0, false},
		{"-1", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseUint16(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		bts, err := ParseUint16([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, bts)
		}
	}
}

func Benchmark_ParseUint16(b *testing.B) {
	input := "65535"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint16(input)
			if err != nil {
				b.Fatal("failed to parse uint16")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint16([]byte(input))
			if err != nil {
				b.Fatal("failed to parse uint16 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseUint(input, 10, 16)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseUint8(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     uint8
		success bool
	}{
		{"0", 0, true},
		{"42", 42, true},
		{"255", 255, true},
		{"256", 0, false},
		{"-1", 0, false},
	}
	for _, tt := range tests {
		v, err := ParseUint8(tt.in)
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, v)
		}
		b, err := ParseUint8([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseUint8(b *testing.B) {
	input := "255"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint8(input)
			if err != nil {
				b.Fatal("failed to parse uint8")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseUint8([]byte(input))
			if err != nil {
				b.Fatal("failed to parse uint8 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseUint(input, 10, 8)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseFloat64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     float64
		success bool
	}{
		{"0", 0, true},
		{"42.5", 42.5, true},
		{"-3.14", -3.14, true},
		{"1e2", 100, true},
		{"-1.2e3", -1200, true},
		{"1.", 1, true},
		{".5", 0.5, true},
		{"-0", math.Copysign(0, -1), true},
		{"+123.45", 123.45, true},
		{strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64), math.MaxFloat64, true},
		{strconv.FormatFloat(-math.MaxFloat64, 'g', -1, 64), -math.MaxFloat64, true},
		{"5e-324", 5e-324, true},
		{"2.4703282292062327e-324", 0, true}, // just under half the smallest subnormal
		{"2.4703282292062328e-324", 5e-324, true},
		{"1e-400", 0, true},
		{"25000000000000000000e-18", 25, true},
		{"1234567891234567891234567", 1.2345678912345679e24, true},
		{"1.2345678912345678", 1.2345678912345678, true}, // large number
		{"0.11111111111111119", 0.11111111111111119, true},
		{".", 0, false},
		{"+.", 0, false},
		{"-.", 0, false},
		{"e5", 0, false},
		{".e5", 0, false},
		{"1.2.3", 0, false},
		{"1e1.0", 0, false},
		{"1e309", 0, false},
		{"1e", 0, false},
		{"1e+", 0, false},
		{"1e-", 0, false},
		{"", 0, false},
		{"abc", 0, false},
		{"+", 0, false},
		{"-", 0, false},
		{"123.", 123, true},
		{"123e", 0, false},
		{"123e+", 0, false},
		{"123e-", 0, false},
		{"123e1a", 0, false},
		{"9999999999999999999", 1e19, true},
		{"1.2.3", 0, false},
		{"0e400", 0, true},
		{"-0e400", math.Copysign(0, -1), true},
		{"1e400", 0, false},
		{"0.1e1000", 0, false},
		{strings.Repeat("9", 400), 0, false},
		{strings.Repeat("9", 400) + "e-390", 1e10, true},
	}
	for _, tt := range tests {
		v, err := ParseFloat64(tt.in)
		require.Equal(t, tt.success, err == nil, "input: %s", tt.in)
		if err == nil {
			require.Equal(t, math.Float64bits(tt.val), math.Float64bits(v), "input: %s", tt.in)
		}
		bts, err := ParseFloat64([]byte(tt.in))
		require.Equal(t, tt.success, err == nil, "input: %s", tt.in)
		if err == nil {
			require.Equal(t, math.Float64bits(tt.val), math.Float64bits(bts), "input: %s", tt.in)
		}
	}
}

func Benchmark_ParseFloat64(b *testing.B) {
	input := "12345.6789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseFloat64(input)
			if err != nil {
				b.Fatal("failed to parse float")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseFloat64([]byte(input))
			if err != nil {
				b.Fatal("failed to parse float from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseFloat(input, 64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseFloat32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		val     float32
		success bool
		negzero bool
	}{
		{"0", 0, true, false},
		{"42.5", 42.5, true, false},
		{"-3.14", -3.14, true, false},
		{"1e2", 100, true, false},
		{"-1.2e3", -1200, true, false},
		{"1.", 1, true, false},
		{".5", 0.5, true, false},
		{"3.4028234e38", 3.4028234e38, true, false},
		{"-3.4028234e38", -3.4028234e38, true, false},
		// Java's Float.MAX_VALUE, less than half an ulp above MaxFloat32: rounds down.
		{"3.4028235e38", math.MaxFloat32, true, false},
		{"3.4028236e38", 0, false, false},
		{"-0", -0, true, true},
		{"1e39", 0, false, false},
		{"1e", 0, false, false},
		{"1e+", 0, false, false},
		{"1e-", 0, false, false},
		{"1234567891234567891234567", 1.2345679e24, true, false},
		{"1.2345678912345678", 1.2345679, true, false},
		{"0.11111111111111119", 0.11111111, true, false},
		{".", 0, false, false},
		{"e5", 0, false, false},
		{"1.2.3", 0, false, false},
		{"1e1.0", 0, false, false},
		{"1e309", 0, false, false},
		{"+123.45", 123.45, true, false},
		{"1e-400", 0, true, false},
		{"", 0, false, false},
		{"abc", 0, false, false},
		{"+", 0, false, false},
		{"-", 0, false, false},
		{"123.", 123, true, false},
		{"123e", 0, false, false},
		{"123e+", 0, false, false},
		{"123e-", 0, false, false},
		{"123e1a", 0, false, false},
		{"9999999999999999999", 1e19, true, false},
		{"1.2.3", 0, false, false},
	}
	for _, tt := range tests {
		v, err := ParseFloat32(tt.in)
		require.Equal(t, tt.success, err == nil, "input: %s", tt.in)
		want := tt.val
		if tt.negzero {
			want = float32(math.Copysign(0, -1))
		}
		if err == nil {
			require.Equal(t, math.Float32bits(want), math.Float32bits(v), "input: %s", tt.in)
		}
		bts, err := ParseFloat32([]byte(tt.in))
		require.Equal(t, tt.success, err == nil, "input: %s", tt.in)
		if err == nil {
			require.Equal(t, math.Float32bits(want), math.Float32bits(bts), "input: %s", tt.in)
		}
	}
}

func Benchmark_ParseFloat32(b *testing.B) {
	input := "12345.6789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseFloat32(input)
			if err != nil {
				b.Fatal("failed to parse float32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := ParseFloat32([]byte(input))
			if err != nil {
				b.Fatal("failed to parse float32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := strconv.ParseFloat(input, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_ParseInt8_ParseUint8_EdgeCases(t *testing.T) {
	t.Parallel()

	// Empty, sign-only, and non-digit inputs on the short fast path.
	for _, s := range []string{"", "+", "-", "1x", "+2y", "-!3"} {
		_, err := ParseInt8(s)
		require.Error(t, err, "ParseInt8(%q)", s)
	}
	for _, s := range []string{"", "2x", "!25"} {
		_, err := ParseUint8(s)
		require.Error(t, err, "ParseUint8(%q)", s)
	}

	// Explicit '+' sign and the exact negative boundary.
	v, err := ParseInt8("+127")
	require.NoError(t, err)
	require.Equal(t, int8(127), v)
	v, err = ParseInt8("-128")
	require.NoError(t, err)
	require.Equal(t, int8(-128), v)
	v, err = ParseInt8("-5")
	require.NoError(t, err)
	require.Equal(t, int8(-5), v)
	_, err = ParseInt8("-129")
	require.Error(t, err)
	_, err = ParseInt8("+128")
	require.Error(t, err)

	// Inputs longer than the short fast path fall back to the generic
	// parsers, including leading zeros and range errors.
	v, err = ParseInt8("0127")
	require.NoError(t, err)
	require.Equal(t, int8(127), v)
	v, err = ParseInt8("-00128")
	require.NoError(t, err)
	require.Equal(t, int8(-128), v)
	_, err = ParseInt8("00128")
	require.Error(t, err)

	u, err := ParseUint8("0255")
	require.NoError(t, err)
	require.Equal(t, uint8(255), u)
	_, err = ParseUint8("0256")
	require.Error(t, err)
	_, err = ParseUint8("025x")
	require.Error(t, err)
}

// The tests below cover the 8-digit SWAR parse step in ParseUint/ParseInt.
//
// A SWAR pre-validation for IsIPv4/IsIPv6 was prototyped alongside this work
// and dropped: it lost on every shape (+10% to +49%) because the scalar
// octet parse still has to run after the word checks on <= 15-byte inputs.
// ParseUint won -16% to -27% at 8+ digits with a one-branch cost on 1-digit
// inputs; the raw comparison lives in the introducing commit message.

func Test_Parse8Digits_And_IsEightDigits(t *testing.T) {
	t.Parallel()
	// checkParse8 zero-pads v to eight digits and pins both SWAR helpers on
	// the resulting word.
	checkParse8 := func(v uint64) {
		t.Helper()
		s := FormatUint(v)
		for len(s) < 8 {
			s = "0" + s
		}
		w := swar.Load8(s, 0)
		require.True(t, isEightDigits(w), s)
		require.Equal(t, v, parse8Digits(w), s)
	}

	for _, v := range []uint64{0, 1, 9, 10, 12345678, 99999999, 10000000, 90000009, 9999990} {
		checkParse8(v)
	}
	// isEightDigits must reject every byte value outside '0'..'9' in every
	// lane, including bytes >= 0xFA whose +6 carries into the next lane.
	for c := range 256 {
		lanes := []byte("11111111")
		for lane := range 8 {
			lanes[lane] = byte(c)
			got := isEightDigits(swar.Load8(lanes, 0))
			require.Equal(t, c >= '0' && c <= '9', got, "byte 0x%02x lane %d", c, lane)
			lanes[lane] = '1'
		}
	}

	// Pin parse8Digits itself on random 8-digit words, independent of the
	// full-parser parity loop.
	rng := rand.New(rand.NewSource(5)) //nolint:gosec // deterministic test data
	for range 10000 {
		checkParse8(uint64(rng.Intn(100000000)))
	}
}

func Test_ParseSWAR_Parity(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(3)) //nolint:gosec // deterministic test data
	for range 200000 {
		v := rng.Uint64() >> rng.Intn(60) // cover all digit counts
		s := strconv.FormatUint(v, 10)
		got, err := ParseUint(s)
		require.NoError(t, err, s)
		require.Equal(t, v, got, s)

		gotI, err := ParseInt(s)
		wantI, wantErr := strconv.ParseInt(s, 10, 64)
		require.Equal(t, wantErr == nil, err == nil, s)
		if wantErr == nil {
			// On range errors strconv returns the clamped value while utils
			// returns 0, so values are only comparable on success.
			require.Equal(t, wantI, gotI, s)
		}
	}

	// Error parity with strconv on syntax/overflow shapes that exercise the
	// SWAR word steps, including the scalar handoff around 16-19 digits.
	for _, s := range []string{
		"", "a", "12345678a", "123456781234567812345", "18446744073709551615",
		"18446744073709551616", "99999999999999999999", "00000000000000000001",
		"1234567x", "12345678\n2345678", "184467440737095516150",
		"0000000000000000000000042", "999999999999999999999999",
	} {
		want, wantErr := strconv.ParseUint(s, 10, 64)
		got, gotErr := ParseUint(s)
		require.Equal(t, wantErr == nil, gotErr == nil, "%q", s)
		if wantErr == nil {
			require.Equal(t, want, got, "%q", s)
		} else {
			require.Zero(t, got, "%q", s)
			// The error must keep its strconv shape: *NumError carrying the
			// function name and the ErrSyntax/ErrRange cause.
			var numErr *strconv.NumError
			require.ErrorAs(t, gotErr, &numErr, "%q", s)
			require.Equal(t, "ParseUint", numErr.Func, "%q", s)
			var wantNum *strconv.NumError
			require.ErrorAs(t, wantErr, &wantNum)
			require.ErrorIs(t, numErr.Err, wantNum.Err, "%q", s)
		}
	}
}

func Test_ParseInt_SWARPaths(t *testing.T) {
	t.Parallel()
	// Fast path: a non-digit inside the first word falls back to the scalar
	// loop, which reports the syntax error.
	for _, s := range []string{"1234x678", "12345678\x00", "123456789012345678x"} {
		_, err := ParseInt(s)
		require.Error(t, err, s)
		want, wantErr := strconv.ParseInt(s, 10, 64)
		require.Error(t, wantErr, s)
		require.Equal(t, want, int64(0), s) // strconv returns 0 on syntax errors too
	}

	// Signed inputs with 8+ digit runs route through parseSigned's word path.
	for _, s := range []string{
		"-123456789", "+123456789012", "-9223372036854775808", "+9223372036854775807",
		"-9223372036854775809", "-99999999999999999999", "-12345678a", "-000000001",
	} {
		got, gotErr := ParseInt(s)
		want, wantErr := strconv.ParseInt(s, 10, 64)
		require.Equal(t, wantErr == nil, gotErr == nil, s)
		if wantErr == nil {
			require.Equal(t, want, got, s)
		} else {
			require.Zero(t, got, s)
		}
	}

	// The same dispatch feeds the narrower signed parsers.
	v32, err := ParseInt32("-214748364")
	require.NoError(t, err)
	require.Equal(t, int32(-214748364), v32)
	_, err = ParseInt32("-21474836480")
	require.Error(t, err)

	// 19-digit boundary on the sign-free fast path: MaxInt64 parses, one
	// past it is a range error (not a syntax error), matching strconv.
	maxV, err := ParseInt("9223372036854775807")
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), maxV)
	_, err = ParseInt("9223372036854775808")
	var numErr *strconv.NumError
	require.ErrorAs(t, err, &numErr)
	require.ErrorIs(t, numErr.Err, strconv.ErrRange)
}

func Benchmark_ParseUint_DigitRuns(b *testing.B) {
	nums := map[string]string{
		"1digit":  "8",
		"4digit":  "8080",
		"6digit":  "123456",
		"8digit":  "12345678",
		"12digit": "123456789012",
		"19digit": "1234567890123456789",
	}
	for name, s := range nums {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseUint(s); err != nil {
					b.Fatal("failed to parse uint")
				}
			}
		})
	}
}

// The integer parsers answer as strconv does, except that an error returns 0 where strconv
// clamps, and the Func of the *strconv.NumError names the parser that failed.

type intParser struct {
	name   string
	bits   int
	signed bool
	str    func(string) (int64, uint64, error)
	bytes  func([]byte) (int64, uint64, error)
}

func signedParser[T Signed](name string, bits int, str func(string) (T, error), bytes func([]byte) (T, error)) intParser {
	return intParser{
		name: name, bits: bits, signed: true,
		str: func(s string) (int64, uint64, error) {
			v, err := str(s)
			return int64(v), 0, err
		},
		bytes: func(b []byte) (int64, uint64, error) {
			v, err := bytes(b)
			return int64(v), 0, err
		},
	}
}

func unsignedParser[T Unsigned](name string, bits int, str func(string) (T, error), bytes func([]byte) (T, error)) intParser {
	return intParser{
		name: name, bits: bits,
		str: func(s string) (int64, uint64, error) {
			v, err := str(s)
			return 0, uint64(v), err
		},
		bytes: func(b []byte) (int64, uint64, error) {
			v, err := bytes(b)
			return 0, uint64(v), err
		},
	}
}

var intParsers = []intParser{
	unsignedParser("ParseUint8", 8, ParseUint8[string], ParseUint8[[]byte]),
	unsignedParser("ParseUint16", 16, ParseUint16[string], ParseUint16[[]byte]),
	unsignedParser("ParseUint32", 32, ParseUint32[string], ParseUint32[[]byte]),
	unsignedParser("ParseUint", 64, ParseUint[string], ParseUint[[]byte]),
	unsignedParser("ParseNativeUint", strconv.IntSize, ParseNativeUint[string], ParseNativeUint[[]byte]),
	signedParser("ParseInt8", 8, ParseInt8[string], ParseInt8[[]byte]),
	signedParser("ParseInt16", 16, ParseInt16[string], ParseInt16[[]byte]),
	signedParser("ParseInt32", 32, ParseInt32[string], ParseInt32[[]byte]),
	signedParser("ParseInt", 64, ParseInt[string], ParseInt[[]byte]),
	signedParser("ParseNativeInt", strconv.IntSize, ParseNativeInt[string], ParseNativeInt[[]byte]),
}

func assertIntParity(tb testing.TB, p intParser, in string) {
	tb.Helper()
	var wantI int64
	var wantU uint64
	var wantErr error
	if p.signed {
		wantI, wantErr = strconv.ParseInt(in, 10, p.bits)
	} else {
		wantU, wantErr = strconv.ParseUint(in, 10, p.bits)
	}
	var wantNum *strconv.NumError
	if wantErr != nil && !errors.As(wantErr, &wantNum) {
		tb.Fatalf("strconv returned %v, not a *strconv.NumError", wantErr)
	}

	gotI, gotU, err := p.str(in)
	gotBI, gotBU, errB := p.bytes([]byte(in))
	if gotI != gotBI || gotU != gotBU || (err == nil) != (errB == nil) {
		tb.Fatalf("%s(%q): string form (%d, %d, %v) differs from bytes form (%d, %d, %v)", p.name, in, gotI, gotU, err, gotBI, gotBU, errB)
	}
	if wantErr != nil && !errors.Is(errB, wantNum.Err) {
		tb.Fatalf("%s(bytes %q) error = %v, strconv gives %v", p.name, in, errB, wantErr)
	}

	if wantErr == nil {
		if err != nil || gotI != wantI || gotU != wantU {
			tb.Fatalf("%s(%q) = (%d, %d, %v), strconv gives (%d, %d, nil)", p.name, in, gotI, gotU, err, wantI, wantU)
		}
		return
	}
	var num *strconv.NumError
	switch {
	case !errors.As(err, &num):
		tb.Fatalf("%s(%q) error = %v, strconv gives %v", p.name, in, err, wantErr)
	case !errors.Is(err, wantNum.Err):
		tb.Fatalf("%s(%q) error = %v, strconv gives %v", p.name, in, err, wantErr)
	case num.Func != p.name || num.Num != in:
		tb.Fatalf("%s(%q) error = %#v, want Func %q and Num %q", p.name, in, num, p.name, in)
	case gotI != 0 || gotU != 0:
		tb.Fatalf("%s(%q) returned (%d, %d) alongside %v", p.name, in, gotI, gotU, err)
	}
}

func assertAllIntParity(tb testing.TB, in string) {
	tb.Helper()
	for _, p := range intParsers {
		assertIntParity(tb, p, in)
	}
}

func Test_ParseInts_StrconvParity_Exhaustive(t *testing.T) {
	t.Parallel()
	forEachString("0159+-x ", 5, func(s string) {
		assertAllIntParity(t, s)
	})
}

func Test_ParseInts_StrconvParity_EveryByte(t *testing.T) {
	t.Parallel()
	// Every byte at every position of each run, on both the word and scalar paths.
	runs := []string{"7", "255", "256", "65536", "4294967296", "12345678", "123456789", "9223372036854775808", "18446744073709551616", "99999999999999999999999"}
	for b := range 256 {
		c := string([]byte{byte(b)})
		for _, run := range runs {
			for cut := 0; cut <= len(run); cut++ {
				assertAllIntParity(t, run[:cut]+c+run[cut:])
			}
			assertAllIntParity(t, c+run)
			assertAllIntParity(t, run+c)
			assertAllIntParity(t, "-"+run+c)
		}
	}
}

func Test_ParseInts_StrconvParity_Boundaries(t *testing.T) {
	t.Parallel()
	var limits []*big.Int
	for _, bits := range []uint{8, 16, 32, 63, 64, 65, 128} {
		hi := new(big.Int).Lsh(big.NewInt(1), bits)
		half := new(big.Int).Rsh(hi, 1)
		for _, base := range []*big.Int{hi, half} {
			for _, delta := range []int64{-2, -1, 0, 1, 2} {
				limits = append(limits, new(big.Int).Add(base, big.NewInt(delta)))
			}
		}
	}
	limits = append(limits, big.NewInt(0), big.NewInt(1), big.NewInt(9), big.NewInt(10), big.NewInt(99), big.NewInt(100), big.NewInt(100000000))
	for _, limit := range limits {
		digits := limit.String()
		for _, sign := range []string{"", "+", "-"} {
			for _, zeros := range []int{0, 1, 7, 8, 9, 19, 20, 25} {
				for _, tail := range []string{"", "x", "0", "9", " ", "+", "-", "_", ".", "e1", "\x00", "\xff"} {
					assertAllIntParity(t, sign+strings.Repeat("0", zeros)+digits+tail)
				}
			}
		}
	}
	for _, in := range []string{"", "+", "-", "++1", "--1", "+-1", "-+1", " 1", "1 ", "0x10", "1_000", "\u0661\u0662\u0663", "1e3", "1.0", "\uff11\uff12\uff13"} {
		assertAllIntParity(t, in)
	}
}

func Test_ParseInts_StrconvParity_Random(t *testing.T) {
	t.Parallel()
	r := randv2.New(randv2.NewPCG(0x1a7, 0x9a45e)) //nolint:gosec // deterministic test data
	const strays = "x+-_. \x00\xff"
	for range 60000 {
		b := make([]byte, r.IntN(42))
		for i := range b {
			b[i] = byte('0' + r.IntN(10))
		}
		switch r.IntN(4) {
		case 0:
			if len(b) > 0 {
				b[r.IntN(len(b))] = strays[r.IntN(len(strays))]
			}
		case 1:
			if len(b) > 0 {
				b[0] = "+-"[r.IntN(2)]
			}
			if len(b) > 1 && r.IntN(2) == 0 {
				b[1+r.IntN(len(b)-1)] = strays[r.IntN(len(strays))]
			}
		}
		assertAllIntParity(t, string(b))
	}
}

func FuzzParseIntWidths(f *testing.F) {
	for _, seed := range []string{
		"256x", "0500x", "65536+", "4294967296 ", "9223372036854775808x", "18446744073709551616x",
		"-129x", "-128", "128", "255", "-256x", "+", "-", "", "00000000000000000000000256x",
		"12345678x", "1234567890123456789012x",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		assertAllIntParity(t, s)
	})
}

// Cases that used to report a syntax error: digits over the unsigned limit before a stray
// byte are a range error.
func Test_ParseInts_RangeBeforeSyntax(t *testing.T) {
	t.Parallel()
	check := func(in string, bits int, signedOnly bool, want error) {
		t.Helper()
		for _, p := range intParsers {
			if p.bits != bits || (signedOnly && !p.signed) {
				continue
			}
			if _, _, err := p.str(in); !errors.Is(err, want) {
				t.Errorf("%s(%q) error = %v, want %v", p.name, in, err, want)
			}
		}
	}
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"256x", 8},
		{"0500x", 8},
		{"2560 ", 8},
		{"256+", 8},
		{"65536x", 16},
		{"00065536x", 16},
		{"4294967296x", 32},
		{"99999999999x", 32},
		{"00000000000000000000000000256x", 8},
		{"18446744073709551616x", 64},
	} {
		check(tc.in, tc.bits, false, strconv.ErrRange)
	}
	// A sign is a syntax error to the unsigned parsers, which stop there.
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"-256x", 8}, {"+256-", 8}, {"-4294967296x", 32}, {"+65536x", 16},
	} {
		check(tc.in, tc.bits, true, strconv.ErrRange)
		for _, p := range intParsers {
			if p.bits == tc.bits && !p.signed {
				if _, _, err := p.str(tc.in); !errors.Is(err, strconv.ErrSyntax) {
					t.Errorf("%s(%q) error = %v, want a syntax error", p.name, tc.in, err)
				}
			}
		}
	}
	// Within the unsigned limit, which strconv compares with, a stray byte is a syntax error
	// even where the signed type overflows.
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"255x", 8},
		{"128x", 8},
		{"65535x", 16},
		{"32768x", 16},
		{"4294967295x", 32},
		{"2147483648x", 32},
		{"9223372036854775808x", 64},
		{"18446744073709551615x", 64},
	} {
		check(tc.in, tc.bits, false, strconv.ErrSyntax)
	}
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"-129x", 8}, {"-32769x", 16}, {"-2147483649x", 32}, {"-9223372036854775809x", 64},
	} {
		check(tc.in, tc.bits, true, strconv.ErrSyntax)
	}
}
