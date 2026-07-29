// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"github.com/gofiber/utils/v2/swar"
)

// MemchrWord returns the index of the first word character — 'A'..'Z',
// 'a'..'z', '0'..'9', or '_' (the regex \w class) — in haystack, or -1 if
// none is present.
func MemchrWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrWordImpl(haystack)
}

// MemchrNotWord is the complement of MemchrWord: it returns the index of the
// first byte that is NOT a word character, or -1 if every byte is one.
func MemchrNotWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrNotWordImpl(haystack)
}

// MemchrInTable returns the index of the first byte b of haystack with
// table[b] true, or -1 if none matches or table is nil.
//
// This is a plain scalar table scan on every architecture — an arbitrary
// 256-entry membership test has no cheap SIMD or SWAR form — so unlike the
// rest of this package it carries no acceleration. For the \w class
// specifically, MemchrWord is much faster.
func MemchrInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if table[b] {
			return i
		}
	}
	return -1
}

// MemchrNotInTable is the complement of MemchrInTable: it returns the index
// of the first byte b with table[b] false, or -1 if every byte is in the
// table or table is nil. Like MemchrInTable it is a plain scalar loop on
// every architecture.
func MemchrNotInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if !table[b] {
			return i
		}
	}
	return -1
}

// isWordChar reports whether b is in [A-Za-z0-9_].
func isWordChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

// wordMask flags the lanes of w holding word characters. The range masks
// are exact per lane and bytes >= 0x80 never match, so the composition is
// exact too.
func wordMask(w uint64) uint64 {
	return swar.MatchRangeMask(w, 'A', 'Z') |
		swar.MatchRangeMask(w, 'a', 'z') |
		swar.MatchRangeMask(w, '0', '9') |
		swar.MatchByteMask(w, '_')
}

// memchrWordGeneric is the portable SWAR fallback for MemchrWord, finishing
// 8+ byte inputs with one overlapping word at n-8.
func memchrWordGeneric(haystack []byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if isWordChar(haystack[i]) {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		// Pinned window: see the note above isASCIIGeneric in ascii.go for
		// why the loads use constant offsets into it.
		w := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := wordMask(swar.Load8(w, 0))
		m1 := wordMask(swar.Load8(w, swar.WordLen))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		if m := wordMask(swar.Load8(haystack, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := wordMask(swar.Load8(haystack, n-swar.WordLen)); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// memchrNotWordGeneric is the portable SWAR fallback for MemchrNotWord.
// wordMask is exact per lane, so its complement within the 0x80 lane
// markers is exact as well.
func memchrNotWordGeneric(haystack []byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if !isWordChar(haystack[i]) {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		w := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := swar.HighBits &^ wordMask(swar.Load8(w, 0))
		m1 := swar.HighBits &^ wordMask(swar.Load8(w, swar.WordLen))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		if m := swar.HighBits &^ wordMask(swar.Load8(haystack, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := swar.HighBits &^ wordMask(swar.Load8(haystack, n-swar.WordLen)); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}
