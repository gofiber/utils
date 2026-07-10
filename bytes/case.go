package bytes

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

// swarMinLen is the smallest input length worth routing through the
// word-at-a-time (SWAR) helpers; shorter inputs are cheaper byte-by-byte.
const swarMinLen = caseconv.WordLen

// ToLower converts an ASCII byte slice to lower-case without modifying the input.
func ToLower(b []byte) []byte {
	n := len(b)
	if n < swarMinLen {
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

	dst := make([]byte, n)
	// Copy the unchanged prefix up to the word containing the first
	// uppercase byte, then convert the rest word-at-a-time.
	from := i &^ (caseconv.WordLen - 1)
	copy(dst, b[:from])
	caseconv.ToLowerCopy(dst, b, from)
	return dst
}

// ToUpper converts an ASCII byte slice to upper-case without modifying the input.
func ToUpper(b []byte) []byte {
	n := len(b)
	if n < swarMinLen {
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

	dst := make([]byte, n)
	from := i &^ (caseconv.WordLen - 1)
	copy(dst, b[:from])
	caseconv.ToUpperCopy(dst, b, from)
	return dst
}

// UnsafeToLower converts an ASCII byte slice to lower-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToLower(b []byte) []byte {
	if len(b) < swarMinLen {
		table := caseconv.ToLowerTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToLowerInPlace(b)
	return b
}

// UnsafeToUpper converts an ASCII byte slice to upper-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToUpper(b []byte) []byte {
	if len(b) < swarMinLen {
		table := caseconv.ToUpperTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToUpperInPlace(b)
	return b
}
