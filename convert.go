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
)

// UnsafeString returns a string pointer without allocation
func UnsafeString(b []byte) string {
	// the new way is slower `return unsafe.String(unsafe.SliceData(b), len(b))`
	// unsafe.Pointer variant: 0.3538 ns/op vs unsafe.String variant: 0.5410 ns/op
	// #nosec G103
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes returns a byte pointer without allocation.
func UnsafeBytes(s string) []byte {
	// #nosec G103
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// CopyString copies a string to make it immutable
func CopyString(s string) string {
	// #nosec G103
	return string(UnsafeBytes(s))
}

// #nosec G103
// CopyBytes copies a slice to make it immutable
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

// ToString Change arg to string
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
		return ToString(v.Interface(), timeFormat...)
	case fmt.Stringer:
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
		var buf strings.Builder
		buf.Grow(len(v) * 8) // Pre-allocate approximate size
		buf.WriteByte('[')
		for i, s := range v {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(s)
		}
		buf.WriteByte(']')
		return buf.String()
	case []int:
		if len(v) == 0 {
			return "[]"
		}
		var buf strings.Builder
		buf.Grow(len(v) * 4) // Pre-allocate approximate size
		buf.WriteByte('[')
		for i, n := range v {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(FormatInt(int64(n)))
		}
		buf.WriteByte(']')
		return buf.String()
	default:
		// Check if the type is a pointer by using reflection
		rv := reflect.ValueOf(arg)
		kind := rv.Kind()
		if kind == reflect.Ptr && !rv.IsNil() {
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
				buf.WriteString(ToString(rv.Index(i).Interface()))
			}
			buf.WriteByte(']')
			return buf.String()
		}

		// For types not explicitly handled, use fmt.Sprint to generate a string representation
		return fmt.Sprint(arg)
	}
}
