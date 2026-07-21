// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// Assembly kernels implemented in memchr_class_amd64.s. Range membership is
// computed with unsigned clamp-and-compare (VPMINUB/VPMAXUB + VPCMPEQB).
//
//go:noescape
func memchrWordAVX2(haystack []byte) int

//go:noescape
func memchrNotWordAVX2(haystack []byte) int

// memchrWord dispatches to the AVX2 kernel for inputs large enough to
// amortize the vector setup. The caller guarantees a non-empty haystack.
func memchrWord(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrWordAVX2(haystack)
	}
	return memchrWordGeneric(haystack)
}

// memchrNotWord dispatches like memchrWord.
func memchrNotWord(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrNotWordAVX2(haystack)
	}
	return memchrNotWordGeneric(haystack)
}
