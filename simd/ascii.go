// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"encoding/binary"
)

// IsASCII reports whether every byte in data is ASCII (< 0x80). An empty
// slice is trivially ASCII.
//
// On amd64 with AVX2 and inputs of 32+ bytes it checks 32 bytes per
// iteration with a single VPMOVMSKB; elsewhere it uses a SWAR loop over
// 8-byte words.
func IsASCII(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return isASCII(data)
}

// FirstNonASCII returns the index of the first byte >= 0x80 in data, or -1
// if data is all ASCII. Use it to find where UTF-8 sequences begin.
func FirstNonASCII(data []byte) int {
	for i, b := range data {
		if b >= 0x80 {
			return i
		}
	}
	return -1
}

// CountNonASCII returns the number of bytes >= 0x80 in data.
func CountNonASCII(data []byte) int {
	count := 0
	for _, b := range data {
		if b >= 0x80 {
			count++
		}
	}
	return count
}

// isASCIIGeneric is the portable SWAR implementation behind IsASCII: it ANDs
// each 8-byte word with 0x80 lane masks, so one test covers eight bytes.
func isASCIIGeneric(data []byte) bool {
	n := len(data)
	i := 0
	for ; i+8 <= n; i += 8 {
		if binary.LittleEndian.Uint64(data[i:])&hi8 != 0 {
			return false
		}
	}
	for ; i < n; i++ {
		if data[i] >= 0x80 {
			return false
		}
	}
	return true
}
