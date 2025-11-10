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
		buf = strconv.AppendUint(buf, bytes, 10)
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

	buf = strconv.AppendUint(buf, integer, 10)
	if fractional > 0 {
		buf = append(buf, '.')
		buf = strconv.AppendUint(buf, fractional, 10)
	}
	buf = append(buf, unit...)
	return UnsafeString(buf)
}

// ToString Change arg to string
func ToString(arg any, timeFormat ...string) string {
	switch v := arg.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
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
	default:
		// Check if the type is a pointer by using reflection
		rv := reflect.ValueOf(arg)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			// Dereference the pointer and recursively call ToString
			return ToString(rv.Elem().Interface(), timeFormat...)
		} else if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			// handle slices
			var buf strings.Builder
			buf.WriteString("[") //nolint:errcheck // no need to check error
			for i := 0; i < rv.Len(); i++ {
				if i > 0 {
					buf.WriteString(" ") //nolint:errcheck // no need to check error
				}
				buf.WriteString(ToString(rv.Index(i).Interface())) //nolint:errcheck // no need to check error
			}
			buf.WriteString("]") //nolint:errcheck // no need to check error
			return buf.String()
		}

		// For types not explicitly handled, use fmt.Sprint to generate a string representation
		return fmt.Sprint(arg)
	}
}
