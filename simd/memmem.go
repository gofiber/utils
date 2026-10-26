// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"
)

// memmemMinPairedHaystack: up to 64 bytes, bytes.Index's brute force wins.
const memmemMinPairedHaystack = 65

// memmemMinSingleHaystack: higher, as the single-byte scan is less selective.
const memmemMinSingleHaystack = 128

// memmemMaxPairNeedle: up to this length, single-byte false hits are frequent.
const memmemMaxPairNeedle = 6

// memmemMaxMisses: caps failed verifications, bounding the worst case.
const memmemMaxMisses = 16

// Memmem returns the index of the first occurrence of needle in haystack,
// or -1 if needle is not present. An empty needle matches at index 0,
// mirroring bytes.Index.
//
// On amd64 with AVX2 it prefilters on the needle's rarest bytes; the worst
// case stays within a constant of bytes.Index's O(n+m).
func Memmem(haystack, needle []byte) int {
	needleLen := len(needle)
	if needleLen == 0 {
		return 0
	}
	if len(haystack) < needleLen {
		return -1
	}
	if needleLen == 1 {
		return bytes.IndexByte(haystack, needle[0])
	}
	if !hasAVX2 || len(haystack) < memmemMinPairedHaystack {
		return bytes.Index(haystack, needle)
	}

	var info RareByteInfo
	pairable := needleLen <= memmemMaxPairNeedle
	if pairable {
		info = SelectRareBytes(needle)
		// Distinct byte values imply distinct indices, so Byte1 != Byte2
		// guarantees the offset >= 1 that MemchrPair needs.
		if info.Byte1 != info.Byte2 {
			return memmemPaired(haystack, needle, info)
		}
	}
	if len(haystack) < memmemMinSingleHaystack {
		return bytes.Index(haystack, needle)
	}
	if !pairable {
		info = SelectRareBytes(needle)
	}
	return memmemSingle(haystack, needle, info.Byte1, info.Index1)
}

func memmemPaired(haystack, needle []byte, info RareByteInfo) int {
	byte1, idx1 := info.Byte1, info.Index1
	byte2, idx2 := info.Byte2, info.Index2
	if idx1 > idx2 {
		byte1, byte2 = byte2, byte1
		idx1, idx2 = idx2, idx1
	}
	offset := idx2 - idx1

	misses := 0
	for start := 0; ; {
		cand := MemchrPair(haystack[start:], byte1, byte2, offset)
		if cand == -1 {
			return -1
		}
		// cand is the position of byte1; the needle would start idx1 earlier.
		pos := start + cand - idx1
		if pos >= 0 && pos+len(needle) <= len(haystack) &&
			bytes.Equal(haystack[pos:pos+len(needle)], needle) {
			return pos
		}
		start += cand + 1
		if misses++; misses >= memmemMaxMisses {
			return memmemFallback(haystack, needle, start, idx1)
		}
	}
}

func memmemSingle(haystack, needle []byte, rareByte byte, rareIdx int) int {
	misses := 0
	for start := 0; ; {
		cand := bytes.IndexByte(haystack[start:], rareByte)
		if cand == -1 {
			return -1
		}
		pos := start + cand - rareIdx
		if pos >= 0 && pos+len(needle) <= len(haystack) &&
			bytes.Equal(haystack[pos:pos+len(needle)], needle) {
			return pos
		}
		start += cand + 1
		if misses++; misses >= memmemMaxMisses {
			return memmemFallback(haystack, needle, start, rareIdx)
		}
	}
}

// memmemFallback searches from start-anchorIdx: anchors before start were
// checked and a remaining match has its anchor at anchorIdx, so none is missed.
func memmemFallback(haystack, needle []byte, start, anchorIdx int) int {
	from := max(start-anchorIdx, 0)
	if pos := bytes.Index(haystack[from:], needle); pos >= 0 {
		return from + pos
	}
	return -1
}
