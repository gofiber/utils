// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"
)

// memmemMinHaystack is the haystack size below which bytes.Index wins: the
// rare-byte selection is a per-call cost that only pays off once the scan
// itself is long enough to be the dominant term. The crossover sits between
// the 96B and 128B points of Benchmark_Memmem (bytes.Index wins at 96B,
// the prefilter from 128B up).
const memmemMinHaystack = 128

// memmemMaxPairNeedle is the needle length up to which the paired-byte
// prefilter beats a single rare-byte scan: short needles produce frequent
// single-byte false positives, and verifying them costs more than the second
// byte check saves.
const memmemMaxPairNeedle = 6

// memmemMaxMisses caps how many failed candidate verifications the
// prefilter loops tolerate before handing the search to bytes.Index. A
// pathological input (rare byte everywhere, long almost-matching needle)
// would otherwise cost O(len(haystack) * len(needle)); the cap bounds the
// wasted work at memmemMaxMisses needle comparisons, keeping the worst case
// at bytes.Index's O(n+m) plus a constant.
const memmemMaxMisses = 16

// Memmem returns the index of the first occurrence of needle in haystack,
// or -1 if needle is not present. An empty needle matches at index 0,
// mirroring bytes.Index.
//
// On amd64 with AVX2 and haystacks of 128+ bytes it scans for the needle's
// two rarest bytes (per ByteRank) at their exact relative distance and
// verifies only those candidates, which typically beats bytes.Index by a
// wide margin on large inputs. In all other cases it delegates to
// bytes.Index directly, and the prefilter itself falls back to bytes.Index
// after a bounded number of failed candidate verifications, so the worst
// case stays within a constant of the stdlib's O(n+m).
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
	if !hasAVX2 || len(haystack) < memmemMinHaystack {
		return bytes.Index(haystack, needle)
	}

	info := SelectRareBytes(needle)
	if needleLen <= memmemMaxPairNeedle && info.Byte1 != info.Byte2 && info.Index1 != info.Index2 {
		return memmemPaired(haystack, needle, info)
	}
	return memmemSingle(haystack, needle, info.Byte1, info.Index1)
}

// memmemPaired scans for the needle's two rarest bytes at their exact
// relative distance with MemchrPair and verifies each candidate. False
// candidates need both bytes at the right distance, so verification is
// rare; if it still misses too often, memmemFallback finishes the search.
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
		if misses++; misses > memmemMaxMisses {
			return memmemFallback(haystack, needle, start, idx1)
		}
	}
}

// memmemSingle scans for the needle's rarest byte and verifies each
// candidate; for longer needles the rare byte alone is selective enough
// that the paired scan's second comparison stops paying for itself. Like
// memmemPaired it hands persistent misses to memmemFallback.
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
		if misses++; misses > memmemMaxMisses {
			return memmemFallback(haystack, needle, start, rareIdx)
		}
	}
}

// memmemFallback finishes a prefiltered search with bytes.Index once the
// candidate loops exceed their miss budget. Every anchor before start has
// been checked, and any remaining occurrence carries its anchor byte at
// anchorIdx, so it starts at or after start-anchorIdx; searching from there
// cannot skip a match or return one out of order.
func memmemFallback(haystack, needle []byte, start, anchorIdx int) int {
	from := max(start-anchorIdx, 0)
	if pos := bytes.Index(haystack[from:], needle); pos >= 0 {
		return from + pos
	}
	return -1
}
