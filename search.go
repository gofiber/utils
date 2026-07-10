package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/swar"
)

// Multi-needle and case-insensitive searching. Word loops scan 8 bytes per
// iteration with i+8 <= n bounds; inputs of 8+ bytes finish with one
// overlapping word at n-8 (pure loads, and lanes already scanned are known
// non-matching, so the first set lane always falls in the new bytes), while
// shorter inputs use scalar loops.

// IndexAny2 returns the index of the first occurrence in s of either a or b,
// or -1 if neither is present.
func IndexAny2[S byteSeq](s S, a, b byte) int {
	n := len(s)
	i := 0
	for ; i+8 <= n; i += 8 {
		w := swar.Load8(s, i)
		if m := swar.MatchByteMask(w, a) | swar.MatchByteMask(w, b); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		w := swar.Load8(s, n-8)
		if m := swar.MatchByteMask(w, a) | swar.MatchByteMask(w, b); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if s[i] == a || s[i] == b {
			return i
		}
	}
	return -1
}

// IndexAny3 returns the index of the first occurrence in s of a, b, or c,
// or -1 if none is present.
func IndexAny3[S byteSeq](s S, a, b, c byte) int {
	n := len(s)
	i := 0
	for ; i+8 <= n; i += 8 {
		w := swar.Load8(s, i)
		if m := swar.MatchByteMask(w, a) | swar.MatchByteMask(w, b) | swar.MatchByteMask(w, c); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		w := swar.Load8(s, n-8)
		if m := swar.MatchByteMask(w, a) | swar.MatchByteMask(w, b) | swar.MatchByteMask(w, c); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if s[i] == a || s[i] == b || s[i] == c {
			return i
		}
	}
	return -1
}

// IndexFold returns the index of the first ASCII case-insensitive occurrence
// of needle in s, or -1 if absent. An empty needle matches at index 0. Only
// 'A'..'Z'/'a'..'z' fold; every other byte (including >= 0x80) must match
// exactly, so e.g. "no\rcache" does NOT match the needle "no-cache".
func IndexFold[S byteSeq](s S, needle string) int {
	n, k := len(s), len(needle)
	if k == 0 {
		return 0
	}
	if k > n {
		return -1
	}

	table := caseconv.ToLowerTable
	first := table[needle[0]]
	last := n - k // last valid start position

	if k <= 8 && n >= 8 {
		// Scan candidate positions 8 at a time: fold a word of s, mark the
		// lanes matching the needle's (folded) first byte, and verify each
		// candidate with a single masked word compare. The folded needle
		// word is built lazily so misses without a single first-byte
		// candidate never pay for it.
		var needleWord, lenMask uint64

		i := 0
		for ; i+8 <= n; i += 8 {
			cand := swar.MatchByteMask(swar.ToLowerWord(swar.Load8(s, i)), first)
			for cand != 0 {
				pos := i + swar.FirstLane(cand)
				if pos > last {
					return -1
				}
				if lenMask == 0 {
					needleWord, lenMask = foldNeedle(needle)
				}
				if foldedWindowAt(s, pos, lenMask) == needleWord {
					return pos
				}
				cand &= cand - 1
			}
		}
		// Candidate starts in the final partial word: positions [i, last].
		for ; i <= last; i++ {
			if table[s[i]] == first {
				if lenMask == 0 {
					needleWord, lenMask = foldNeedle(needle)
				}
				if foldedWindowAt(s, i, lenMask) == needleWord {
					return i
				}
			}
		}
		return -1
	}

	// Scalar path: needle longer than a word, or s shorter than a word.
outer:
	for i := 0; i <= last; i++ {
		if table[s[i]] != first {
			continue
		}
		for j := 1; j < k; j++ {
			if table[s[i+j]] != table[needle[j]] {
				continue outer
			}
		}
		return i
	}
	return -1
}

// foldNeedle returns needle lower-cased into one little-endian word (lane 0
// is needle[0]) plus the mask covering its lanes. len(needle) must be 1..8;
// the mask is never zero, which IndexFold uses as its "not built yet" state.
func foldNeedle(needle string) (uint64, uint64) {
	k := len(needle)
	if k == 8 {
		return swar.ToLowerWord(swar.Load8(needle, 0)), ^uint64(0)
	}
	table := caseconv.ToLowerTable
	var word uint64
	for j := k - 1; j >= 0; j-- {
		word = word<<8 | uint64(table[needle[j]])
	}
	return word, ^uint64(0) >> ((8 - uint(k)) * 8)
}

// foldedWindowAt returns the lower-cased window of s starting at pos, masked
// to the needle length. Requires len(s) >= 8 and pos+popcount(lenMask)/8 <=
// len(s); windows running past len(s)-8 are read through one overlapping
// load at len(s)-8 and shifted so lane 0 is s[pos].
func foldedWindowAt[S byteSeq](s S, pos int, lenMask uint64) uint64 {
	if pos+8 <= len(s) {
		return swar.ToLowerWord(swar.Load8(s, pos)) & lenMask
	}
	from := len(s) - 8
	w := swar.ToLowerWord(swar.Load8(s, from)) >> (8 * uint(pos-from))
	return w & lenMask
}

// ContainsFold reports whether s contains needle, ASCII case-insensitively.
func ContainsFold[S byteSeq](s S, needle string) bool {
	return IndexFold(s, needle) >= 0
}
