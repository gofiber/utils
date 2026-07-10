package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

type byteSeq interface {
	~string | ~[]byte
}

// loadWord assembles 8 bytes of s starting at i into a little-endian uint64.
// The 8-byte reslice pins the length so the constant-index reads below are
// bounds-check free, and the compiler fuses them into a single 8-byte load on
// little-endian platforms; on big-endian ones it stays correct, just slower.
func loadWord[S byteSeq](s S, i int) uint64 {
	w := s[i : i+8]
	return uint64(w[0]) |
		uint64(w[1])<<8 |
		uint64(w[2])<<16 |
		uint64(w[3])<<24 |
		uint64(w[4])<<32 |
		uint64(w[5])<<40 |
		uint64(w[6])<<48 |
		uint64(w[7])<<56
}

// EqualFold tests ascii strings or bytes for equality case-insensitively
func EqualFold[S byteSeq](b, s S) bool {
	n := len(b)
	if n != len(s) {
		return false
	}

	// Compare 8 bytes per iteration; case-fold with SWAR only when the raw
	// words differ, so byte-identical input skips both folds entirely.
	i := 0
	for ; i+8 <= n; i += 8 {
		x := loadWord(b, i)
		y := loadWord(s, i)
		if x != y && caseconv.ToUpperWord(x) != caseconv.ToUpperWord(y) {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		// Handle the tail with one overlapping word compare; re-checking
		// bytes that were already equal cannot change the outcome.
		x := loadWord(b, n-8)
		y := loadWord(s, n-8)
		return x == y || caseconv.ToUpperWord(x) == caseconv.ToUpperWord(y)
	}

	table := caseconv.ToUpperTable
	for ; i < n; i++ {
		if table[b[i]] != table[s[i]] {
			return false
		}
	}
	return true
}

// TrimLeft removes all leading occurrences of the byte cutset from s.
// Unlike strings/bytes.TrimLeft, cutset is a single byte, not a set of characters.
func TrimLeft[S byteSeq](s S, cutset byte) S {
	lenStr, start := len(s), 0
	for start < lenStr && s[start] == cutset {
		start++
	}
	return s[start:]
}

// Trim removes all leading and trailing occurrences of the byte cutset from s.
// Unlike strings/bytes.Trim, cutset is a single byte, not a set of characters.
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

// TrimRight removes all trailing occurrences of the byte cutset from s.
// Unlike strings/bytes.TrimRight, cutset is a single byte, not a set of characters.
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
	n := len(s)
	if n == 0 {
		return s
	}

	i, j := 0, n-1
	if !whitespaceTable[s[i]] && !whitespaceTable[s[j]] {
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
