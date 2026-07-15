package caseconv

import (
	"encoding/binary"

	"github.com/gofiber/utils/v2/swar"
)

// Word-at-a-time ASCII case conversion built on the public SWAR primitives
// in github.com/gofiber/utils/v2/swar. Byte order inside a word does not
// matter for the case folds themselves (every operation is per byte lane),
// so these helpers are correct on both little- and big-endian platforms.
//
// Words move through encoding/binary rather than swar.Load8/swar.Store8 on
// purpose: the binary.LittleEndian forms measure ~9% faster in the in-place
// word loops here (benchstat -count=10, Go 1.25, Apple M2 Pro). The lane
// order is identical, so swar words and binary.LittleEndian loads/stores
// can be mixed freely.

// WordLen is the SWAR word width in bytes, re-exported for the case
// conversion callers in the strings and bytes packages: they route inputs
// shorter than WordLen to byte-wise paths and align offsets they pass (such
// as ToLowerCopy's from) with WordLen-derived masks.
const WordLen = swar.WordLen

// FirstUpperIndex returns the index of the first ASCII uppercase byte in b,
// or -1 if b contains none. b is only ever read; it may be an unsafe view
// over immutable string memory.
func FirstUpperIndex(b []byte) int {
	n := len(b)
	i := 0
	if n >= 32 {
		// Check the first word alone: canonical mixed-case values
		// ("Content-Type") change case right at the start, and that path
		// must stay one mask.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[0:8]), 'A', 'Z'); m != 0 {
			return swar.FirstLane(m)
		}
		// Then four words per branch: the masks compute independently, so
		// the common no-uppercase case pays one test per 32 bytes. The
		// MatchRangeMask computations stay inline so their constant ranges
		// fold; only the cold hit-resolution ladder is shared, and it
		// inlines (the no-match loop never reaches it).
		for i = 8; i+32 <= n; i += 32 {
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'A', 'Z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+8:i+16]), 'A', 'Z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+16:i+24]), 'A', 'Z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+24:i+32]), 'A', 'Z')
			if m0|m1|m2|m3 != 0 {
				return firstSetLane32(i, m0, m1, m2, m3)
			}
		}
		if i == n {
			return -1
		}
		// Finish with one overlapping block at n-32. MatchRangeMask is
		// per-lane exact, so re-scanned lanes are known non-matching and the
		// first set lane always falls in the new bytes.
		i = n - 32
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'A', 'Z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+8:i+16]), 'A', 'Z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+16:i+24]), 'A', 'Z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+24:i+32]), 'A', 'Z')
		if m0|m1|m2|m3 != 0 {
			return firstSetLane32(i, m0, m1, m2, m3)
		}
		return -1
	}
	for ; i+8 <= n; i += 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'A', 'Z'); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		// Scan the tail with one overlapping word; lanes that were already
		// scanned are known to be zero, so FirstLane lands in the new bytes.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[n-8:n]), 'A', 'Z'); m != 0 {
			return n - 8 + swar.FirstLane(m)
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
// or -1 if b contains none. b is only ever read; it may be an unsafe view
// over immutable string memory.
// firstSetLane32 resolves a hit inside a 32-byte block: given the block's
// four per-word masks, it returns the index of the first marked byte
// relative to base, the block's starting offset. It runs at most once per
// scan — only after a block's combined mask tested non-zero — and is small
// enough to inline, so the hot no-match loops pay nothing for the sharing.
func firstSetLane32(base int, m0, m1, m2, m3 uint64) int {
	if m0 != 0 {
		return base + swar.FirstLane(m0)
	}
	if m1 != 0 {
		return base + 8 + swar.FirstLane(m1)
	}
	if m2 != 0 {
		return base + 16 + swar.FirstLane(m2)
	}
	return base + 24 + swar.FirstLane(m3)
}

func FirstLowerIndex(b []byte) int {
	n := len(b)
	i := 0
	if n >= 32 {
		// Structure and rationale mirror FirstUpperIndex.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[0:8]), 'a', 'z'); m != 0 {
			return swar.FirstLane(m)
		}
		for i = 8; i+32 <= n; i += 32 {
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'a', 'z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+8:i+16]), 'a', 'z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+16:i+24]), 'a', 'z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+24:i+32]), 'a', 'z')
			if m0|m1|m2|m3 != 0 {
				return firstSetLane32(i, m0, m1, m2, m3)
			}
		}
		if i == n {
			return -1
		}
		i = n - 32
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'a', 'z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+8:i+16]), 'a', 'z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+16:i+24]), 'a', 'z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i+24:i+32]), 'a', 'z')
		if m0|m1|m2|m3 != 0 {
			return firstSetLane32(i, m0, m1, m2, m3)
		}
		return -1
	}
	for ; i+8 <= n; i += 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'a', 'z'); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[n-8:n]), 'a', 'z'); m != 0 {
			return n - 8 + swar.FirstLane(m)
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

// ToLowerInPlace lower-cases every ASCII uppercase byte of b in place.
func ToLowerInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(b[i:i+8], swar.ToLowerWord(binary.LittleEndian.Uint64(b[i:i+8])))
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
		binary.LittleEndian.PutUint64(b[i:i+8], swar.ToUpperWord(binary.LittleEndian.Uint64(b[i:i+8])))
	}
	for ; i < n; i++ {
		b[i] = ToUpperTable[b[i]]
	}
}

// ToLowerCopy writes the lower-cased content of src into dst starting at
// byte offset from. Bytes before from must already be present in dst, and
// len(dst) must equal len(src). src is only ever read — the strings package
// passes an unsafe view over a string's immutable backing memory — so no
// code path here may write through src.
//
// from must be a multiple of WordLen and at or below the index of the first
// byte that lower-casing changes: when len(src) >= WordLen the overlapping
// tail store rewrites dst[n-WordLen:from) with case-converted bytes, which is
// only a no-op under that precondition.
//
// Callers currently invoke this only with len(src) >= WordLen; the byte-wise
// tail below keeps the helper correct standalone for shorter inputs.
func ToLowerCopy(dst, src []byte, from int) {
	n := len(src)
	i := from
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:i+8], swar.ToLowerWord(binary.LittleEndian.Uint64(src[i:i+8])))
	}
	if i == n {
		return
	}
	if n >= 8 {
		binary.LittleEndian.PutUint64(dst[n-8:n], swar.ToLowerWord(binary.LittleEndian.Uint64(src[n-8:n])))
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
// byte that upper-casing changes; see ToLowerCopy for why, including the
// note on the byte-wise tail and the requirement that src is never written.
func ToUpperCopy(dst, src []byte, from int) {
	n := len(src)
	i := from
	for ; i+8 <= n; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:i+8], swar.ToUpperWord(binary.LittleEndian.Uint64(src[i:i+8])))
	}
	if i == n {
		return
	}
	if n >= 8 {
		binary.LittleEndian.PutUint64(dst[n-8:n], swar.ToUpperWord(binary.LittleEndian.Uint64(src[n-8:n])))
		return
	}
	for ; i < n; i++ {
		dst[i] = ToUpperTable[src[i]]
	}
}
