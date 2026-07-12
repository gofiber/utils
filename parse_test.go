package utils

import (
	"math"
	"math/rand"
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
		tests = append(tests,
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
		tests = append(tests,
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
		{"-0", -0, true},
		{"+123.45", 123.45, true},
		{strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64), math.MaxFloat64, true},
		{strconv.FormatFloat(-math.MaxFloat64, 'g', -1, 64), -math.MaxFloat64, true},
		{"5e-324", 0, true},
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
		{"-0e400", 0, true},
		{"1e400", 0, false},
		{"0.1e1000", 0, false},
		{strings.Repeat("9", 400), 0, false},
		{strings.Repeat("9", 400) + "e-390", 1e10, true},
	}
	for _, tt := range tests {
		v, err := ParseFloat64(tt.in)
		require.Equal(t, tt.success, err == nil, "input: %s", tt.in)
		if err == nil {
			if tt.val == 0 {
				require.InDelta(t, tt.val, v, 1e-9, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, v, 1e-9, "input: %s", tt.in)
			}
		}
		bts, err := ParseFloat64([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
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
		{"3.4028235e38", 0, false, false},
		{"-0", -0, true, true},
		{"1e39", 0, false, false},
		{"1e", 0, false, false},
		{"1e+", 0, false, false},
		{"1e-", 0, false, false},
		{"1234567891234567891234567", 1.2345679e24, true, false},
		{"1.2345678912345678", 1.2345679, true, false},
		{"0.11111111111111119", 0.11111112, true, false},
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
		if err == nil {
			if tt.negzero {
				require.InDelta(t, float32(0), v, 1e-6, "input: %s", tt.in)
				require.True(t, math.Signbit(float64(v)))
			} else if tt.val == 0 {
				require.InDelta(t, tt.val, v, 1e-6, "input: %s", tt.in)
			} else {
				require.InEpsilon(t, tt.val, v, 1e-6, "input: %s", tt.in)
			}
		}
		bts, err := ParseFloat32([]byte(tt.in))
		require.Equal(t, tt.success, err == nil)
		if err == nil {
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
	for _, v := range []uint64{0, 1, 9, 10, 12345678, 99999999, 10000000, 90000009, 9999990} {
		s := FormatUint(v)
		for len(s) < 8 {
			s = "0" + s
		}
		w := swar.Load8(s, 0)
		require.True(t, isEightDigits(w), s)
		require.Equal(t, v, parse8Digits(w), s)
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
		v := uint64(rng.Intn(100000000))
		s := FormatUint(v)
		for len(s) < 8 {
			s = "0" + s
		}
		w := swar.Load8(s, 0)
		require.True(t, isEightDigits(w), s)
		require.Equal(t, v, parse8Digits(w), s)
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
