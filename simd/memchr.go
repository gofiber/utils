// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"

	"github.com/gofiber/utils/v2/swar"
)

// The portable fallbacks below are composed from the swar package's
// primitives rather than re-deriving the bit tricks; see swar's package
// documentation for the ZeroLanes contract the first-match scans rely on
// (false positives only above the first true match, so the first set lane
// is always exact).

// Memchr2 returns the index of the first occurrence in haystack of either
// needle1 or needle2, or -1 if neither is present. Both needles are checked
// in the same pass, so it costs the same as a single-byte scan.
//
// There is deliberately no single-needle Memchr: use bytes.IndexByte, which
// the Go runtime already accelerates on all major architectures.
func Memchr2(haystack []byte, needle1, needle2 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr2Impl(haystack, needle1, needle2)
}

// Memchr3 returns the index of the first occurrence in haystack of needle1,
// needle2, or needle3, or -1 if none is present. All three needles are
// checked in the same pass.
func Memchr3(haystack []byte, needle1, needle2, needle3 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr3Impl(haystack, needle1, needle2, needle3)
}

// MemchrPair returns the first index i such that haystack[i] == byte1 and
// haystack[i+offset] == byte2, or -1 if no such position exists. offset must
// be >= 0; a negative offset returns -1.
//
// Requiring two bytes at a fixed distance makes the scan far more selective
// than a single-byte search, which is what makes it a good substring
// prefilter: false candidates need both bytes at exactly the right distance.
// Memmem uses it for short needles.
func MemchrPair(haystack []byte, byte1, byte2 byte, offset int) int {
	if offset < 0 || len(haystack) <= offset {
		return -1
	}
	if offset == 0 {
		// Both bytes would occupy the same position, so they must be equal.
		if byte1 != byte2 {
			return -1
		}
		return bytes.IndexByte(haystack, byte1)
	}
	return memchrPairImpl(haystack, byte1, byte2, offset)
}

// memchr2Generic is the portable SWAR fallback for Memchr2: a ZeroLanes
// first-match scan over both needle broadcasts, finishing 8+ byte inputs
// with one overlapping word at n-8 (already-scanned lanes are known
// non-matching, so the first set lane falls in the new bytes).
func memchr2Generic(haystack []byte, needle1, needle2 byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(needle1)
	bc2 := swar.Broadcast(needle2)
	i := 0
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		// Pinned window: see the note above isASCIIGeneric in ascii.go for
		// why the loads use constant offsets into it.
		win := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		w0, w1 := swar.Load8(win, 0), swar.Load8(win, swar.WordLen)
		m0 := swar.ZeroLanes(w0^bc1) | swar.ZeroLanes(w0^bc2)
		m1 := swar.ZeroLanes(w1^bc1) | swar.ZeroLanes(w1^bc2)
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		w := swar.Load8(haystack, i)
		if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	w := swar.Load8(haystack, n-swar.WordLen)
	if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// memchr3Generic is the portable SWAR fallback for Memchr3; see
// memchr2Generic for the scan shape.
func memchr3Generic(haystack []byte, needle1, needle2, needle3 byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 || haystack[i] == needle3 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(needle1)
	bc2 := swar.Broadcast(needle2)
	bc3 := swar.Broadcast(needle3)
	i := 0
	for ; i+2*swar.WordLen <= n; i += 2 * swar.WordLen {
		win := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		w0, w1 := swar.Load8(win, 0), swar.Load8(win, swar.WordLen)
		m0 := swar.ZeroLanes(w0^bc1) | swar.ZeroLanes(w0^bc2) | swar.ZeroLanes(w0^bc3)
		m1 := swar.ZeroLanes(w1^bc1) | swar.ZeroLanes(w1^bc2) | swar.ZeroLanes(w1^bc3)
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		w := swar.Load8(haystack, i)
		if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2) | swar.ZeroLanes(w^bc3); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	w := swar.Load8(haystack, n-swar.WordLen)
	if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2) | swar.ZeroLanes(w^bc3); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// memchrPairGeneric is the portable SWAR fallback for MemchrPair. It leans
// on ZeroLanes flagging every true zero lane (part of its contract), so
// ANDing the two per-position masks can never drop a real pair; what the
// AND does lose is lane exactness — a flagged lane may be a false
// candidate sitting strictly below the first true pair — so every
// candidate is re-verified scalar and refuted ones are stepped past.
// ZeroLanes plus this rare verify step beats a pair of exact
// swar.MatchByteMask masks on the no-candidate fast path, which is the
// common case. The caller guarantees offset >= 1 and
// len(haystack) > offset.
func memchrPairGeneric(haystack []byte, byte1, byte2 byte, offset int) int {
	n := len(haystack)
	if n < swar.WordLen+offset {
		for i := 0; i+offset < n; i++ {
			if haystack[i] == byte1 && haystack[i+offset] == byte2 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(byte1)
	bc2 := swar.Broadcast(byte2)
	i := 0
	for ; i+swar.WordLen+offset <= n; i += swar.WordLen {
		cand := swar.ZeroLanes(swar.Load8(haystack, i)^bc1) &
			swar.ZeroLanes(swar.Load8(haystack, i+offset)^bc2)
		for cand != 0 {
			pos := i + swar.FirstLane(cand)
			if haystack[pos] == byte1 && haystack[pos+offset] == byte2 {
				return pos
			}
			cand &= cand - 1
		}
	}
	for ; i+offset < n; i++ {
		if haystack[i] == byte1 && haystack[i+offset] == byte2 {
			return i
		}
	}
	return -1
}
