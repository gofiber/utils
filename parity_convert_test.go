package utils

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/netip"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ToString is fmt.Sprint with a faster road for the types a web server meets
// most. These tests pin where the two agree, which is nearly everywhere, and
// the few places the difference is the point.

// nilReceiver implements the interfaces ToString dispatches on with methods
// that do not guard against a nil pointer, as most do not.
type nilReceiver struct{ name string }

func (n *nilReceiver) String() string { return n.name }

type nilReceiverError struct{ msg string }

func (n *nilReceiverError) Error() string { return n.msg }

// nilSafe is a Stringer written to be called on nil, which fmt honors.
type nilSafe struct{}

func (n *nilSafe) String() string {
	if n == nil {
		return "no value"
	}
	return "a value"
}

// nilMap has a method that works on a nil map, one of the other nil values
// whose interface data word is nil.
type nilMap map[string]int

func (m nilMap) String() string { return "entries: " + strconv.Itoa(len(m)) }

// nilFunc is the same for a func.
type nilFunc func()

func (f nilFunc) String() string { return "func is nil: " + strconv.FormatBool(f == nil) }

// panicky panics in its method whatever its receiver.
type panicky struct{}

func (panicky) String() string { panic("boom") }

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
		// The pointer types with a road of their own.
		{"*string", &s, "value"},
		{"*int", &i, "42"},
		{"*int64", &i64, "-64"},
		{"*uint64", &u64, "64"},
		{"*float64", &f64, "1.5"},
		{"*bool", &b, "true"},
		// Every other pointer is followed by reflection.
		{"*int8", &i8, "-8"},
		{"*uint8", &u8, "8"},
		{"*int32", &i32, "-32"},
		{"*float32", &f32, "0.1"},
		{"**int", &ptrToInt, "42"},
		{"*[]int", &slice, "[1 2]"},
		{"*struct", &strct, "{example 42}"},
		// Empty collections of every kind that has a road.
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
		// Collections of the less common kinds.
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

// Test_ToString_FmtParity requires ToString to answer as fmt.Sprint does for
// everything whose rendering is the same by design: integers of every width
// and sign, strings, booleans, errors, Stringers, and collections of them.
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
	check(namedFloat(0.5)) // a named type with no road of its own is rendered by fmt
}

// Test_ToString_FloatsAreNeverExponential pins the first deliberate
// difference: floats are written in plain decimal at their shortest round-trip
// length, where fmt switches to an exponent for very large and small values.
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

// Test_ToString_Differences pins the rest of the deliberate differences.
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
		fmt   string // what fmt.Sprint gives, to prove the difference is real
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

	// The layout argument is only for time.Time, wherever it is reached.
	require.Equal(t, "2001", ToString(moment, "2006"))
	require.Equal(t, "2001", ToString(reflect.ValueOf(moment), "2006"))
}

// Test_ToString_NilReceivers requires a nil pointer whose type implements
// error or fmt.Stringer to be rendered as fmt renders it, where calling its
// method would panic: "<nil>", or what a method written for nil returns.
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

	// A pointer that is not nil still gets its method.
	require.Equal(t, "named", ToString(&nilReceiver{name: "named"}))
	require.Equal(t, "broken", ToString(&nilReceiverError{msg: "broken"}))
	require.Equal(t, "a value", ToString(&nilSafe{}))
	require.Equal(t, "entries: 1", ToString(nilMap{"k": 1}))

	// A panic that a nil receiver did not cause is not hidden: fmt reports such a
	// panic in its output, and ToString lets it through as it always has.
	require.PanicsWithValue(t, "boom", func() { ToString(panicky{}) })
	require.PanicsWithValue(t, "boom", func() { ToString(&panicky{}) })
}
