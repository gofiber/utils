package bytes

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

// Sub-word inputs skip caseconv: its call costs more than the table lookups. Delegating them
// measured +47%/+43% on unchanged three-byte input and +65%/+60% for the in-place forms
// (benchstat n=15, Go 1.26, linux/amd64).

// ToLower returns an ASCII lower-cased copy of b, or b itself if it has no upper-case letter: a
// write to the result then writes to b, unlike with bytes.ToLower, which always copies.
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

// ToUpper returns an ASCII upper-cased copy of b, or b itself if it has no lower-case letter: a
// write to the result then writes to b, unlike with bytes.ToUpper, which always copies.
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
