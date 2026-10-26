// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// The kernels end with an overlapping 32-byte rescan, so inputs are gated at
// MinLen (MinLen+offset for the pair).
//
//go:noescape
func memchr2AVX2(haystack []byte, needle1, needle2 byte) int

//go:noescape
func memchr3AVX2(haystack []byte, needle1, needle2, needle3 byte) int

//go:noescape
func memchrPairAVX2(haystack []byte, byte1, byte2 byte, offset int) int

// memchr2Impl and memchr3Impl: the caller guarantees a non-empty haystack.
func memchr2Impl(haystack []byte, needle1, needle2 byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchr2AVX2(haystack, needle1, needle2)
	}
	return memchr2Generic(haystack, needle1, needle2)
}

func memchr3Impl(haystack []byte, needle1, needle2, needle3 byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchr3AVX2(haystack, needle1, needle2, needle3)
	}
	return memchr3Generic(haystack, needle1, needle2, needle3)
}

// memchrPairImpl: the caller guarantees offset >= 1 and len(haystack) > offset.
func memchrPairImpl(haystack []byte, byte1, byte2 byte, offset int) int {
	if hasAVX2 && len(haystack) >= MinLen+offset {
		return memchrPairAVX2(haystack, byte1, byte2, offset)
	}
	return memchrPairGeneric(haystack, byte1, byte2, offset)
}
