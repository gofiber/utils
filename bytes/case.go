package bytes

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

// ToLower converts an ASCII byte slice to lower-case without modifying the input.
func ToLower(b []byte) []byte {
	n := len(b)
	if n == 0 {
		return b
	}

	table := caseconv.ToLowerTable
	i := 0
	for i < n {
		c := b[i]
		low := table[c]
		if low != c {
			dst := make([]byte, n)
			copy(dst, b[:i])
			dst[i] = low
			i++
			limit := i + ((n - i) &^ 3)
			for i < limit {
				dst[i+0] = table[b[i+0]]
				dst[i+1] = table[b[i+1]]
				dst[i+2] = table[b[i+2]]
				dst[i+3] = table[b[i+3]]
				i += 4
			}
			for i < n {
				dst[i] = table[b[i]]
				i++
			}
			return dst
		}
		i++
	}
	return b
}

// ToUpper converts an ASCII byte slice to upper-case without modifying the input.
func ToUpper(b []byte) []byte {
	n := len(b)
	if n == 0 {
		return b
	}

	table := caseconv.ToUpperTable
	i := 0
	for i < n {
		c := b[i]
		up := table[c]
		if up != c {
			dst := make([]byte, n)
			copy(dst, b[:i])
			dst[i] = up
			i++
			limit := i + ((n - i) &^ 3)
			for i < limit {
				dst[i+0] = table[b[i+0]]
				dst[i+1] = table[b[i+1]]
				dst[i+2] = table[b[i+2]]
				dst[i+3] = table[b[i+3]]
				i += 4
			}
			for i < n {
				dst[i] = table[b[i]]
				i++
			}
			return dst
		}
		i++
	}
	return b
}

// UnsafeToLower converts an ASCII byte slice to lower-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToLower(b []byte) []byte {
	table := caseconv.ToLowerTable
	n := len(b)
	i := 0
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		b1 := b[i+1]
		b2 := b[i+2]
		b3 := b[i+3]

		b[i+0] = table[b0]
		b[i+1] = table[b1]
		b[i+2] = table[b2]
		b[i+3] = table[b3]

		i += 4
	}
	for i < n {
		b[i] = table[b[i]]
		i++
	}
	return b
}

// UnsafeToUpper converts an ASCII byte slice to upper-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToUpper(b []byte) []byte {
	table := caseconv.ToUpperTable
	n := len(b)
	i := 0
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		b1 := b[i+1]
		b2 := b[i+2]
		b3 := b[i+3]

		b[i+0] = table[b0]
		b[i+1] = table[b1]
		b[i+2] = table[b2]
		b[i+3] = table[b3]

		i += 4
	}
	for i < n {
		b[i] = table[b[i]]
		i++
	}
	return b
}
