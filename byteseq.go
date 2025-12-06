package utils

import (
	"reflect"
	"slices"
)

type byteSeq interface {
	~string | ~[]byte
}

// EqualFold tests ascii strings or bytes for equality case-insensitively
func EqualFold[S byteSeq](b, s S) bool {
	if len(b) != len(s) {
		return false
	}

	table := toUpperTable
	n := len(b)
	i := 0

	// Unroll by 4 to match other hot paths and drive instruction-level parallelism.
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		s0 := s[i+0]
		if table[b0] != table[s0] {
			return false
		}

		b1 := b[i+1]
		s1 := s[i+1]
		if table[b1] != table[s1] {
			return false
		}

		b2 := b[i+2]
		s2 := s[i+2]
		if table[b2] != table[s2] {
			return false
		}

		b3 := b[i+3]
		s3 := s[i+3]
		if table[b3] != table[s3] {
			return false
		}

		i += 4
	}

	for i < n {
		if table[b[i]] != table[s[i]] {
			return false
		}
		i++
	}
	return true
}

// TrimLeft is the equivalent of strings/bytes.TrimLeft
func TrimLeft[S byteSeq](s S, cutset byte) S {
	lenStr, start := len(s), 0
	for start < lenStr && s[start] == cutset {
		start++
	}
	return s[start:]
}

// Trim is the equivalent of strings/bytes.Trim
func Trim[S byteSeq](s S, cutset byte) S {
	i, j := 0, len(s)-1
	for ; i <= j; i++ {
		if s[i] != cutset {
			break
		}
	}
	for ; i < j; j-- {
		if s[j] != cutset {
			break
		}
	}

	return s[i : j+1]
}

// TrimRight is the equivalent of strings/bytes.TrimRight
func TrimRight[S byteSeq](s S, cutset byte) S {
	lenStr := len(s)
	for lenStr > 0 && s[lenStr-1] == cutset {
		lenStr--
	}
	return s[:lenStr]
}

// TrimSpace removes leading and trailing whitespace from a string or byte slice.
// This is an optimized version that's faster than strings/bytes.TrimSpace for ASCII strings.
// It removes the following ASCII whitespace characters: space, tab, newline, carriage return, vertical tab, and form feed.
func TrimSpace[S byteSeq](s S) S {
	i, j := 0, len(s)-1

	// fast path for empty string
	if j < 0 {
		return s
	}

	// Find first non-whitespace from start
	for ; i <= j && whitespaceTable[s[i]]; i++ { //nolint:revive // we want to check for multiple whitespace chars
	}

	// Find first non-whitespace from end
	for ; i < j && whitespaceTable[s[j]]; j-- { //nolint:revive // we want to check for multiple whitespace chars
	}

	return s[i : j+1]
}

// AddTrailingSlash appends a trailing '/' to v if it does not already end with one.
//
// For string inputs, a new string is returned only when a '/' needs to be appended.
// If the input already ends with '/', the original string is returned.
//
// For []byte inputs, a new slice is returned when a '/' is appended.
// If the input already ends with '/', the original slice is returned unchanged.
// The original slice is never modified.
func AddTrailingSlash[T byteSeq](v T) T {
	n := len(v)
	if n > 0 && v[n-1] == '/' {
		return v
	}
	// Type-specific optimization
	switch x := any(v).(type) {
	case string:
		// For strings: allocate exact size, use UnsafeString to avoid double alloc
		buf := make([]byte, n+1)
		copy(buf, x)
		buf[n] = '/'
		return any(UnsafeString(buf)).(T) //nolint:forcetypeassert,errcheck // type is guaranteed
	case []byte:
		// For bytes: use append which may reuse capacity
		return any(append(x, '/')).(T) //nolint:forcetypeassert,errcheck // type is guaranteed
	default:
		// Fallback for named types (e.g., type MyString string) using reflection
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.String {
			s := val.String() + "/"
			return reflect.ValueOf(s).Convert(val.Type()).Interface().(T) //nolint:forcetypeassert,errcheck // type is guaranteed
		}
		// Assumed to be a slice of bytes
		b := val.Bytes()
		res := append(slices.Clone(b), '/')
		newSlice := reflect.MakeSlice(val.Type(), len(res), len(res))
		reflect.Copy(newSlice, reflect.ValueOf(res))
		return newSlice.Interface().(T) //nolint:forcetypeassert,errcheck // type is guaranteed
	}
}
