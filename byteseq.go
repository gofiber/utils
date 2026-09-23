package utils

import (
	"math/bits"

	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/swar"
)

type byteSeq interface {
	~string | ~[]byte
}

// load4 assembles s[i:i+4] into a little-endian uint32; the caller guarantees i+4 <= len(s).
func load4[S byteSeq](s S, i int) uint32 {
	w := s[i : i+4]
	return uint32(w[0]) | uint32(w[1])<<8 | uint32(w[2])<<16 | uint32(w[3])<<24
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
		x := swar.Load8(b, i)
		y := swar.Load8(s, i)
		if x != y && swar.ToUpperWord(x) != swar.ToUpperWord(y) {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		// Handle the tail with one overlapping word compare; re-checking
		// bytes that were already equal cannot change the outcome.
		x := swar.Load8(b, n-8)
		y := swar.Load8(s, n-8)
		return x == y || swar.ToUpperWord(x) == swar.ToUpperWord(y)
	}
	if n >= 4 {
		// 4..7 bytes: two overlapping 4-byte windows packed into one word
		// fold with a single ToUpperWord; checked after the word paths so
		// longer inputs pay nothing for it.
		x := uint64(load4(b, 0)) | uint64(load4(b, n-4))<<32
		y := uint64(load4(s, 0)) | uint64(load4(s, n-4))<<32
		return x == y || swar.ToUpperWord(x) == swar.ToUpperWord(y)
	}

	table := caseconv.ToUpperTable
	for ; i < n; i++ {
		if table[b[i]] != table[s[i]] {
			return false
		}
	}
	return true
}

// HashFold returns a 64-bit hash of s in which ASCII letters are folded to one
// case: strings that EqualFold reports equal hash equal. It is the hashing
// half of a case-insensitive lookup table whose other half is EqualFold, for
// keys such as header or field names, and reads s a word at a time without
// writing a folded copy of it anywhere. As in EqualFold, only 'A'..'Z' and
// 'a'..'z' fold; every other byte, including those >= 0x80, hashes as is.
//
// It is not a cryptographic hash, and it takes no seed, so it offers no
// protection against collisions planned in advance: use it to index keys of
// the program's own, such as the fields of a struct, and not to store keys an
// attacker picks. Its values may change between releases; do not persist them.
func HashFold[S byteSeq](s S) uint64 {
	n := len(s)
	if n < 8 {
		// A key shorter than a word, as most are, packs into one: from four
		// bytes up as two overlapping halves, below that as its first,
		// middle and last byte, which cover one to three bytes between them.
		// Keys of different lengths can pack alike, "a" and "aa" say, so the
		// length goes into the multiplier, where no byte of the key reaches.
		var w uint64
		switch {
		case n >= 4:
			w = uint64(load4(s, 0)) | uint64(load4(s, n-4))<<32
		case n > 0:
			w = uint64(s[0]) | uint64(s[n>>1])<<8 | uint64(s[n-1])<<16
		}
		return hashFoldMix(swar.ToLowerWord(w)^hashFoldK0, hashFoldK1^uint64(n)<<56)
	}
	// Longer keys go as wyhash does: sixteen bytes per multiply, then the
	// last sixteen, overlapping what came before when fewer are left, and a
	// final mix that takes in the length.
	h := hashFoldK2
	for i := 0; n-i > 16; i += 16 {
		h = hashFoldMix(swar.ToLowerWord(swar.Load8(s, i))^hashFoldK1, swar.ToLowerWord(swar.Load8(s, i+8))^h)
	}
	hi, lo := bits.Mul64(
		swar.ToLowerWord(swar.Load8(s, max(n-16, 0)))^hashFoldK1,
		swar.ToLowerWord(swar.Load8(s, n-8))^h,
	)
	return hashFoldMix(lo^hashFoldK0^uint64(n), hi^hashFoldK1)
}

// The constants of HashFold, wyhash's.
const (
	hashFoldK0 uint64 = 0xa0761d6478bd642f
	hashFoldK1 uint64 = 0xe7037ed1a0b428db
	hashFoldK2 uint64 = 0x8ebc6af09c88c6e3
)

// hashFoldMix folds the 128-bit product of a and b into 64 bits, the mixing
// step of the wyhash family: a single widening multiply spreads every bit of
// a and b across the result.
func hashFoldMix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
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
