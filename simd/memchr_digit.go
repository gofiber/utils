// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"github.com/gofiber/utils/v2/swar"
)

// MemchrDigit returns the index of the first ASCII digit '0'..'9' in
// haystack, or -1 if no digit is present. Only ASCII digits match; Unicode
// digit characters do not.
//
// On amd64 with AVX2 the range check runs on 32 bytes per iteration, which
// is useful for locating numeric fields (ports, lengths, IP octets) in
// request data.
func MemchrDigit(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrDigit(haystack)
}

// MemchrDigitAt is MemchrDigit(haystack[at:]) with the result translated
// back to an absolute index in haystack. It returns -1 when no digit is
// found or when at is out of bounds.
func MemchrDigitAt(haystack []byte, at int) int {
	if at < 0 || at >= len(haystack) {
		return -1
	}
	pos := MemchrDigit(haystack[at:])
	if pos < 0 {
		return -1
	}
	return pos + at
}

// memchrDigitGeneric is the portable SWAR fallback for MemchrDigit: a
// MatchRangeMask first-match scan (exact per lane; bytes >= 0x80 never
// match), finishing 8+ byte inputs with one overlapping word at n-8.
func memchrDigitGeneric(haystack []byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i, b := range haystack {
			if b >= '0' && b <= '9' {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		if m := swar.MatchRangeMask(swar.Load8(haystack, i), '0', '9'); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := swar.MatchRangeMask(swar.Load8(haystack, n-swar.WordLen), '0', '9'); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}
