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

// wordMask flags the lanes of w holding word characters, exactly per lane;
// bytes >= 0x80 never match. It is the swar.MatchRangeMask construction
// with the three range tests fused into one expression:
//
//   - the ASCII gate (&^ w & HighBits) is applied once to the OR of the
//     tests instead of once per test;
//   - 'A'..'Z' and 'a'..'z' collapse into a single 'a'..'z' test on a
//     case-folded copy (bit 5 set in every lane), which maps exactly the
//     two letter ranges onto 0x61..0x7A;
//   - '_' (0x5F) is the only byte besides DEL (0x7F) that folds to 0x7F,
//     so "folded lane == 0x7F" is one biased add, and the DEL lanes are
//     removed by w<<2, which moves each lane's bit 5 — clear for '_', set
//     for DEL — into the bit-7 marker position.
//
// Every addition below adds a per-lane bias of at most 0x50 to lanes that
// are at most 0x7F, so no lane can carry into its neighbor and the bit-7
// markers stay exact. That is 15 ALU operations per word instead of 26
// for the four separate masks.
func wordMask(w uint64) uint64 {
	b := w & swar.LowSeven
	bf := b | 0x20*swar.Ones
	letters := (bf + (0x80-'a')*swar.Ones) &^ (bf + (0x80-'z'-1)*swar.Ones)
	digits := (b + (0x80-'0')*swar.Ones) &^ (b + (0x80-'9'-1)*swar.Ones)
	underscore := (bf + swar.Ones) &^ (w << 2)
	return (letters | digits | underscore) &^ w & swar.HighBits
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
