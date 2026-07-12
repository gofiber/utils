package strings

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// swarMinLen is the smallest input length worth routing through the
// word-at-a-time (SWAR) helpers; shorter inputs are cheaper byte-by-byte.
const swarMinLen = caseconv.WordLen

// ToLower converts an ASCII string to lower-case without modifying the input.
func ToLower(s string) string {
	n := len(s)
	if n < swarMinLen {
		table := caseconv.ToLowerTable
		for i := 0; i < n; i++ {
			c := s[i]
			low := table[c]
			if low != c {
				res := make([]byte, n)
				copy(res, s[:i])
				res[i] = low
				for i++; i < n; i++ {
					res[i] = table[s[i]]
				}
				return unsafeconv.UnsafeString(res)
			}
		}
		return s
	}

	src := unsafeconv.UnsafeBytes(s)
	i := caseconv.FirstUpperIndex(src)
	if i < 0 {
		return s
	}

	res := make([]byte, n)
	// Copy the unchanged prefix up to the word containing the first
	// uppercase byte, then convert the rest word-at-a-time.
	from := i &^ (caseconv.WordLen - 1)
	copy(res, src[:from])
	caseconv.ToLowerCopy(res, src, from)
	return unsafeconv.UnsafeString(res)
}

// ToUpper converts an ASCII string to upper-case without modifying the input.
func ToUpper(s string) string {
	n := len(s)
	if n < swarMinLen {
		table := caseconv.ToUpperTable
		for i := 0; i < n; i++ {
			c := s[i]
			up := table[c]
			if up != c {
				res := make([]byte, n)
				copy(res, s[:i])
				res[i] = up
				for i++; i < n; i++ {
					res[i] = table[s[i]]
				}
				return unsafeconv.UnsafeString(res)
			}
		}
		return s
	}

	src := unsafeconv.UnsafeBytes(s)
	i := caseconv.FirstLowerIndex(src)
	if i < 0 {
		return s
	}

	res := make([]byte, n)
	from := i &^ (caseconv.WordLen - 1)
	copy(res, src[:from])
	caseconv.ToUpperCopy(res, src, from)
	return unsafeconv.UnsafeString(res)
}

// UnsafeToLower converts an ASCII string to lower-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func UnsafeToLower(s string) string {
	b := unsafeconv.UnsafeBytes(s)
	if len(b) < swarMinLen {
		table := caseconv.ToLowerTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return s
	}
	caseconv.ToLowerInPlace(b)
	return s
}

// UnsafeToUpper converts an ASCII string to upper-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func UnsafeToUpper(s string) string {
	b := unsafeconv.UnsafeBytes(s)
	if len(b) < swarMinLen {
		table := caseconv.ToUpperTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return s
	}
	caseconv.ToUpperInPlace(b)
	return s
}
