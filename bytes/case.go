package bytes

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

// Sub-word inputs skip caseconv: its call costs more than the table lookups.

// ToLower returns an ASCII lower-cased copy of b, or b itself if it has no upper-case letter.
func ToLower(b []byte) []byte {
	n := len(b)
	if n < caseconv.WordLen {
		table := caseconv.ToLowerTable
		for i := 0; i < n; i++ {
			c := b[i]
			low := table[c]
			if low != c {
				dst := make([]byte, n)
				copy(dst, b[:i])
				dst[i] = low
				for i++; i < n; i++ {
					dst[i] = table[b[i]]
				}
				return dst
			}
		}
		return b
	}

	i := caseconv.FirstUpperIndex(b)
	if i < 0 {
		return b
	}
	return caseconv.ToLowerFrom(b, i)
}

// ToUpper returns an ASCII upper-cased copy of b, or b itself if it has no lower-case letter.
func ToUpper(b []byte) []byte {
	n := len(b)
	if n < caseconv.WordLen {
		table := caseconv.ToUpperTable
		for i := 0; i < n; i++ {
			c := b[i]
			up := table[c]
			if up != c {
				dst := make([]byte, n)
				copy(dst, b[:i])
				dst[i] = up
				for i++; i < n; i++ {
					dst[i] = table[b[i]]
				}
				return dst
			}
		}
		return b
	}

	i := caseconv.FirstLowerIndex(b)
	if i < 0 {
		return b
	}
	return caseconv.ToUpperFrom(b, i)
}

// UnsafeToLower lower-cases b's ASCII letters in place and returns b.
func UnsafeToLower(b []byte) []byte {
	if len(b) < caseconv.WordLen {
		table := caseconv.ToLowerTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToLowerInPlace(b)
	return b
}

// UnsafeToUpper upper-cases b's ASCII letters in place and returns b.
func UnsafeToUpper(b []byte) []byte {
	if len(b) < caseconv.WordLen {
		table := caseconv.ToUpperTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToUpperInPlace(b)
	return b
}
