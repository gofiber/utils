// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// UnsafeString returns a string pointer without allocation
func UnsafeString(b []byte) string {
	return unsafeconv.UnsafeString(b)
}

// UnsafeBytes returns a byte pointer without allocation.
// The slice shares the memory of s and must only be read: a string is immutable, and one
// that is a literal lives in read-only memory, where a write is a fatal fault.
func UnsafeBytes(s string) []byte {
	return unsafeconv.UnsafeBytes(s)
}

// CopyString returns a copy of s in memory of its own. The copy stays as it was when s is a
// view of a buffer that is reused or modified later, such as one made with UnsafeString.
func CopyString(s string) string {
	// #nosec G103
	return string(UnsafeBytes(s))
}

// #nosec G103
// CopyBytes returns a copy of b in memory of its own, with capacity equal to its length:
// writes to either slice no longer reach the other. The copy of a nil or empty b is a
// non-nil empty slice.
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

const (
	uByte = 1 << (10 * iota)
	uKilobyte
	uMegabyte
	uGigabyte
	uTerabyte
	uPetabyte
	uExabyte
)

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
// Maximum supported input is math.MaxUint64 / 10 (≈ 1844674407370955161).
func ByteSize(bytes uint64) string {
	const maxSafe = math.MaxUint64 / 10
	unit := ""
	div := uint64(1)
	switch {
	case bytes >= uExabyte:
		unit = "EB"
		div = uExabyte
	case bytes >= uPetabyte:
		unit = "PB"
		div = uPetabyte
	case bytes >= uTerabyte:
		unit = "TB"
		div = uTerabyte
	case bytes >= uGigabyte:
		unit = "GB"
		div = uGigabyte
	case bytes >= uMegabyte:
		unit = "MB"
		div = uMegabyte
	case bytes >= uKilobyte:
		unit = "KB"
		div = uKilobyte
	case bytes >= uByte:
		unit = "B"
	default:
		return "0B"
	}

	buf := make([]byte, 0, 16)

	if div == 1 {
		buf = AppendUint(buf, bytes)
		buf = append(buf, unit...)
		return UnsafeString(buf)
	}

	// Fix: cap bytes to maxSafe for overflow, but format as fractional
	if bytes > maxSafe {
		bytes = maxSafe
	}

	scaled := (bytes/div)*10 + ((bytes%div)*10+div/2)/div
	integer := scaled / 10
	fractional := scaled % 10

	buf = AppendUint(buf, integer)
	if fractional > 0 {
		buf = append(buf, '.')
		buf = AppendUint(buf, fractional)
	}
	buf = append(buf, unit...)
	return UnsafeString(buf)
}

// ToString Change arg to string.
// A time.Time is formatted with timeFormat[0], or "2006-01-02 15:04:05" without one;
// so is each time.Time that is an element of a slice or array.
func ToString(arg any, timeFormat ...string) string {
	switch v := arg.(type) {
	case int:
		return FormatInt(int64(v))
	case int8:
		return FormatInt8(v)
	case int16:
		return FormatInt16(v)
	case int32:
		return FormatInt32(v)
	case int64:
		return FormatInt(v)
	case uint:
		return FormatUint(uint64(v))
	case uint8:
		return FormatUint8(v)
	case uint16:
		return FormatUint16(v)
	case uint32:
		return FormatUint32(v)
	case uint64:
		return FormatUint(v)
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeFormat) > 0 {
			return v.Format(timeFormat[0])
		}
		return v.Format("2006-01-02 15:04:05")
	case reflect.Value:
		if !v.IsValid() || !v.CanInterface() {
			return ""
		}
		return ToString(v.Interface(), timeFormat...)
	// error before fmt.Stringer to match the fmt package's precedence
	case error:
		if isNilPointer(arg) {
			return nilReceiverString(arg, v.Error)
		}
		return v.Error()
	case fmt.Stringer:
		if isNilPointer(arg) {
			return nilReceiverString(arg, v.String)
		}
		return v.String()
	// Handle common pointer types directly to avoid reflection
	case *string:
		if v != nil {
			return *v
		}
		return ""
	case *int:
		if v != nil {
			return FormatInt(int64(*v))
		}
		return "0"
	case *int64:
		if v != nil {
			return FormatInt(*v)
		}
		return "0"
	case *uint64:
		if v != nil {
			return FormatUint(*v)
		}
		return "0"
	case *float64:
		if v != nil {
			return strconv.FormatFloat(*v, 'f', -1, 64)
		}
		return "0"
	case *bool:
		if v != nil {
			return strconv.FormatBool(*v)
		}
		return "false"
	// Handle common slice types directly to avoid reflection
	case []string:
		if len(v) == 0 {
			return "[]"
		}
		// Compute the exact output size for a single allocation:
		// '[' + ']' + one separating space per gap + all element bytes.
		size := len(v) + 1
		for _, s := range v {
			size += len(s)
		}
		buf := make([]byte, 0, size)
		buf = append(buf, '[')
		for i, s := range v {
			if i > 0 {
				buf = append(buf, ' ')
			}
			buf = append(buf, s...)
		}
		buf = append(buf, ']')
		return UnsafeString(buf)
	case []int:
		if len(v) == 0 {
			return "[]"
		}
		// Compute the exact output size for a single allocation:
		// '[' + ']' + one separating space per gap + all element digits.
		size := len(v) + 1
		for _, n := range v {
			size += intDigits(int64(n))
		}
		buf := make([]byte, 0, size)
		buf = append(buf, '[')
		for i, n := range v {
			if i > 0 {
				buf = append(buf, ' ')
			}
			buf = AppendInt(buf, int64(n))
		}
		buf = append(buf, ']')
		return UnsafeString(buf)
	default:
		// Check if the type is a pointer by using reflection
		rv := reflect.ValueOf(arg)
		if !rv.IsValid() {
			return ""
		}
		kind := rv.Kind()
		if kind == reflect.Pointer && !rv.IsNil() {
			// Dereference the pointer and recursively call ToString
			return ToString(rv.Elem().Interface(), timeFormat...)
		} else if kind == reflect.Slice || kind == reflect.Array {
			// handle slices
			n := rv.Len()
			if n == 0 {
				return "[]"
			}
			var buf strings.Builder
			buf.Grow(n * 8) // Pre-allocate approximate size
			buf.WriteByte('[')
			for i := range n {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(ToString(rv.Index(i).Interface(), timeFormat...))
			}
			buf.WriteByte(']')
			return buf.String()
		}

		// For types not explicitly handled, use fmt.Sprint to generate a string representation
		return fmt.Sprint(arg)
	}
}

// isNilPointer reports whether v may hold a nil pointer, by the data word of the
// interface: for a pointer, and for the other pointer-shaped types, that word
// is the value itself. Such a value can still be an error or a fmt.Stringer,
// because the method set of *T includes the methods of T, and calling those on
// it dereferences the nil. The answer is only "maybe": a nil map, func or chan
// and a struct holding a single nil pointer have a nil data word too. A true
// answer sends the call through nilReceiverString, which tells them apart, and
// only when the method panics. Asking reflect instead costs about 1.5ns on
// every error and Stringer.
func isNilPointer(v any) bool {
	// #nosec G103
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1] == nil
}

// nilReceiverString calls the Error or String method of arg, whose data word is
// nil. It returns what the method returns, which a method written to be called
// on nil may well do. If the method panics, a nil pointer gets "<nil>", as in
// fmt, most methods being unguarded; any other value gets its panic back,
// which is not a nil receiver's doing and which ToString lets through for every
// other value.
func nilReceiverString(arg any, method func() string) (s string) { //nolint:nonamedreturns // the deferred recover sets the result
	defer func() {
		if r := recover(); r != nil {
			if v := reflect.ValueOf(arg); v.Kind() != reflect.Pointer || !v.IsNil() {
				panic(r)
			}
			s = "<nil>"
		}
	}()
	return method()
}
