// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build !amd64

package simd

// memchrDigit uses the portable implementation off amd64.
func memchrDigit(haystack []byte) int {
	return memchrDigitGeneric(haystack)
}
