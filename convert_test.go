// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_UnsafeString(t *testing.T) {
	t.Parallel()
	res := UnsafeString([]byte("Hello, World!"))
	require.Equal(t, "Hello, World!", res)
}

// go test -v -run=^$ -bench=UnsafeString -benchmem -count=2

func Benchmark_UnsafeString(b *testing.B) {
	hello := []byte("Hello, World!")
	var res string
	b.Run("unsafe", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = UnsafeString(hello)
		}
		require.Equal(b, "Hello, World!", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = string(hello)
		}
		require.Equal(b, "Hello, World!", res)
	})
}

func Test_UnsafeBytes(t *testing.T) {
	t.Parallel()
	res := UnsafeBytes("Hello, World!")
	require.Equal(t, []byte("Hello, World!"), res)
}

// go test -v -run=^$ -bench=UnsafeBytes -benchmem -count=4

func Benchmark_UnsafeBytes(b *testing.B) {
	hello := "Hello, World!"
	var res []byte
	b.Run("unsafe", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = UnsafeBytes(hello)
		}
		require.Equal(b, []byte("Hello, World!"), res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = []byte(hello)
		}
		require.Equal(b, []byte("Hello, World!"), res)
	})
}

func Test_CopyString(t *testing.T) {
	t.Parallel()
	res := CopyString("Hello, World!")
	require.Equal(t, "Hello, World!", res)
}

func Test_ToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    interface{}
		expected string
	}{
		{[]byte("Hello, World!"), "Hello, World!"},
		{true, "true"},
		{uint(100), "100"},
		{int(42), "42"},
		{int8(42), "42"},
		{int16(42), "42"},
		{int32(42), "42"},
		{int64(42), "42"},
		{uint8(100), "100"},
		{uint16(100), "100"},
		{uint32(100), "100"},
		{uint64(100), "100"},
		{"test string", "test string"},
		{float32(3.14), "3.14"},
		{float64(3.14159), "3.14159"},
		{time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC), "2000-01-01 12:34:56"},
		{struct{ Name string }{"John"}, "{John}"},
	}

	for _, tc := range tests {
		t.Run(reflect.TypeOf(tc.input).String(), func(t *testing.T) {
			res := ToString(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}

	// Testing pointer to int
	intPtr := 42
	testsPtr := []struct {
		input    interface{}
		expected string
	}{
		{&intPtr, "42"},
	}
	for _, tc := range testsPtr {
		t.Run("pointer to "+reflect.TypeOf(tc.input).Elem().String(), func(t *testing.T) {
			res := ToString(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}
}

// go test -v -run=^$ -bench=ToString -benchmem -count=4
func Benchmark_ToString(b *testing.B) {
	values := []interface{}{
		42,
		int8(42),
		int16(42),
		int32(42),
		int64(42),
		uint(42),
		uint8(42),
		uint16(42),
		uint32(42),
		uint64(42),
		"Hello, World!",
		[]byte("Hello, World!"),
		true,
		float32(3.14),
		float64(3.14),
		time.Now(),
	}
	for n := 0; n < b.N; n++ {
		for _, value := range values {
			_ = ToString(value)
		}
	}
}
