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
//
// The word loops pin each group of four words with a three-index reslice
// and address the words at constant offsets inside it. A two-index
// reslice at a variable offset costs a capacity check plus the
// four-instruction pointer clamp the compiler emits when the new
// capacity is not a known nonzero constant; the pinned form pays one
// check per group and no clamp, and the constant-offset accesses inside
// it are check-free. See swar.Load8 for the same note.

// WordLen is the SWAR word width in bytes, and the length from which these
// helpers are worth a call: below it the conversion is a table lookup per byte,
// cheaper inline. The case packages gate on it; the From converters require it.
const WordLen = swar.WordLen

// ToLowerFrom returns a lower-cased copy of src. first must be FirstUpperIndex's
// own result, with 0 <= first < len(src) and len(src) >= WordLen: a larger first
// silently mis-converts, so it must not be cached or derived elsewhere.
//
// src is only read, so it may be an unsafe view over immutable string memory;
// the result is freshly allocated and unaliased. strings.ToLower needs both.
func ToLowerFrom(src []byte, first int) []byte {
	n := len(src)
	dst := make([]byte, n)

	// Copy verbatim up to the word holding first, then convert from there.
	i := first &^ (WordLen - 1)
	copy(dst, src[:i])

	for ; i+4*WordLen <= n; i += 4 * WordLen {
		sw := src[i : i+4*WordLen : i+4*WordLen]
		d := dst[i : i+4*WordLen : i+4*WordLen]
		binary.LittleEndian.PutUint64(d[0:WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[0:WordLen])))
		binary.LittleEndian.PutUint64(d[WordLen:2*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[WordLen:2*WordLen])))
		binary.LittleEndian.PutUint64(d[2*WordLen:3*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[2*WordLen:3*WordLen])))
		binary.LittleEndian.PutUint64(d[3*WordLen:4*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[3*WordLen:4*WordLen])))
	}
	for ; i+WordLen <= n; i += WordLen {
		binary.LittleEndian.PutUint64(dst[i:i+WordLen:i+WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(src[i:i+WordLen:i+WordLen])))
	}
	if i == n {
		return dst
	}
	// One overlapping word finishes the tail. It re-reads src, so above the
	// aligned start it recomputes what the loops wrote, and below it folds over
	// the verbatim prefix -- harmless only because every byte under first is
	// unchanged by the fold, which is what makes first's upper bound matter.
	binary.LittleEndian.PutUint64(dst[n-WordLen:n:n], swar.ToLowerWord(binary.LittleEndian.Uint64(src[n-WordLen:n:n])))
	return dst
}

// ToUpperFrom mirrors ToLowerFrom, with first from FirstLowerIndex.
func ToUpperFrom(src []byte, first int) []byte {
	n := len(src)
	dst := make([]byte, n)

	// Copy verbatim up to the word holding first, then convert from there.
	i := first &^ (WordLen - 1)
	copy(dst, src[:i])

	for ; i+4*WordLen <= n; i += 4 * WordLen {
		sw := src[i : i+4*WordLen : i+4*WordLen]
		d := dst[i : i+4*WordLen : i+4*WordLen]
		binary.LittleEndian.PutUint64(d[0:WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[0:WordLen])))
		binary.LittleEndian.PutUint64(d[WordLen:2*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[WordLen:2*WordLen])))
		binary.LittleEndian.PutUint64(d[2*WordLen:3*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[2*WordLen:3*WordLen])))
		binary.LittleEndian.PutUint64(d[3*WordLen:4*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[3*WordLen:4*WordLen])))
	}
	for ; i+WordLen <= n; i += WordLen {
		binary.LittleEndian.PutUint64(dst[i:i+WordLen:i+WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(src[i:i+WordLen:i+WordLen])))
	}
	if i == n {
		return dst
	}
	// See ToLowerFrom for why the overlapping store is harmless.
	binary.LittleEndian.PutUint64(dst[n-WordLen:n:n], swar.ToUpperWord(binary.LittleEndian.Uint64(src[n-WordLen:n:n])))
	return dst
}

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
			w := b[i : i+32 : i+32]
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'A', 'Z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'A', 'Z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'A', 'Z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'A', 'Z')
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
		w := b[i : i+32 : i+32]
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'A', 'Z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'A', 'Z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'A', 'Z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'A', 'Z')
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

// FirstLowerIndex returns the index of the first ASCII lowercase byte in b,
// or -1 if b contains none. b is only ever read; it may be an unsafe view
// over immutable string memory.
func FirstLowerIndex(b []byte) int {
	n := len(b)
	i := 0
	if n >= 32 {
		// Structure and rationale mirror FirstUpperIndex.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[0:8]), 'a', 'z'); m != 0 {
			return swar.FirstLane(m)
		}
		for i = 8; i+32 <= n; i += 32 {
			w := b[i : i+32 : i+32]
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'a', 'z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'a', 'z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'a', 'z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'a', 'z')
			if m0|m1|m2|m3 != 0 {
				return firstSetLane32(i, m0, m1, m2, m3)
			}
		}
		if i == n {
			return -1
		}
		i = n - 32
		w := b[i : i+32 : i+32]
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'a', 'z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'a', 'z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'a', 'z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'a', 'z')
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

// ToLowerInPlace lower-cases every ASCII uppercase byte of b in place:
// four words per pinned group, then single words, then bytes.
func ToLowerInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i+32 <= n; i += 32 {
		w := b[i : i+32 : i+32]
		binary.LittleEndian.PutUint64(w[0:8], swar.ToLowerWord(binary.LittleEndian.Uint64(w[0:8])))
		binary.LittleEndian.PutUint64(w[8:16], swar.ToLowerWord(binary.LittleEndian.Uint64(w[8:16])))
		binary.LittleEndian.PutUint64(w[16:24], swar.ToLowerWord(binary.LittleEndian.Uint64(w[16:24])))
		binary.LittleEndian.PutUint64(w[24:32], swar.ToLowerWord(binary.LittleEndian.Uint64(w[24:32])))
	}
	for ; i+8 <= n; i += 8 {
		w := b[i : i+8 : i+8]
		binary.LittleEndian.PutUint64(w, swar.ToLowerWord(binary.LittleEndian.Uint64(w)))
	}
	// Finish byte-wise: an overlapping word here would partially overlap the
	// previous 8-byte store and stall on store-to-load forwarding.
	for ; i < n; i++ {
		b[i] = ToLowerTable[b[i]]
	}
}

// ToUpperInPlace upper-cases every ASCII lowercase byte of b in place; the
// loop shape mirrors ToLowerInPlace.
func ToUpperInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i+32 <= n; i += 32 {
		w := b[i : i+32 : i+32]
		binary.LittleEndian.PutUint64(w[0:8], swar.ToUpperWord(binary.LittleEndian.Uint64(w[0:8])))
		binary.LittleEndian.PutUint64(w[8:16], swar.ToUpperWord(binary.LittleEndian.Uint64(w[8:16])))
		binary.LittleEndian.PutUint64(w[16:24], swar.ToUpperWord(binary.LittleEndian.Uint64(w[16:24])))
		binary.LittleEndian.PutUint64(w[24:32], swar.ToUpperWord(binary.LittleEndian.Uint64(w[24:32])))
	}
	for ; i+8 <= n; i += 8 {
		w := b[i : i+8 : i+8]
		binary.LittleEndian.PutUint64(w, swar.ToUpperWord(binary.LittleEndian.Uint64(w)))
	}
	for ; i < n; i++ {
		b[i] = ToUpperTable[b[i]]
	}
}
