// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

// ByteFrequencies holds empirical byte frequency ranks derived from English
// text, source code, and binary data samples. Lower rank means rarer byte,
// which makes it a better anchor for a search prefilter. The approach
// mirrors Rust's memchr crate (https://github.com/BurntSushi/memchr).
var ByteFrequencies = [256]byte{
	// 0x00-0x0F: control characters (generally rare)
	0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0,
	// 0x10-0x1F: more control characters
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0x20-0x2F: space and punctuation
	255, 60, 140, 50, 40, 35, 30, 160, 130, 130, 80, 55, 200, 140, 210, 100,
	// 0x30-0x3F: digits and more punctuation
	180, 190, 170, 150, 140, 140, 130, 120, 120, 120, 150, 100, 70, 160, 70, 50,
	// 0x40-0x4F: '@' and uppercase A-O
	25, 120, 80, 90, 85, 130, 75, 70, 80, 115, 30, 35, 90, 85, 100, 105,
	// 0x50-0x5F: uppercase P-Z and brackets
	80, 15, 100, 110, 115, 70, 45, 55, 20, 50, 10, 90, 60, 90, 20, 110,
	// 0x60-0x6F: backtick and lowercase a-o
	30, 225, 140, 170, 165, 245, 135, 130, 150, 200, 25, 65, 175, 155, 195, 205,
	// 0x70-0x7F: lowercase p-z and braces
	145, 15, 195, 200, 215, 150, 75, 95, 45, 120, 20, 85, 40, 85, 15, 0,
	// 0x80-0xFF: UTF-8 continuation/lead bytes (rare in typical text)
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
}

// ByteRank returns the frequency rank of b. Lower values indicate rarer
// bytes, which are better anchors for search prefilters.
func ByteRank(b byte) byte {
	return ByteFrequencies[b]
}

// RareByteInfo describes the two rarest bytes of a needle, as selected by
// SelectRareBytes. For needles shorter than two bytes, or needles whose
// bytes are all identical, Byte2/Index2 may equal Byte1/Index1.
type RareByteInfo struct {
	// Byte1 is the rarest byte found in the needle.
	Byte1 byte
	// Byte2 is the second-rarest byte, different from Byte1 when possible.
	Byte2 byte
	// Index1 is the position of Byte1 in the needle.
	Index1 int
	// Index2 is the position of Byte2 in the needle.
	Index2 int
}

// SelectRareBytes returns the two rarest bytes of needle according to
// ByteFrequencies, preferring distinct byte values. Memmem feeds the result
// to MemchrPair to build a highly selective substring prefilter; it is
// exported so downstream code can build its own prefilters the same way.
func SelectRareBytes(needle []byte) RareByteInfo {
	n := len(needle)
	if n == 0 {
		return RareByteInfo{}
	}
	if n == 1 {
		return RareByteInfo{Byte1: needle[0], Byte2: needle[0]}
	}

	byte1, idx1 := needle[0], 0
	byte2, idx2 := needle[1], 1
	if ByteFrequencies[byte2] < ByteFrequencies[byte1] {
		byte1, byte2 = byte2, byte1
		idx1, idx2 = idx2, idx1
	}

	for i := 2; i < n; i++ {
		b := needle[i]
		rank := ByteFrequencies[b]
		switch {
		case rank < ByteFrequencies[byte1]:
			byte2, idx2 = byte1, idx1
			byte1, idx1 = b, i
		case b != byte1 && rank < ByteFrequencies[byte2]:
			byte2, idx2 = b, i
		}
	}

	return RareByteInfo{Byte1: byte1, Byte2: byte2, Index1: idx1, Index2: idx2}
}
