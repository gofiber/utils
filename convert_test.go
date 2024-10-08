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

var dataTypeExamples []any

// MyStringer for testing fmt.Stringer support.
type MyStringer struct {
	value string
}

// CustomType not handled by switch cases.
type CustomType struct {
	Field1 string
	Field2 int
}

func (ms MyStringer) String() string {
	return ms.value
}

func init() {
	dataTypeExamples = []any{
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
		time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
		[]string{"Hello", "World"},
		[]int{42, 21},
		[2]int{42, 21},
		[][]int{{42, 21}, {42, 21}},
		[]any{[]int{42, 21}, 42, "Hello", true, []string{"Hello", "World"}},
		MyStringer{value: "Stringer Value"},
		CustomType{Field1: "example", Field2: 42},
	}
}

func Test_UnsafeString(t *testing.T) {
	t.Parallel()
	res := UnsafeString([]byte("Hello, World!"))
	require.Equal(t, "Hello, World!", res)
}

func Test_UnsafeBytes(t *testing.T) {
	t.Parallel()
	res := UnsafeBytes("Hello, World!")
	require.Equal(t, []byte("Hello, World!"), res)
}

func Test_CopyString(t *testing.T) {
	t.Parallel()
	res := CopyString("Hello, World!")
	require.Equal(t, "Hello, World!", res)
}

func Test_ToString(t *testing.T) {
	t.Parallel()

	customInstance := CustomType{Field1: "example", Field2: 42}
	timeSample := time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC)
	stringerSample := MyStringer{value: "Stringer Value"}

	tests := []struct {
		name       string
		input      any
		timeFormat []string
		expected   string
	}{
		// Primitive types and string
		{name: "int", input: int(42), expected: "42"},
		{name: "int8", input: int8(42), expected: "42"},
		{name: "int16", input: int16(42), expected: "42"},
		{name: "int32", input: int32(42), expected: "42"},
		{name: "int64", input: int64(42), expected: "42"},
		{name: "uint", input: uint(100), expected: "100"},
		{name: "uint8", input: uint8(100), expected: "100"},
		{name: "uint16", input: uint16(100), expected: "100"},
		{name: "uint32", input: uint32(100), expected: "100"},
		{name: "uint64", input: uint64(100), expected: "100"},
		{name: "string", input: "test string", expected: "test string"},
		{name: "[]byte", input: []byte("Hello, World!"), expected: "Hello, World!"},
		{name: "bool", input: true, expected: "true"},
		{name: "float32", input: float32(3.14), expected: "3.14"},
		{name: "float64", input: float64(3.14159), expected: "3.14159"},

		// time.Time
		{name: "time.Time default format", input: timeSample, expected: "2000-01-01 12:34:56"},
		{name: "time.Time custom format", input: timeSample, timeFormat: []string{"Jan 02, 2006"}, expected: "Jan 01, 2000"},

		// reflect.Value
		{name: "reflect.Value", input: reflect.ValueOf(42), expected: "42"},

		// fmt.Stringer
		{name: "fmt.Stringer", input: stringerSample, expected: "Stringer Value"},

		// Composite types (arrays, slices)
		{name: "[]string", input: []string{"Hello", "World"}, expected: "[Hello World]"},
		{name: "[]int", input: []int{42, 21}, expected: "[42 21]"},
		{name: "[][]int", input: [][]int{{42, 21}, {42, 21}}, expected: "[[42 21] [42 21]]"},
		{name: "[]any", input: []any{[]int{42, 21}, 42, "Hello", true, []string{"Hello", "World"}}, expected: "[[42 21] 42 Hello true [Hello World]]"},

		// Custom unhandled type
		{name: "CustomType", input: customInstance, expected: "{example 42}"},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var res string
			if len(tc.timeFormat) > 0 {
				res = ToString(tc.input, tc.timeFormat...)
			} else {
				res = ToString(tc.input)
			}
			require.Equal(t, tc.expected, res)
		})
	}

	// Testing pointer to int
	intPtr := 42
	testsPtr := []struct {
		input    any
		expected string
	}{
		{&intPtr, "42"},
	}
	for _, tc := range testsPtr {
		tc := tc
		t.Run("pointer to "+reflect.TypeOf(tc.input).Elem().String(), func(t *testing.T) {
			t.Parallel()
			res := ToString(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}
}

func TestCopyBytes(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		input := []byte{}
		copied := CopyBytes(input)
		require.Equal(t, input, copied)
	})

	t.Run("single element", func(t *testing.T) {
		input := []byte{42}
		copied := CopyBytes(input)
		require.Equal(t, input, copied)
		input[0] = 0 // Modify the input to ensure the copied slice does not change
		require.NotEqual(t, input[0], copied[0])
	})

	t.Run("multiple elements", func(t *testing.T) {
		input := []byte{1, 2, 3, 4, 5}
		copied := CopyBytes(input)
		require.Equal(t, input, copied)
		input[0] = 0 // Modify the input to ensure the copied slice does not change
		require.NotEqual(t, input, copied)
	})

	t.Run("deep copy validation", func(t *testing.T) {
		input := []byte{1, 2, 3, 4, 5}
		copied := CopyBytes(input)
		input[0] = 0 // Modify the input to ensure the copied slice does not change
		require.NotEqual(t, input[0], copied[0])
	})

	t.Run("nil slice", func(t *testing.T) {
		copied := CopyBytes(nil)
		require.NotNil(t, copied)
		require.Empty(t, copied)
		require.Equal(t, 0, cap(copied))
	})
}

func TestByteSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"Zero Byte", 0, "0B"},
		{"One Byte", 1, "1B"},
		{"Kilobytes", 1024, "1KB"},
		{"Megabytes", 1024 * 1024, "1MB"},
		{"Gigabytes", 1024 * 1024 * 1024, "1GB"},
		{"Terabytes", 1024 * 1024 * 1024 * 1024, "1TB"},
		{"Petabytes", 1024 * 1024 * 1024 * 1024 * 1024, "1PB"},
		{"Exabytes", 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1EB"},
		{"Bytes less than KB", 500, "500B"},
		{"KB less than MB", 500 * 1024, "500KB"},
		{"MB less than GB", 500 * 1024 * 1024, "500MB"},
		{"GB less than TB", 500 * 1024 * 1024 * 1024, "500GB"},
		{"TB less than PB", 500 * 1024 * 1024 * 1024 * 1024, "500TB"},
		{"PB less than EB", 500 * 1024 * 1024 * 1024 * 1024 * 1024, "500PB"},
		{"Fractional KB", 1126, "1.1KB"},
		{"Fractional MB", 1126 * 1024, "1.1MB"},
		{"Fractional GB", 1126 * 1024 * 1024, "1.1GB"},
		{"Fractional TB", 1126 * 1024 * 1024 * 1024, "1.1TB"},
		{"Fractional PB", 1126 * 1024 * 1024 * 1024 * 1024, "1.1PB"},
		{"Fractional EB", 1126 * 1024 * 1024 * 1024 * 1024 * 1024, "1.1EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByteSize(tt.bytes)
			require.Equal(t, tt.expected, result)
		})
	}
}

// go test -v -run=^$ -bench=ToString -benchmem -count=4
func Benchmark_ToString(b *testing.B) {
	for _, value := range dataTypeExamples {
		b.Run(reflect.TypeOf(value).String(), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ToString(value)
			}
		})
	}
}

// go test -v -run=^$ -bench=ToString_concurrency -benchmem -count=4
func Benchmark_ToString_concurrency(b *testing.B) {
	for _, value := range dataTypeExamples {
		b.Run(reflect.TypeOf(value).String(), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = ToString(value)
				}
			})
		})
	}
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
