package utils

import (
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

func Benchmark_ParseUint(b *testing.B) {
	input := "123456789"

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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

func Benchmark_ParseInt(b *testing.B) {
	input := "123456789"

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseInt(input, 10, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseInt(input, 10, 8)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseUint(input, 10, 32)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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

	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := strconv.ParseUint(input, 10, 8)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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
}
