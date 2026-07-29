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
// On amd64 with AVX2 and inputs of MinLen+ bytes it ORs four 32-byte
// vectors and tests all 128 bytes with a single VPMOVMSKB; elsewhere it
// uses a SWAR loop over 8-byte words.
func IsASCII(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return isASCIIImpl(data)
}

// FirstNonASCII returns the index of the first byte >= 0x80 in data, or -1
// if data is all ASCII. Use it to find where UTF-8 sequences begin.
//
// It reads the same signal as IsASCII — a non-ASCII byte is exactly a byte
// with its high bit set — so on amd64 with AVX2 and inputs of MinLen+
// bytes it also tests 128 bytes per VPMOVMSKB, keeping the four vectors
// live so the hit path can pinpoint the byte; elsewhere it falls back to
// SWAR words.
func FirstNonASCII(data []byte) int {
	if len(data) == 0 {
		return -1
	}
	return firstNonASCIIImpl(data)
}

// CountNonASCII returns the number of bytes >= 0x80 in data.
//
// On amd64 with AVX2 and inputs of MinLen+ bytes it turns each 32-byte
// vector into a high-bit mask and adds its population count; elsewhere it
// counts SWAR-word-wise with a popcount per word.
func CountNonASCII(data []byte) int {
	return countNonASCIIImpl(data)
}

// The unrolled loops below all pin their group of words with a
// three-index reslice before loading from it. swar.Load8's own reslice is
// bounds-checked against the backing array's capacity, which the loop
// condition (stated in terms of length) does not imply, so a variable
// index costs a compare, a branch and a four-instruction empty-slice
// pointer clamp per load. Pinning the window once pays that at group
// granularity and leaves the loads inside it at constant offsets, where
// the checks fold away entirely: measured -20% at 32B rising to -62% at
// 4KiB for the ASCII scan, and the same shape for the first-match scans.

// isASCIIGeneric is the portable SWAR implementation behind IsASCII. No
// index is needed, so detection can be deferred: whole groups of words are
// OR-ed into one accumulator and tested once, four words per branch and
// then two, which is what turns the per-word test into a per-32-byte one.
// Inputs of 8..15 bytes skip the loops entirely — two overlapping words
// already cover them — and every 8+ byte input finishes with one
// overlapping word at n-8.
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
	if n < 2*swar.WordLen {
		return (swar.Load8(data, 0)|swar.Load8(data, n-swar.WordLen))&swar.HighBits == 0
	}
	i := 0
	for ; i+4*swar.WordLen <= n; i += 4 * swar.WordLen {
		w := data[i : i+4*swar.WordLen : i+4*swar.WordLen]
		acc := swar.Load8(w, 0) |
			swar.Load8(w, swar.WordLen) |
			swar.Load8(w, 2*swar.WordLen) |
			swar.Load8(w, 3*swar.WordLen)
		if acc&swar.HighBits != 0 {
			return false
		}
	}
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		w := data[i : i+2*swar.WordLen : i+2*swar.WordLen]
		if (swar.Load8(w, 0)|swar.Load8(w, swar.WordLen))&swar.HighBits != 0 {
			return false
		}
	}
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

// firstNonASCIIGeneric is the portable SWAR implementation behind
// FirstNonASCII: two words per branch, then one overlapping word at n-8.
func firstNonASCIIGeneric(data []byte) int {
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
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		w := data[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := swar.Load8(w, 0) & swar.HighBits
		m1 := swar.Load8(w, swar.WordLen) & swar.HighBits
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
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

// countNonASCIIGeneric is the portable SWAR implementation behind
// CountNonASCII: a popcount of the 0x80 lane bits per word, four words per
// pinned window. Every word has to be visited whatever happens, so unlike
// the first-match scans there is no branch to amortize — the win here is
// purely the bounds checks the window removes.
func countNonASCIIGeneric(data []byte) int {
	count := 0
	i := 0
	n := len(data)
	for ; i+4*swar.WordLen <= n; i += 4 * swar.WordLen {
		w := data[i : i+4*swar.WordLen : i+4*swar.WordLen]
		count += bits.OnesCount64(swar.Load8(w, 0)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, swar.WordLen)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, 2*swar.WordLen)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, 3*swar.WordLen)&swar.HighBits)
	}
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
