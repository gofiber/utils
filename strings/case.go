package strings

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// ToLower converts an ASCII string to lower-case.
func ToLower(s string) string {
	n := len(s)
	if n < caseconv.WordLen {
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

	// caseconv only reads src; its fresh result is never written again.
	src := unsafeconv.UnsafeBytes(s)
	i := caseconv.FirstUpperIndex(src)
	if i < 0 {
		return s
	}
	return unsafeconv.UnsafeString(caseconv.ToLowerFrom(src, i))
}

// ToUpper converts an ASCII string to upper-case.
func ToUpper(s string) string {
	n := len(s)
	if n < caseconv.WordLen {
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
	return unsafeconv.UnsafeString(caseconv.ToUpperFrom(src, i))
}

// UnsafeToLower lower-cases the ASCII letters of s in place; s must reference mutable memory.
func UnsafeToLower(s string) string {
	b := unsafeconv.UnsafeBytes(s)
	if len(b) < caseconv.WordLen {
		table := caseconv.ToLowerTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return s
	}
	caseconv.ToLowerInPlace(b)
	return s
}

// UnsafeToUpper upper-cases the ASCII letters of s in place; s must reference mutable memory.
func UnsafeToUpper(s string) string {
	b := unsafeconv.UnsafeBytes(s)
	if len(b) < caseconv.WordLen {
		table := caseconv.ToUpperTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return s
	}
	caseconv.ToUpperInPlace(b)
	return s
}
