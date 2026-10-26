// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

//go:noescape
func memchrWordAVX2(haystack []byte) int

//go:noescape
func memchrNotWordAVX2(haystack []byte) int

// The caller guarantees a non-empty haystack (also for memchrNotWordImpl).
func memchrWordImpl(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrWordAVX2(haystack)
	}
	return memchrWordGeneric(haystack)
}

func memchrNotWordImpl(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrNotWordAVX2(haystack)
	}
	return memchrNotWordGeneric(haystack)
}
