package caseconv

import (
	"encoding/binary"
	"math/bits"
)

// SWAR (SIMD within a register) helpers that process eight ASCII bytes per
// uint64 word. Byte order inside the word does not matter: every operation is
// carried out independently per byte lane, so the helpers are correct on both
// little- and big-endian platforms.
const (
	wordOnes  = 0x0101010101010101
	wordHighs = 0x8080808080808080
)

// WordLen is the SWAR word width in bytes. It is the single source of truth
// for the word size assumed by every helper in this file: callers must route
// inputs shorter than WordLen to a byte-wise path and align offsets they pass
// (such as ToLowerCopy's from) with WordLen-derived masks.
const WordLen = 8

// UpperMask returns a word with 0x80 set in every byte lane of v that holds
// an ASCII uppercase letter ('A'..'Z') and zero in all other lanes.
func UpperMask(v uint64) uint64 {
	w := v &^ uint64(wordHighs)
	// After masking, every lane is <= 0x7F, so the additions below cannot
	// carry into a neighbouring lane.
	geA := w + (0x80-'A')*wordOnes // high bit set where lane >= 'A'
	gtZ := w + (0x80-'Z'-1)*wordOnes
	// Lanes with the original high bit set are not ASCII letters.
	return geA &^ gtZ &^ v & wordHighs
}

// LowerMask returns a word with 0x80 set in every byte lane of v that holds
// an ASCII lowercase letter ('a'..'z') and zero in all other lanes.
func LowerMask(v uint64) uint64 {
	w := v &^ uint64(wordHighs)
	geA := w + (0x80-'a')*wordOnes
	gtZ := w + (0x80-'z'-1)*wordOnes
	return geA &^ gtZ &^ v & wordHighs
}

// ToUpperWord upper-cases every ASCII lowercase byte lane in v.
func ToUpperWord(v uint64) uint64 {
	// 0x80 >> 2 == 0x20: toggling bit 5 flips the case of an ASCII letter.
	return v ^ LowerMask(v)>>2
}

// ToLowerWord lower-cases every ASCII uppercase byte lane in v.
func ToLowerWord(v uint64) uint64 {
	return v ^ UpperMask(v)>>2
}

// FirstUpperIndex returns the index of the first ASCII uppercase byte in b,
// or -1 if b contains none.
func FirstUpperIndex(b []byte) int {
	n := len(b)
	i := 0
	for ; i+8 <= n; i += 8 {
		if m := UpperMask(binary.LittleEndian.Uint64(b[i : i+8])); m != 0 {
			return i + firstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		// Scan the tail with one overlapping word; lanes that were already
		// scanned are known to be zero, so firstLane lands in the new bytes.
		if m := UpperMask(binary.LittleEndian.Uint64(b[n-8 : n])); m != 0 {
			return n - 8 + firstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := b[i]; c >= 'A' && c <= 'Z' {
			return i
		}
	}
	return -1
}

// FirstLowerIndex returns the index of the first ASCII lowercase byte in b,
// or -1 if b contains none.
func FirstLowerIndex(b []byte) int {
	n := len(b)
	i := 0
	for ; i+8 <= n; i += 8 {
		if m := LowerMask(binary.LittleEndian.Uint64(b[i : i+8])); m != 0 {
			return i + firstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := LowerMask(binary.LittleEndian.Uint64(b[n-8 : n])); m != 0 {
			return n - 8 + firstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := b[i]; c >= 'a' && c <= 'z' {
			return i
		}
	}
	return -1
}

// firstLane returns the index (0-7) of the lowest byte lane whose high bit is
// set in mask. mask must be non-zero and contain only 0x80 lane markers.
func firstLane(mask uint64) int {
	// The word was loaded little-endian, so the lowest set lane corresponds
	// to the smallest byte index regardless of platform byte order.
	return bits.TrailingZeros64(mask) >> 3
}

// ToLowerInPlace lower-cases every ASCII uppercase byte of b in place.
func ToLowerInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(b[i:i+8], ToLowerWord(binary.LittleEndian.Uint64(b[i:i+8])))
	}
	// Finish byte-wise: an overlapping word here would partially overlap the
	// previous 8-byte store and stall on store-to-load forwarding.
	for ; i < n; i++ {
		b[i] = ToLowerTable[b[i]]
	}
}

// ToUpperInPlace upper-cases every ASCII lowercase byte of b in place.
func ToUpperInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(b[i:i+8], ToUpperWord(binary.LittleEndian.Uint64(b[i:i+8])))
	}
	for ; i < n; i++ {
		b[i] = ToUpperTable[b[i]]
	}
}

// ToLowerCopy writes the lower-cased content of src into dst starting at
// byte offset from. Bytes before from must already be present in dst, and
// len(dst) must equal len(src).
//
// from must be a multiple of WordLen and at or below the index of the first
// byte that lower-casing changes: when len(src) >= WordLen the overlapping
// tail store rewrites dst[n-WordLen:from) with case-converted bytes, which is
// only a no-op under that precondition.
func ToLowerCopy(dst, src []byte, from int) {
	n := len(src)
	i := from
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:i+8], ToLowerWord(binary.LittleEndian.Uint64(src[i:i+8])))
	}
	if i == n {
		return
	}
	if n >= 8 {
		binary.LittleEndian.PutUint64(dst[n-8:n], ToLowerWord(binary.LittleEndian.Uint64(src[n-8:n])))
		return
	}
	for ; i < n; i++ {
		dst[i] = ToLowerTable[src[i]]
	}
}

// ToUpperCopy writes the upper-cased content of src into dst starting at
// byte offset from. Bytes before from must already be present in dst, and
// len(dst) must equal len(src).
//
// from must be a multiple of WordLen and at or below the index of the first
// byte that upper-casing changes; see ToLowerCopy for why.
func ToUpperCopy(dst, src []byte, from int) {
	n := len(src)
	i := from
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:i+8], ToUpperWord(binary.LittleEndian.Uint64(src[i:i+8])))
	}
	if i == n {
		return
	}
	if n >= 8 {
		binary.LittleEndian.PutUint64(dst[n-8:n], ToUpperWord(binary.LittleEndian.Uint64(src[n-8:n])))
		return
	}
	for ; i < n; i++ {
		dst[i] = ToUpperTable[src[i]]
	}
}
