// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"math/bits"

	"github.com/gofiber/utils/v2/swar"
)

// IsASCII reports whether every byte in data is ASCII (< 0x80). An empty
// slice is trivially ASCII.
//
// On amd64 with AVX2 and inputs of MinLen+ bytes it checks 32 bytes per
// iteration with a single VPMOVMSKB; elsewhere it uses a SWAR loop over
// 8-byte words.
func IsASCII(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return isASCIIImpl(data)
}

// FirstNonASCII returns the index of the first byte >= 0x80 in data, or -1
// if data is all ASCII. Use it to find where UTF-8 sequences begin. It scans
// with SWAR words (no AVX2 kernel).
func FirstNonASCII(data []byte) int {
	n := len(data)
	if n < swar.WordLen {
		for i := range n {
			if data[i] >= 0x80 {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		if m := swar.Load8(data, i) & swar.HighBits; m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	// Overlapping final word: lanes before i are known ASCII, so any set
	// lane falls in the new bytes.
	if m := swar.Load8(data, n-swar.WordLen) & swar.HighBits; m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// CountNonASCII returns the number of bytes >= 0x80 in data. It counts
// SWAR-word-wise with a popcount per word (no AVX2 kernel).
func CountNonASCII(data []byte) int {
	count := 0
	i := 0
	n := len(data)
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		count += bits.OnesCount64(swar.Load8(data, i) & swar.HighBits)
	}
	// No overlapping word here: it would double-count. The tail is at most
	// seven bytes.
	for ; i < n; i++ {
		if data[i] >= 0x80 {
			count++
		}
	}
	return count
}

// isASCIIGeneric is the portable SWAR implementation behind IsASCII: it
// masks each 8-byte word with the 0x80 lane bits, finishing 8+ byte inputs
// with one overlapping word at n-8.
func isASCIIGeneric(data []byte) bool {
	n := len(data)
	if n < swar.WordLen {
		for i := range n {
			if data[i] >= 0x80 {
				return false
			}
		}
		return true
	}
	i := 0
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		if swar.Load8(data, i)&swar.HighBits != 0 {
			return false
		}
	}
	if i == n {
		return true
	}
	return swar.Load8(data, n-swar.WordLen)&swar.HighBits == 0
}
