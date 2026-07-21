// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// Assembly kernels implemented in memchr_amd64.s. Each processes 32 bytes
// per iteration with AVX2 and finishes the tail with a scalar loop.
//
//go:noescape
func memchr2AVX2(haystack []byte, needle1, needle2 byte) int

//go:noescape
func memchr3AVX2(haystack []byte, needle1, needle2, needle3 byte) int

//go:noescape
func memchrPairAVX2(haystack []byte, byte1, byte2 byte, offset int) int

// memchr2 dispatches to the AVX2 kernel for inputs large enough to amortize
// the vector setup. The caller guarantees a non-empty haystack.
func memchr2(haystack []byte, needle1, needle2 byte) int {
	if hasAVX2 && len(haystack) >= 32 {
		return memchr2AVX2(haystack, needle1, needle2)
	}
	return memchr2Generic(haystack, needle1, needle2)
}

// memchr3 dispatches like memchr2.
func memchr3(haystack []byte, needle1, needle2, needle3 byte) int {
	if hasAVX2 && len(haystack) >= 32 {
		return memchr3AVX2(haystack, needle1, needle2, needle3)
	}
	return memchr3Generic(haystack, needle1, needle2, needle3)
}

// memchrPair dispatches to the AVX2 kernel when there is room for at least
// one 32-byte vector of pair positions. The caller guarantees offset >= 1
// and len(haystack) > offset.
func memchrPair(haystack []byte, byte1, byte2 byte, offset int) int {
	if hasAVX2 && len(haystack) >= 32+offset {
		return memchrPairAVX2(haystack, byte1, byte2, offset)
	}
	return memchrPairGeneric(haystack, byte1, byte2, offset)
}
