package utils

import (
	"math"
	"strconv"
	"testing"

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
		v, ok := ParseUint(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseUint([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Test_ParseUint_Whitespace(t *testing.T) {
	t.Parallel()

	v, ok := ParseUint(" 123")
	require.False(t, ok)
	require.Equal(t, uint64(0), v)
}

func Benchmark_ParseUint(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint(input)
			if !ok {
				b.Fatal("failed to parse uint")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint([]byte(input))
			if !ok {
				b.Fatal("failed to parse uint from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseInt(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseInt([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Test_ParseInt_SignOnly(t *testing.T) {
	t.Parallel()

	tests := []string{"+", "-"}
	for _, in := range tests {
		v, ok := ParseInt(in)
		require.False(t, ok)
		require.Equal(t, int64(0), v)
		b, ok := ParseInt([]byte(in))
		require.False(t, ok)
		require.Equal(t, int64(0), b)
	}
}

func Test_ParseInt_Whitespace(t *testing.T) {
	t.Parallel()

	v, ok := ParseInt(" 42")
	require.False(t, ok)
	require.Equal(t, int64(0), v)
}

func Test_ParseUnsigned_SignOnly(t *testing.T) {
	t.Parallel()

	in := "+"
	v, ok := ParseUint(in)
	require.False(t, ok)
	require.Equal(t, uint64(0), v)
	b, ok := ParseUint([]byte(in))
	require.False(t, ok)
	require.Equal(t, uint64(0), b)
}

func Benchmark_ParseInt(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt(input)
			if !ok {
				b.Fatal("failed to parse int")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt([]byte(input))
			if !ok {
				b.Fatal("failed to parse int from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseInt32(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseInt32([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseInt32(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt32(input)
			if !ok {
				b.Fatal("failed to parse int32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt32([]byte(input))
			if !ok {
				b.Fatal("failed to parse int32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseInt16(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		bts, ok := ParseInt16([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, bts)
		}
	}
}

func Benchmark_ParseInt16(b *testing.B) {
	input := "12345"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt16(input)
			if !ok {
				b.Fatal("failed to parse int16")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt16([]byte(input))
			if !ok {
				b.Fatal("failed to parse int16 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseInt8(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseInt8([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseInt8(b *testing.B) {
	input := "127"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt8(input)
			if !ok {
				b.Fatal("failed to parse int8")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseInt8([]byte(input))
			if !ok {
				b.Fatal("failed to parse int8 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseUint32(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseUint32([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseUint32(b *testing.B) {
	input := "123456789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint32(input)
			if !ok {
				b.Fatal("failed to parse uint32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint32([]byte(input))
			if !ok {
				b.Fatal("failed to parse uint32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseUint16(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		bts, ok := ParseUint16([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, bts)
		}
	}
}

func Benchmark_ParseUint16(b *testing.B) {
	input := "65535"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint16(input)
			if !ok {
				b.Fatal("failed to parse uint16")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint16([]byte(input))
			if !ok {
				b.Fatal("failed to parse uint16 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		v, ok := ParseUint8(tt.in)
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, v)
		}
		b, ok := ParseUint8([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			require.Equal(t, tt.val, b)
		}
	}
}

func Benchmark_ParseUint8(b *testing.B) {
	input := "255"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint8(input)
			if !ok {
				b.Fatal("failed to parse uint8")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseUint8([]byte(input))
			if !ok {
				b.Fatal("failed to parse uint8 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		{"-0", -0, true},
		{"+123.45", 123.45, true},
		{strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64), math.MaxFloat64, true},
		{strconv.FormatFloat(-math.MaxFloat64, 'g', -1, 64), -math.MaxFloat64, true},
		{"5e-324", 0, true},
		{"1e-400", 0, true},
		{"1234567891234567891234567", 0, false},
		{"1.2345678912345678", 1.2345678912345678, true}, // large number
		{"0.11111111111111119", 0, false},
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
	}
	for _, tt := range tests {
		v, ok := ParseFloat64(tt.in)
		require.Equal(t, tt.success, ok, "input: %s", tt.in)
		if ok {
			if tt.val == 0 {
				require.InDelta(t, tt.val, v, 1e-9, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, v, 1e-9, "input: %s", tt.in)
			}
		}
		bts, ok := ParseFloat64([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			if tt.val == 0 {
				require.InDelta(t, tt.val, bts, 1e-9, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, bts, 1e-9, "input: %s", tt.in)
			}
		}
	}
}

func Benchmark_ParseFloat64(b *testing.B) {
	input := "12345.6789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseFloat64(input)
			if !ok {
				b.Fatal("failed to parse float")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseFloat64([]byte(input))
			if !ok {
				b.Fatal("failed to parse float from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
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
		{"3.4028235e38", 0, false, false},
		{"-0", -0, true, true},
		{"1e39", 0, false, false},
		{"1e", 0, false, false},
		{"1e+", 0, false, false},
		{"1e-", 0, false, false},
		{"1234567891234567891234567", 0, false, false},
		{"1.2345678912345678", 1.2345679, true, false},
		{"0.11111111111111119", 0, false, false},
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
		v, ok := ParseFloat32(tt.in)
		require.Equal(t, tt.success, ok, "input: %s", tt.in)
		if ok {
			if tt.negzero {
				require.InDelta(t, float32(0), v, 1e-6, "input: %s", tt.in)
				require.True(t, math.Signbit(float64(v)))
			} else if tt.val == 0 {
				require.InDelta(t, tt.val, v, 1e-6, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, v, 1e-6, "input: %s", tt.in)
			}
		}
		bts, ok := ParseFloat32([]byte(tt.in))
		require.Equal(t, tt.success, ok)
		if ok {
			if tt.negzero {
				require.InDelta(t, float32(0), bts, 1e-6, "input: %s", tt.in)
				require.True(t, math.Signbit(float64(bts)))
			} else if tt.val == 0 {
				require.InDelta(t, tt.val, bts, 1e-6, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, bts, 1e-6, "input: %s", tt.in)
			}
		}
	}
}

func Benchmark_ParseFloat32(b *testing.B) {
	input := "12345.6789"

	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseFloat32(input)
			if !ok {
				b.Fatal("failed to parse float32")
			}
		}
	})
	b.Run("fiber_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, ok := ParseFloat32([]byte(input))
			if !ok {
				b.Fatal("failed to parse float32 from bytes")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseFloat(input, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
