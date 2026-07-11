package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/swar"
)

// Multi-needle and case-insensitive searching, built on swar.ZeroLanes
// first-match masks with needle broadcasts hoisted out of the word loops.
// IndexAny2/IndexAny3 scan words while i+8 <= n and finish 8+ byte inputs
// with one overlapping word at n-8 (pure loads; already-scanned lanes are
// known non-matching, so the first set lane falls in the new bytes).
// IndexFold instead scans candidate start positions only up to n-k, so its
// word loop hands the final partial word to a scalar remainder; the
// overlapping-window trick appears there only inside the verify reads.
// Inputs shorter than a word use scalar loops throughout.

// IndexAny2 returns the index of the first occurrence in s of either a or b,
// or -1 if neither is present.
func IndexAny2[S byteSeq](s S, a, b byte) int {
	n := len(s)
	if n >= 8 {
		// The needle broadcasts are hoisted here so sub-word inputs never
		// pay for the multiplies.
		bcA := uint64(a) * swar.Ones
		bcB := uint64(b) * swar.Ones
		i := 0
		for ; i+8 <= n; i += 8 {
			w := swar.Load8(s, i)
			if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		if i == n {
			return -1
		}
		w := swar.Load8(s, n-8)
		if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for i := range n {
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
	if n >= 8 {
		bcA := uint64(a) * swar.Ones
		bcB := uint64(b) * swar.Ones
		bcC := uint64(c) * swar.Ones
		i := 0
		for ; i+8 <= n; i += 8 {
			w := swar.Load8(s, i)
			if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB) | swar.ZeroLanes(w^bcC); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		if i == n {
			return -1
		}
		w := swar.Load8(s, n-8)
		if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB) | swar.ZeroLanes(w^bcC); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for i := range n {
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

	if n >= 8 {
		// Scan candidate positions 8 at a time by matching the needle's
		// first byte in both case variants (no need to fold the haystack),
		// then verify each candidate. Approximate-mask false positives just
		// cost a verify; the candidate order stays low-to-high, so the
		// pos > last cutoff remains a valid global exit.
		bc1 := uint64(first) * swar.Ones
		bc2 := bc1
		if first >= 'a' && first <= 'z' {
			bc2 = uint64(first-0x20) * swar.Ones
		}

		// For needles up to one word the verify is a single masked word
		// compare against the folded needle, built lazily on the first
		// candidate; longer needles fold-compare word-at-a-time. The verify
		// block is deliberately spelled out in both loops below: routing it
		// through a closure taxes every call with capture setup even on the
		// no-candidate path, so keep the two copies in sync.
		var needleWord, lenMask uint64
		built := false

		i := 0
		for ; i+8 <= n; i += 8 {
			w := swar.Load8(s, i)
			cand := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2)
			for cand != 0 {
				pos := i + swar.FirstLane(cand)
				if pos > last {
					return -1
				}
				if k <= 8 {
					if !built {
						needleWord = foldNeedle(needle)
						lenMask = ^uint64(0) >> ((8 - k) * 8)
						built = true
					}
					if foldedWindowAt(s, pos, lenMask) == needleWord {
						return pos
					}
				} else if foldEqualAt(s, pos, needle) {
					return pos
				}
				cand &= cand - 1
			}
		}
		// Candidate starts in the final partial word: positions [i, last].
		// Here i is the word loop's exit index n - n%8, so this region is
		// non-empty only when k <= n%8 <= 7 — the needle always fits in one
		// word and only the short-needle verify can apply.
		for ; i <= last; i++ {
			if table[s[i]] == first {
				if !built {
					needleWord = foldNeedle(needle)
					lenMask = ^uint64(0) >> ((8 - k) * 8)
					built = true
				}
				if foldedWindowAt(s, i, lenMask) == needleWord {
					return i
				}
			}
		}
		return -1
	}

	// Scalar path: s shorter than a word (so k <= n < 8).
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

// foldNeedle returns needle lower-cased into one little-endian word, lane 0
// holding needle[0]. len(needle) must be 1..8; lanes past the needle are
// zero, matching the ^uint64(0) >> ((8-k)*8) mask callers apply to windows.
func foldNeedle(needle string) uint64 {
	k := len(needle)
	if k == 8 {
		return swar.ToLowerWord(swar.Load8(needle, 0))
	}
	table := caseconv.ToLowerTable
	var word uint64
	for j := k - 1; j >= 0; j-- {
		word = word<<8 | uint64(table[needle[j]])
	}
	return word
}

// foldedWindowAt returns the lower-cased window of s starting at pos, masked
// to the needle length. Requires len(s) >= 8 and the window to fit in s;
// windows running past len(s)-8 are read through one overlapping load at
// len(s)-8 and shifted so lane 0 is s[pos].
func foldedWindowAt[S byteSeq](s S, pos int, lenMask uint64) uint64 {
	if pos+8 <= len(s) {
		return swar.ToLowerWord(swar.Load8(s, pos)) & lenMask
	}
	from := len(s) - 8
	w := swar.ToLowerWord(swar.Load8(s, from)) >> (8 * (pos - from))
	return w & lenMask
}

// foldEqualAt reports whether s[pos:pos+len(needle)] equals needle ASCII
// case-insensitively. len(needle) must be >= 8 and pos+len(needle) must be
// within s; the tail is one overlapping word compare.
func foldEqualAt[S byteSeq](s S, pos int, needle string) bool {
	k := len(needle)
	i := 0
	for ; i+8 <= k; i += 8 {
		if swar.ToLowerWord(swar.Load8(s, pos+i)) != swar.ToLowerWord(swar.Load8(needle, i)) {
			return false
		}
	}
	if i == k {
		return true
	}
	return swar.ToLowerWord(swar.Load8(s, pos+k-8)) == swar.ToLowerWord(swar.Load8(needle, k-8))
}

// ContainsFold reports whether s contains needle, ASCII case-insensitively.
func ContainsFold[S byteSeq](s S, needle string) bool {
	return IndexFold(s, needle) >= 0
}
