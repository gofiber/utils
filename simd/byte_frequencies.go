// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

// byteFrequencies holds empirical byte frequency ranks (English text, source
// code, binary samples), as in Rust's memchr crate.
var byteFrequencies = [256]byte{
	// 16 ranks per row, from byte 0x00; 0x80-0xFF, UTF-8 bytes, all rank 5.
	0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	255, 60, 140, 50, 40, 35, 30, 160, 130, 130, 80, 55, 200, 140, 210, 100,
	180, 190, 170, 150, 140, 140, 130, 120, 120, 120, 150, 100, 70, 160, 70, 50,
	25, 120, 80, 90, 85, 130, 75, 70, 80, 115, 30, 35, 90, 85, 100, 105,
	80, 15, 100, 110, 115, 70, 45, 55, 20, 50, 10, 90, 60, 90, 20, 110,
	30, 225, 140, 170, 165, 245, 135, 130, 150, 200, 25, 65, 175, 155, 195, 205,
	145, 15, 195, 200, 215, 150, 75, 95, 45, 120, 20, 85, 40, 85, 15, 0,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
}

// ByteRank returns the empirical frequency rank of b; lower means rarer.
func ByteRank(b byte) byte {
	return byteFrequencies[b]
}

// RareByteInfo describes the two rarest bytes of a needle. Byte2 differs from
// Byte1 whenever the needle has two distinct values; otherwise Byte2/Index2
// equal Byte1/Index1.
type RareByteInfo struct {
	// Byte1 is the needle's rarest byte.
	Byte1 byte
	// Byte2 is the rarest byte other than Byte1.
	Byte2 byte
	// Index1 and Index2 are the first positions of Byte1 and Byte2.
	Index1 int
	Index2 int
}

// SelectRareBytes returns the two rarest bytes of needle according to ByteRank,
// preferring distinct values so repeats of the rarest byte never crowd out the
// second anchor. Memmem feeds the result to MemchrPair.
func SelectRareBytes(needle []byte) RareByteInfo {
	n := len(needle)
	if n == 0 {
		return RareByteInfo{}
	}

	byte1, idx1 := needle[0], 0
	byte2, idx2 := needle[0], 0
	have2 := false

	for i := 1; i < n; i++ {
		b := needle[i]
		switch {
		case byteFrequencies[b] < byteFrequencies[byte1]:
			// Demoted byte1 differs from b and is no more common than byte2.
			byte2, idx2 = byte1, idx1
			byte1, idx1 = b, i
			have2 = true
		case b != byte1 && (!have2 || byteFrequencies[b] < byteFrequencies[byte2]):
			byte2, idx2 = b, i
			have2 = true
		}
	}
	if !have2 {
		byte2, idx2 = byte1, idx1
	}

	return RareByteInfo{Byte1: byte1, Byte2: byte2, Index1: idx1, Index2: idx2}
}
