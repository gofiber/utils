// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

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

// memchrDigitGeneric is the portable fallback for MemchrDigit.
func memchrDigitGeneric(haystack []byte) int {
	for i, b := range haystack {
		if b >= '0' && b <= '9' {
			return i
		}
	}
	return -1
}
