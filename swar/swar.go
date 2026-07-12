// Package swar provides SWAR (SIMD within a register) primitives for
// processing eight ASCII bytes per uint64 word. The building blocks here are
// the ones the rest of gofiber/utils composes into case conversion,
// multi-needle scanning, and validation loops, and they are exported so that
// downstream packages (fiber itself, middleware) can fuse their own scans —
// e.g. finding a delimiter while simultaneously classifying the bytes before
// it — without re-deriving the bit tricks.
//
// Every operation yields per-lane results (exactly for the Match* masks,
// approximately for ZeroLanes, whose contract allows false positives above
// the first hit), so results are identical on little- and big-endian
// platforms. Load8 and Store8 always use little-endian lane order: lane 0
// (the least significant byte) is the lowest byte index, which is what
// FirstLane and LastLane assume. The empty-mask sentinels differ on
// purpose: FirstLane(0) == 8 lets forward scans step past the current word
// arithmetically, while LastLane(0) == -1 is the conventional not-found
// result for reverse scans.
//
// No unsafe is used anywhere in this package.
package swar

import (
	"math/bits"
)

const (
	// WordLen is the SWAR word width in bytes. Word loops should iterate
	// while i+WordLen <= len(s) and route shorter inputs to scalar code.
	WordLen = 8

	// Ones has 0x01 in every byte lane.
	Ones = 0x0101010101010101
	// HighBits has 0x80 in every byte lane; w&HighBits != 0 iff some byte
	// has its high bit set, and Match* masks equal HighBits iff every lane
	// matched.
	HighBits = 0x8080808080808080
	// LowSeven has 0x7F in every byte lane.
	LowSeven = 0x7f7f7f7f7f7f7f7f
)

// Load8 assembles s[i:i+8] into a little-endian uint64: s[i] becomes the
// least significant byte (lane 0). The caller must guarantee i+8 <= len(s);
// the reslice below turns a violation into a bounds panic. The 8-byte
// reslice pins the length so the constant-index reads are bounds-check free,
// and the compiler fuses them into a single 8-byte load on little-endian
// targets (verified for both the string and []byte instantiations).
func Load8[S ~string | ~[]byte](s S, i int) uint64 {
	w := s[i : i+8]
	return uint64(w[0]) |
		uint64(w[1])<<8 |
		uint64(w[2])<<16 |
		uint64(w[3])<<24 |
		uint64(w[4])<<32 |
		uint64(w[5])<<40 |
		uint64(w[6])<<48 |
		uint64(w[7])<<56
}

// Store8 writes w into b[i:i+8] in Load8's lane order: lane 0 (the least
// significant byte) lands at the lowest byte index, so a Load8/Store8 round
// trip is the identity. The caller must guarantee i+8 <= len(b); the
// reslice below turns a violation into a bounds panic. Like Load8, the
// constant-index writes are bounds-check free and fuse into a single 8-byte
// store on little-endian targets.
func Store8(b []byte, i int, w uint64) {
	d := b[i : i+8]
	d[0] = byte(w)
	d[1] = byte(w >> 8)
	d[2] = byte(w >> 16)
	d[3] = byte(w >> 24)
	d[4] = byte(w >> 32)
	d[5] = byte(w >> 40)
	d[6] = byte(w >> 48)
	d[7] = byte(w >> 56)
}

// Broadcast returns a word with c in every byte lane, the needle form that
// ZeroLanes scans consume. Hoist the result out of word loops; it compiles
// to a single multiply.
func Broadcast(c byte) uint64 {
	return uint64(c) * Ones
}

// ToLowerWord lower-cases every 'A'..'Z' lane in w. All other bytes,
// including those >= 0x80, pass through unchanged — bit-identical to the
// ASCII case-conversion tables used elsewhere in this module.
func ToLowerWord(w uint64) uint64 {
	return w | MatchRangeMask(w, 'A', 'Z')>>2 // 0x80>>2 == 0x20: set bit 5
}

// ToUpperWord upper-cases every 'a'..'z' lane in w. All other bytes,
// including those >= 0x80, pass through unchanged.
func ToUpperWord(w uint64) uint64 {
	return w &^ (MatchRangeMask(w, 'a', 'z') >> 2) // clear bit 5 on matched lanes
}

// ZeroLanes returns a mask whose lowest set lane is the lowest zero-byte
// lane of x; higher lanes may carry false positives (borrows propagate
// strictly upward from true zero lanes), and the mask is zero iff x has no
// zero byte. It is two ops cheaper than MatchByteMask, which makes it the
// right primitive for first-match scans, typically as
// ZeroLanes(w ^ Broadcast(c)) with the broadcast hoisted out of the word
// loop. Use MatchByteMask when every lane must be exact.
func ZeroLanes(x uint64) uint64 {
	return (x - Ones) &^ x & HighBits
}

// MatchByteMask returns a word with 0x80 set in exactly the lanes of w whose
// byte equals c, and zero in all other lanes. The result is per-lane exact
// (no false positives in lanes above a match), so it is safe to feed to
// LastLane or to arbitrary lane masking, not just first-match scans; for
// pure first-match scanning ZeroLanes is cheaper.
func MatchByteMask(w uint64, c byte) uint64 {
	x := w ^ (uint64(c) * Ones)
	// Setting the high bit before the per-lane subtraction keeps every lane
	// >= 0x80, so borrows cannot cross lanes and the test stays exact:
	// after the subtraction the high bit is clear only where the lane was 0,
	// and OR-ing x back in rejects lanes that had their own high bit set.
	return ^(((x | HighBits) - Ones) | x) & HighBits
}

// MatchRangeMask returns a word with 0x80 set in exactly the lanes of w
// whose byte is within [lo, hi]. It requires lo <= hi <= 0x7F; lanes with
// the high bit set (bytes >= 0x80) never match.
func MatchRangeMask(w uint64, lo, hi byte) uint64 {
	b := w & LowSeven
	// After masking, every lane is <= 0x7F, so the biased additions below
	// cannot carry into a neighboring lane.
	ge := b + (0x80-uint64(lo))*Ones   // high bit set where lane >= lo
	gt := b + (0x80-uint64(hi)-1)*Ones // high bit set where lane > hi
	return ge &^ gt &^ w & HighBits
}

// FirstLane returns the index (0..7) of the lowest-addressed set lane in
// mask — the smallest byte index, given Load8's lane order — or 8 if mask
// is zero. mask must contain only 0x80 lane markers, as produced by the
// Match* functions.
func FirstLane(mask uint64) int {
	return bits.TrailingZeros64(mask) >> 3
}

// LastLane returns the index (0..7) of the highest-addressed set lane in
// mask, or -1 if mask is zero. mask must contain only 0x80 lane markers.
func LastLane(mask uint64) int {
	return (63 - bits.LeadingZeros64(mask)) >> 3
}
