// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

// MemchrWord returns the index of the first word character — 'A'..'Z',
// 'a'..'z', '0'..'9', or '_' (the regex \w class) — in haystack, or -1 if
// none is present.
func MemchrWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrWord(haystack)
}

// MemchrNotWord is the complement of MemchrWord: it returns the index of the
// first byte that is NOT a word character, or -1 if every byte is one.
func MemchrNotWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrNotWord(haystack)
}

// MemchrInTable returns the index of the first byte b of haystack with
// table[b] true, or -1 if none matches or table is nil. It is a
// general-purpose character-class scan; for the \w class specifically,
// MemchrWord is faster.
func MemchrInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if table[b] {
			return i
		}
	}
	return -1
}

// MemchrNotInTable is the complement of MemchrInTable: it returns the index
// of the first byte b with table[b] false, or -1 if every byte is in the
// table or table is nil.
func MemchrNotInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if !table[b] {
			return i
		}
	}
	return -1
}

// isWordChar reports whether b is in [A-Za-z0-9_].
func isWordChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

// memchrWordGeneric is the portable fallback for MemchrWord.
func memchrWordGeneric(haystack []byte) int {
	for i, b := range haystack {
		if isWordChar(b) {
			return i
		}
	}
	return -1
}

// memchrNotWordGeneric is the portable fallback for MemchrNotWord.
func memchrNotWordGeneric(haystack []byte) int {
	for i, b := range haystack {
		if !isWordChar(b) {
			return i
		}
	}
	return -1
}
