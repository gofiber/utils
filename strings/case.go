package strings

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// ToLower converts an ASCII string to lower-case without modifying the input.
func ToLower(s string) string {
	n := len(s)
	if n == 0 {
		return s
	}

	table := caseconv.ToLowerTable
	for i := range n {
		c := s[i]
		low := table[c]
		if low != c {
			res := make([]byte, n)
			copy(res, s[:i])
			res[i] = low
			j := i + 1
			for ; j+3 < n; j += 4 {
				res[j+0] = table[s[j+0]]
				res[j+1] = table[s[j+1]]
				res[j+2] = table[s[j+2]]
				res[j+3] = table[s[j+3]]
			}
			for ; j < n; j++ {
				res[j] = table[s[j]]
			}
			return unsafeconv.UnsafeString(res)
		}
	}
	return s
}

// ToUpper converts an ASCII string to upper-case without modifying the input.
func ToUpper(s string) string {
	n := len(s)
	if n == 0 {
		return s
	}

	table := caseconv.ToUpperTable
	for i := range n {
		c := s[i]
		up := table[c]
		if up != c {
			res := make([]byte, n)
			copy(res, s[:i])
			res[i] = up
			j := i + 1
			for ; j+3 < n; j += 4 {
				res[j+0] = table[s[j+0]]
				res[j+1] = table[s[j+1]]
				res[j+2] = table[s[j+2]]
				res[j+3] = table[s[j+3]]
			}
			for ; j < n; j++ {
				res[j] = table[s[j]]
			}
			return unsafeconv.UnsafeString(res)
		}
	}
	return s
}

// UnsafeToLower converts an ASCII string to lower-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func UnsafeToLower(s string) string {
	b := unsafeconv.UnsafeBytes(s)
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
	return s
}

// UnsafeToUpper converts an ASCII string to upper-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func UnsafeToUpper(s string) string {
	b := unsafeconv.UnsafeBytes(s)
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
	return s
}
