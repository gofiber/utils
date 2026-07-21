// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"
	"encoding/binary"
	"math/bits"
)

// Memchr2 returns the index of the first occurrence in haystack of either
// needle1 or needle2, or -1 if neither is present. Both needles are checked
// in the same pass, so it costs the same as a single-byte scan.
//
// There is deliberately no single-needle Memchr: use bytes.IndexByte, which
// the Go runtime already accelerates on all major architectures.
func Memchr2(haystack []byte, needle1, needle2 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr2(haystack, needle1, needle2)
}

// Memchr3 returns the index of the first occurrence in haystack of needle1,
// needle2, or needle3, or -1 if none is present. All three needles are
// checked in the same pass.
func Memchr3(haystack []byte, needle1, needle2, needle3 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr3(haystack, needle1, needle2, needle3)
}

// MemchrPair returns the first index i such that haystack[i] == byte1 and
// haystack[i+offset] == byte2, or -1 if no such position exists. offset must
// be >= 0; a negative offset returns -1.
//
// Requiring two bytes at a fixed distance makes the scan far more selective
// than a single-byte search, which is what makes it a good substring
// prefilter: false candidates need both bytes at exactly the right distance.
// Memmem uses it for short needles.
func MemchrPair(haystack []byte, byte1, byte2 byte, offset int) int {
	if offset < 0 || len(haystack) <= offset {
		return -1
	}
	if offset == 0 {
		// Both bytes would occupy the same position, so they must be equal.
		if byte1 != byte2 {
			return -1
		}
		return bytes.IndexByte(haystack, byte1)
	}
	return memchrPair(haystack, byte1, byte2, offset)
}

// memchr2Generic is the portable SWAR fallback for Memchr2. Matching lanes
// XOR to zero; the zero-lane formula can flag lanes above the first true
// match, but the trailing-zero pick always lands on a true match.
func memchr2Generic(haystack []byte, needle1, needle2 byte) int {
	n := len(haystack)
	if n < 8 {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 {
				return i
			}
		}
		return -1
	}

	mask1 := uint64(needle1) * lo8
	mask2 := uint64(needle2) * lo8

	i := 0
	for ; i+8 <= n; i += 8 {
		chunk := binary.LittleEndian.Uint64(haystack[i:])
		x1 := chunk ^ mask1
		x2 := chunk ^ mask2
		zero := ((x1 - lo8) &^ x1 & hi8) | ((x2 - lo8) &^ x2 & hi8)
		if zero != 0 {
			return i + bits.TrailingZeros64(zero)/8
		}
	}
	if i == n {
		return -1
	}
	// Finish with one overlapping word at n-8: lanes already scanned are
	// known non-matching, so the first set lane falls in the new bytes.
	chunk := binary.LittleEndian.Uint64(haystack[n-8:])
	x1 := chunk ^ mask1
	x2 := chunk ^ mask2
	zero := ((x1 - lo8) &^ x1 & hi8) | ((x2 - lo8) &^ x2 & hi8)
	if zero != 0 {
		return n - 8 + bits.TrailingZeros64(zero)/8
	}
	return -1
}

// memchr3Generic is the portable SWAR fallback for Memchr3; see
// memchr2Generic for the technique.
func memchr3Generic(haystack []byte, needle1, needle2, needle3 byte) int {
	n := len(haystack)
	if n < 8 {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 || haystack[i] == needle3 {
				return i
			}
		}
		return -1
	}

	mask1 := uint64(needle1) * lo8
	mask2 := uint64(needle2) * lo8
	mask3 := uint64(needle3) * lo8

	i := 0
	for ; i+8 <= n; i += 8 {
		chunk := binary.LittleEndian.Uint64(haystack[i:])
		x1 := chunk ^ mask1
		x2 := chunk ^ mask2
		x3 := chunk ^ mask3
		zero := ((x1 - lo8) &^ x1 & hi8) | ((x2 - lo8) &^ x2 & hi8) | ((x3 - lo8) &^ x3 & hi8)
		if zero != 0 {
			return i + bits.TrailingZeros64(zero)/8
		}
	}
	if i == n {
		return -1
	}
	chunk := binary.LittleEndian.Uint64(haystack[n-8:])
	x1 := chunk ^ mask1
	x2 := chunk ^ mask2
	x3 := chunk ^ mask3
	zero := ((x1 - lo8) &^ x1 & hi8) | ((x2 - lo8) &^ x2 & hi8) | ((x3 - lo8) &^ x3 & hi8)
	if zero != 0 {
		return n - 8 + bits.TrailingZeros64(zero)/8
	}
	return -1
}

// memchrPairGeneric is the portable SWAR fallback for MemchrPair. The two
// per-position masks are ANDed, so a lane can be flagged without both bytes
// truly matching (zero-lane false positives no longer sit above a true match
// after the AND); every candidate is therefore re-verified scalar before
// being returned. The caller guarantees offset >= 1 and len > offset.
func memchrPairGeneric(haystack []byte, byte1, byte2 byte, offset int) int {
	n := len(haystack)
	if n < 8+offset {
		for i := 0; i+offset < n; i++ {
			if haystack[i] == byte1 && haystack[i+offset] == byte2 {
				return i
			}
		}
		return -1
	}

	mask1 := uint64(byte1) * lo8
	mask2 := uint64(byte2) * lo8

	i := 0
	for ; i+8+offset <= n; i += 8 {
		chunk1 := binary.LittleEndian.Uint64(haystack[i:])
		chunk2 := binary.LittleEndian.Uint64(haystack[i+offset:])
		x1 := chunk1 ^ mask1
		x2 := chunk2 ^ mask2
		cand := ((x1 - lo8) &^ x1 & hi8) & ((x2 - lo8) &^ x2 & hi8)
		for cand != 0 {
			pos := i + bits.TrailingZeros64(cand)/8
			if haystack[pos] == byte1 && haystack[pos+offset] == byte2 {
				return pos
			}
			cand &= cand - 1
		}
	}
	for ; i+offset < n; i++ {
		if haystack[i] == byte1 && haystack[i+offset] == byte2 {
			return i
		}
	}
	return -1
}
