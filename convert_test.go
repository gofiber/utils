// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"math/rand/v2"
	"net"
	"net/netip"
	"reflect"
	"strconv"
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

		{name: "time.Time default format", input: timeSample, expected: "2000-01-01 12:34:56"},
		{name: "time.Time custom format", input: timeSample, timeFormat: []string{"Jan 02, 2006"}, expected: "Jan 01, 2000"},

		{name: "reflect.Value", input: reflect.ValueOf(42), expected: "42"},

		{name: "fmt.Stringer", input: stringerSample, expected: "Stringer Value"},

		{name: "error", input: errors.New("something failed"), expected: "something failed"},

		{name: "[]string", input: []string{"Hello", "World"}, expected: "[Hello World]"},
		{name: "[]int", input: []int{42, 21}, expected: "[42 21]"},
		{name: "[][]int", input: [][]int{{42, 21}, {42, 21}}, expected: "[[42 21] [42 21]]"},
		{name: "[]any", input: []any{[]int{42, 21}, 42, "Hello", true, []string{"Hello", "World"}}, expected: "[[42 21] 42 Hello true [Hello World]]"},

		{name: "CustomType", input: customInstance, expected: "{example 42}"},
	}

	for _, tc := range tests {
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

	intPtr := 42
	testsPtr := []struct {
		input    any
		expected string
	}{
		{&intPtr, "42"},
	}
	for _, tc := range testsPtr {
		t.Run("pointer to "+reflect.TypeOf(tc.input).Elem().String(), func(t *testing.T) {
			t.Parallel()
			res := ToString(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		var res string
		require.NotPanics(t, func() {
			res = ToString(nil)
		})
		require.Empty(t, res)
	})

	t.Run("invalid reflect.Value", func(t *testing.T) {
		t.Parallel()
		var res string
		require.NotPanics(t, func() {
			res = ToString(reflect.Value{})
		})
		require.Empty(t, res)
	})

	t.Run("non-interfaceable reflect.Value", func(t *testing.T) {
		t.Parallel()
		type sample struct {
			value string
		}

		field := reflect.ValueOf(sample{value: "secret"}).FieldByName("value")
		require.False(t, field.CanInterface())

		var res string
		require.NotPanics(t, func() {
			res = ToString(field)
		})
		require.Empty(t, res)
	})

	// A nil pointer renders as its type's zero value.
	t.Run("nil pointer to string", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, ToString((*string)(nil)))
	})
	t.Run("nil pointer to int", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "0", ToString((*int)(nil)))
	})
	t.Run("nil pointer to int64", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "0", ToString((*int64)(nil)))
	})
	t.Run("nil pointer to uint64", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "0", ToString((*uint64)(nil)))
	})
	t.Run("nil pointer to float64", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "0", ToString((*float64)(nil)))
	})
	t.Run("nil pointer to bool", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "false", ToString((*bool)(nil)))
	})
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
		input[0] = 0
		require.NotEqual(t, input[0], copied[0])
	})

	t.Run("multiple elements", func(t *testing.T) {
		input := []byte{1, 2, 3, 4, 5}
		copied := CopyBytes(input)
		require.Equal(t, input, copied)
		input[0] = 0
		require.NotEqual(t, input, copied)
	})

	t.Run("deep copy validation", func(t *testing.T) {
		input := []byte{1, 2, 3, 4, 5}
		copied := CopyBytes(input)
		input[0] = 0
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

	// Above math.MaxUint64/10, where the arithmetic was once capped, sizes are written as they are.
	for _, tt := range []struct {
		bytes    uint64
		expected string
	}{
		{math.MaxUint64 / 10, "1.6EB"},
		{math.MaxUint64/10 + 1, "1.6EB"},
		{1902320482601297510, "1.6EB"},
		{1902320482601297511, "1.7EB"},
		{2 << 60, "2EB"},
		{1 << 62, "4EB"},
		{1 << 63, "8EB"},
		{math.MaxUint64, "16EB"},
	} {
		require.Equal(t, tt.expected, ByteSize(tt.bytes), "%d", tt.bytes)
	}
}

// Test_ByteSize_Rounding pins a unit's upper edge: 1023.95 of a unit is written as 1 of the next.
func Test_ByteSize_Rounding(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		bytes    uint64
		expected string
	}{
		{1023, "1023B"},
		{1048524, "1023.9KB"},
		{1048525, "1MB"},
		{1048575, "1MB"},
		{1 << 20, "1MB"},
		{1<<30 - 1, "1GB"},
		{1<<40 - 1, "1TB"},
		{1<<50 - 1, "1PB"},
		{1<<60 - 1, "1EB"},
		{1280, "1.3KB"}, // 1.25 exactly: half rounds up
		{1792, "1.8KB"},
		{1126, "1.1KB"},
	} {
		require.Equal(t, tt.expected, ByteSize(tt.bytes), "%d", tt.bytes)
	}
}

// byteSizeModel recomputes ByteSize with math/big from its documented rules alone.
func byteSizeModel(n uint64) string {
	if n < 1024 {
		return strconv.FormatUint(n, 10) + "B"
	}
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	k := (bits.Len64(n) - 1) / 10
	unit := new(big.Int).Lsh(big.NewInt(1), uint(10*k))
	tenths := new(big.Int).Mul(new(big.Int).SetUint64(n), big.NewInt(10))
	tenths.Add(tenths, new(big.Int).Rsh(unit, 1))
	tenths.Quo(tenths, unit)
	if tenths.Cmp(big.NewInt(10240)) >= 0 {
		tenths.SetInt64(10)
		k++
	}
	whole, frac := tenths.Int64()/10, tenths.Int64()%10
	s := strconv.FormatInt(whole, 10)
	if frac != 0 {
		s += "." + strconv.FormatInt(frac, 10)
	}
	return s + units[k-1]
}

// Test_ByteSize_Model compares ByteSize with byteSizeModel, densely around the unit edges.
func Test_ByteSize_Model(t *testing.T) {
	t.Parallel()
	check := func(n uint64) {
		t.Helper()
		if got, want := ByteSize(n), byteSizeModel(n); got != want {
			t.Fatalf("ByteSize(%d) = %q, want %q", n, got, want)
		}
	}
	for n := range uint64(100_000) {
		check(n)
	}
	for shift := range 64 {
		for d := range uint64(300) {
			base := uint64(1) << shift
			check(base + d)
			check(base - 1 - d)
		}
	}
	for d := range uint64(2000) {
		check(math.MaxUint64 - d)
	}
	rng := rand.New(rand.NewPCG(1, 2)) //nolint:gosec // a fixed sequence for a model test
	for k := range 6 {
		edge := uint64(1) << (10 * (k + 1))
		for range 5000 {
			check(edge - 1 - rng.Uint64N(edge/16)) // 1023.95 of the unit below is edge - edge/20480
		}
	}
	for range 40_000 {
		check(rng.Uint64() >> rng.UintN(64))
	}
}

// Ensure ConvertToBytes caps values exceeding math.MaxInt
func Test_ConvertToBytes_MaxInt(t *testing.T) {
	t.Parallel()
	huge := fmt.Sprintf("%dT", int64(math.MaxInt/1000+1))
	require.Equal(t, math.MaxInt, ConvertToBytes(huge))
}

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

// Test_mayBeNilPointer holds the gc shortcut to nilPointerReflect; a nil map, func or chan may match.
func Test_mayBeNilPointer(t *testing.T) {
	t.Parallel()
	n := 7
	var (
		nilInt  *int
		nilType *CustomType
		nilTime *time.Time
		nilErr  *nilReceiverError
		nilStr  *nilReceiver
	)
	for name, v := range map[string]any{
		"nil *int":           nilInt,
		"nil *struct":        nilType,
		"nil *time.Time":     nilTime,
		"nil error receiver": nilErr,
		"nil Stringer":       nilStr,
		"nil **int":          (**int)(nil),
		"nil *[]int":         (*[]int)(nil),
		"nil *uintptr":       (*uintptr)(nil),
	} {
		require.True(t, nilPointerReflect(v), name)
		require.True(t, mayBeNilPointer(v), name)
	}
	for name, v := range map[string]any{
		"non-nil *int":    &n,
		"non-nil *struct": &CustomType{},
		"int":             0,
		"string":          "",
		"bool":            false,
		"struct":          CustomType{},
		"empty struct":    struct{}{},
		"nil slice":       []int(nil),
		"slice":           []int{1},
		"map":             map[string]int{},
		"func":            func() {},
		"error":           errors.New("x"),
		"time":            time.Time{},
	} {
		require.False(t, nilPointerReflect(v), name)
		require.False(t, mayBeNilPointer(v), name)
	}
	for name, v := range map[string]any{
		"nil map":  map[string]int(nil),
		"nil func": (func())(nil),
		"nil chan": (chan int)(nil),
	} {
		require.False(t, nilPointerReflect(v), name)
	}
}

// Benchmark_ToString_Methods measures the types that go through an Error or String method,
// where ToString first checks for a nil pointer receiver.
func Benchmark_ToString_Methods(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value any
	}{
		{"error", errors.New("boom")},
		{"stringer", MyStringer{value: "Stringer Value"}},
		{"stringer-ptr", &MyStringer{value: "Stringer Value"}},
		{"duration", 1500 * time.Millisecond},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = ToString(tc.value)
			}
		})
	}
}

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

func Benchmark_UnsafeBytes(b *testing.B) {
	hello := "Hello, World!"
	var res []byte
	b.Run("unsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = UnsafeBytes(hello)
		}
		require.Equal(b, []byte("Hello, World!"), res)
	})
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = []byte(hello)
		}
		require.Equal(b, []byte("Hello, World!"), res)
	})
}

func Benchmark_UnsafeString(b *testing.B) {
	hello := []byte("Hello, World!")
	var res string
	b.Run("unsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = UnsafeString(hello)
		}
		require.Equal(b, "Hello, World!", res)
	})
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = string(hello)
		}
		require.Equal(b, "Hello, World!", res)
	})
}

const (
	uKilobyte = 1 << (10 * (iota + 1))
	uMegabyte
	uGigabyte
	uTerabyte
	uPetabyte
	uExabyte
)

func Benchmark_ByteSize(b *testing.B) {
	testCases := []uint64{
		0,
		1,
		500,
		uKilobyte,
		1126,
		uMegabyte,
		1126 * uKilobyte,
		uGigabyte,
		1126 * uMegabyte,
		uTerabyte,
		1126 * uGigabyte,
		uPetabyte,
		1126 * uTerabyte,
		uExabyte,
		1126 * uPetabyte,
	}
	for _, bytes := range testCases {
		b.Run(strconv.FormatUint(bytes, 10), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ByteSize(bytes)
			}
		})
	}
}

func Test_UnsafeConversions_Parity(t *testing.T) {
	t.Parallel()
	for n := 0; n <= 70; n++ {
		src := make([]byte, n)
		for i := range src {
			src[i] = byte(i*37 + 1)
		}

		view := UnsafeString(src)
		require.Equal(t, string(src), view, n)
		require.Len(t, view, n)
		if n > 0 {
			src[0] ^= 0xff
			require.Equal(t, string(src), view, n)
			require.Same(t, &src[0], &UnsafeBytes(view)[0])
			src[0] ^= 0xff
		}

		str := string(src)
		bs := UnsafeBytes(str)
		require.Equal(t, str, string(bs)) // an empty string may give a nil slice
		require.Len(t, bs, n)
		require.Equal(t, n, cap(bs))

		copied := CopyString(str)
		require.Equal(t, str, copied)
		if n > 1 { // one-byte strings are shared by the runtime
			require.NotSame(t, &UnsafeBytes(str)[0], &UnsafeBytes(copied)[0], "CopyString must not share the original's bytes")
		}

		copiedBytes := CopyBytes(src)
		require.Equal(t, src, copiedBytes)
		require.Len(t, copiedBytes, n)
		if n > 0 {
			require.NotSame(t, &src[0], &copiedBytes[0])
			src[0] ^= 0xff
			require.NotEqual(t, src[0], copiedBytes[0])
		}
	}
}

// ToString is fmt.Sprint with faster paths; these tests pin where they agree and where they differ.

type nilReceiver struct{ name string }

func (n *nilReceiver) String() string { return n.name }

type nilReceiverError struct{ msg string }

func (n *nilReceiverError) Error() string { return n.msg }

type nilSafe struct{}

func (n *nilSafe) String() string {
	if n == nil {
		return "no value"
	}
	return "a value"
}

// nilMap, nilFunc, panickyMap and holder are not pointers, but zero ones have a nil data word in
// an interface, as nil pointers do (a struct with one pointer field is stored as that pointer).
type nilMap map[string]int

func (m nilMap) String() string { return "entries: " + strconv.Itoa(len(m)) }

type nilFunc func()

func (f nilFunc) String() string { return "func is nil: " + strconv.FormatBool(f == nil) }

type panicky struct{}

func (panicky) String() string { panic("boom") }

type panickyMap map[string]int

func (panickyMap) String() string { panic("boom") }

type holder struct{ p *int }

func (h holder) String() string { return strconv.Itoa(*h.p) }

type holderError struct{ p *int }

func (h holderError) Error() string { return strconv.Itoa(*h.p) }

type (
	namedInt    int
	namedString string
	namedFloat  float64
	namedBool   bool
)

func Test_ToString_PointersAndEmptyCollections(t *testing.T) {
	t.Parallel()

	s, i, i64, u64, f64, b := "value", 42, int64(-64), uint64(64), 1.5, true
	i8, u8, i32, f32 := int8(-8), uint8(8), int32(-32), float32(0.1)
	slice := []int{1, 2}
	ptrToInt := &i
	strct := CustomType{Field1: "example", Field2: 42}

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"*string", &s, "value"},
		{"*int", &i, "42"},
		{"*int64", &i64, "-64"},
		{"*uint64", &u64, "64"},
		{"*float64", &f64, "1.5"},
		{"*bool", &b, "true"},
		{"*int8", &i8, "-8"},
		{"*uint8", &u8, "8"},
		{"*int32", &i32, "-32"},
		{"*float32", &f32, "0.1"},
		{"**int", &ptrToInt, "42"},
		{"*[]int", &slice, "[1 2]"},
		{"*struct", &strct, "{example 42}"},
		{"[]string{}", []string{}, "[]"},
		{"[]string(nil)", []string(nil), "[]"},
		{"[]int{}", []int{}, "[]"},
		{"[]int(nil)", []int(nil), "[]"},
		{"[]float64{}", []float64{}, "[]"},
		{"[]any{}", []any{}, "[]"},
		{"[0]int", [0]int{}, "[]"},
		{"[0]string", [0]string{}, "[]"},
		{"[]CustomType{}", []CustomType{}, "[]"},
		{"empty []byte", []byte{}, ""},
		{"nil []byte", []byte(nil), ""},
		{"[]float64", []float64{1.5, 2, 1e21}, "[1.5 2 1000000000000000000000]"},
		{"[]bool", []bool{true, false}, "[true false]"},
		{"[2]string", [2]string{"a", "b"}, "[a b]"},
		{"[]*int", []*int{&i, nil}, "[42 0]"},
		{"[]error", []error{errors.New("e1"), nil}, "[e1 ]"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, ToString(tc.input))
		})
	}
}

// Test_ToString_FmtParity requires ToString to answer as fmt.Sprint does wherever they agree by design.
func Test_ToString_FmtParity(t *testing.T) {
	t.Parallel()

	const maxFloat32Int = 1 << 24
	check := func(v any) {
		t.Helper()
		require.Equal(t, fmt.Sprint(v), ToString(v), "%T %v", v, v)
	}

	for _, n := range []int64{0, 1, -1, 9, 10, 99, 100, 101, -99, -100, math.MaxInt8, math.MinInt8, math.MaxInt16, math.MinInt16, math.MaxInt32, math.MinInt32, math.MaxInt64, math.MinInt64, maxFloat32Int} {
		check(int(n))
		check(n)
		check(int32(n))
		check(int16(n))
		check(int8(n))
		check(uint(n))
		check(uint64(n))
		check(uint32(n))
		check(uint16(n))
		check(uint8(n))
		check(namedInt(n))
		check([]int{int(n), 7})
	}
	for _, n := range []uint64{math.MaxUint64, math.MaxUint32, math.MaxUint16, math.MaxUint8, 1 << 63, 10000000000000000000} {
		check(n)
		check(uint(n))
	}
	for _, v := range []string{"", "text", "multi word", "\x00\xff", "\u00fcn\u00efc\u00f6d\u00e9", "[brackets]"} {
		check(v)
		check(namedString(v))
		check([]string{v, v})
		check([2]string{v, "x"})
		check(MyStringer{value: v})
		check(errors.New(v))
		check(fmt.Errorf("wrapped: %w", errors.New(v)))
	}
	for _, v := range []bool{true, false} {
		check(v)
		check(namedBool(v))
		check([]bool{v, !v})
	}
	check(time.Second)
	check(time.March)
	check(net.IPv4(1, 2, 3, 4))
	check(netip.MustParseAddr("2001:db8::1"))
	check([]any{1, "a", true, 1.5, MyStringer{value: "s"}})
	check([][]int{{1}, {2, 3}, {}})
	check(map[string]int{"a": 1, "b": 2})
	check(struct{ A, B int }{1, 2})
	check(CustomType{Field1: "x", Field2: 2})
	check(complex(1, 2))
	check(uintptr(7))
	check('x')
	check(namedFloat(0.5))
}

// Test_ToString_FloatsAreNeverExponential: plain decimal at the shortest round-trip length, where fmt uses an exponent.
func Test_ToString_FloatsAreNeverExponential(t *testing.T) {
	t.Parallel()
	for _, v := range []float64{0, 1, -1, 0.1, 1.5, 100, 123456789, 1e20, 1e21, 1e-6, 1e-7, 1.7976931348623157e308, 5e-324, math.Pi, math.Inf(1), math.Inf(-1), math.NaN()} {
		want := strconv.FormatFloat(v, 'f', -1, 64)
		require.Equal(t, want, ToString(v), "float64 %v", v)
		require.Equal(t, want, ToString(&v), "*float64 %v", v)
		f32 := float32(v)
		if !math.IsInf(float64(f32), 0) || math.IsInf(v, 0) {
			require.Equal(t, strconv.FormatFloat(float64(f32), 'f', -1, 32), ToString(f32), "float32 %v", f32)
		}
	}
	require.NotEqual(t, fmt.Sprint(1e21), ToString(1e21), "fmt no longer switches to an exponent: the documented difference is gone")
	require.Equal(t, "1000000000000000000000", ToString(1e21))
	require.Equal(t, "0.0000001", ToString(1e-7))
}

func Test_ToString_Differences(t *testing.T) {
	t.Parallel()
	var nilInt *int
	var nilString *string
	moment := time.Date(2001, 2, 3, 4, 5, 6, 789, time.UTC)
	n := 5
	i := 7
	tests := []struct {
		name  string
		input any
		want  string
		fmt   string
	}{
		{"nil", nil, "", "<nil>"},
		{"[]byte is text", []byte("bytes"), "bytes", "[98 121 116 101 115]"},
		{"nil *int is zero", nilInt, "0", "<nil>"},
		{"nil *string is empty", nilString, "", "<nil>"},
		{"pointers are followed", &n, "5", ""},
		{"nil elements are empty", []any{1, nil}, "[1 ]", "[1 <nil>]"},
		{"time.Time has a short layout", moment, "2001-02-03 04:05:06", "2001-02-03 04:05:06.000000789 +0000 UTC"},
		{"pointer to a struct", &CustomType{Field1: "x", Field2: 1}, "{x 1}", "&{x 1}"},
		{"slice of pointers", []*int{&i}, "[7]", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, ToString(tc.input))
			if tc.fmt != "" {
				require.Equal(t, tc.fmt, fmt.Sprint(tc.input), "fmt no longer renders this differently: the documented difference is gone")
			}
		})
	}

	// The layout applies to time.Time wherever reached: directly, via reflect.Value, in slices and arrays.
	require.Equal(t, "2001", ToString(moment, "2006"))
	require.Equal(t, "2001", ToString(reflect.ValueOf(moment), "2006"))
	require.Equal(t, "[2001 2001]", ToString([]time.Time{moment, moment}, "2006"))
	require.Equal(t, "[2001]", ToString([1]time.Time{moment}, "2006"))
	require.Equal(t, "[2001 x]", ToString([]any{moment, "x"}, "2006"))
}

type (
	selfPointer *selfPointer
	pingPointer *pongPointer
	pongPointer *pingPointer
)

// Test_ToString_EndlessPointers: a type like T *T never ends when followed, and once overflowed the
// stack. ToString prints the address, as fmt does.
func Test_ToString_EndlessPointers(t *testing.T) {
	t.Parallel()
	var self selfPointer
	self = &self
	var ping pingPointer
	var pong pongPointer
	ping, pong = &pong, &ping

	for name, v := range map[string]any{"self": self, "ping": ping, "pong": pong} {
		var got string
		require.NotPanics(t, func() { got = ToString(v) }, name)
		require.Equal(t, fmt.Sprint(v), got, name)
	}

	n := 7
	p := &n
	pp := &p
	ppp := &pp
	require.Equal(t, "7", ToString(ppp))

	// No fixed depth cap: only a type that comes back to itself is endless.
	deep := reflect.ValueOf(&n)
	for range 39 {
		next := reflect.New(deep.Type())
		next.Elem().Set(deep)
		deep = next
	}
	require.Equal(t, "7", ToString(deep.Interface()))
}

// Test_ToString_NilReceivers requires a nil error or Stringer pointer to render as fmt does, not panic.
func Test_ToString_NilReceivers(t *testing.T) {
	t.Parallel()
	var nilStringer *nilReceiver
	var nilError *nilReceiverError
	var nilTime *time.Time
	var nilIP *net.IPNet
	var nilSafeValue *nilSafe

	for name, v := range map[string]any{
		"nil Stringer":        nilStringer,
		"nil error":           nilError,
		"nil *time.Time":      nilTime,
		"nil *net.IPNet":      nilIP,
		"nil-safe Stringer":   nilSafeValue,
		"nil map":             nilMap(nil),
		"nil func":            nilFunc(nil),
		"in a slice":          []any{nilStringer, nilError, nilSafeValue},
		"slice of nil errors": []error{nilError},
		"in an array":         [2]fmt.Stringer{nilStringer, nilSafeValue},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got string
			require.NotPanics(t, func() { got = ToString(v) })
			require.Equal(t, fmt.Sprint(v), got)
		})
	}

	// The cases that tell a guarded call from a blanket "<nil>".
	require.Equal(t, "no value", ToString(nilSafeValue))
	require.Equal(t, "entries: 0", ToString(nilMap(nil)))
	require.Equal(t, "func is nil: true", ToString(nilFunc(nil)))
	require.Equal(t, "<nil>", ToString(nilStringer))
	require.Equal(t, "<nil>", ToString(nilError))

	require.Equal(t, "named", ToString(&nilReceiver{name: "named"}))
	require.Equal(t, "broken", ToString(&nilReceiverError{msg: "broken"}))
	require.Equal(t, "a value", ToString(&nilSafe{}))
	require.Equal(t, "entries: 1", ToString(nilMap{"k": 1}))

	// A panic that is not a nil receiver's is not hidden.
	require.PanicsWithValue(t, "boom", func() { ToString(panicky{}) })
	require.PanicsWithValue(t, "boom", func() { ToString(&panicky{}) })

	// Only a nil pointer gets "<nil>" for a panic; a nil map or a zero holder panics like any other value.
	const nilDereference = "runtime error: invalid memory address or nil pointer dereference"
	require.PanicsWithValue(t, "boom", func() { ToString(panickyMap(nil)) })
	require.PanicsWithError(t, nilDereference, func() { ToString(holder{}) })
	require.PanicsWithError(t, nilDereference, func() { ToString(holderError{}) })
	seven := 7
	require.Equal(t, "7", ToString(holder{p: &seven}))
	require.Equal(t, "7", ToString(holderError{p: &seven}))
}
