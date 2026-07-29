package utils

import (
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/simd"
	"github.com/gofiber/utils/v2/swar"
)

// quoteLanes and backslashLanes broadcast '"' and '\\' to every lane.
const (
	quoteLanes     = uint64('"') * swar.Ones
	backslashLanes = uint64('\\') * swar.Ones
)

// The dispatch prologues below gate on simd.MinLen — the measured length
// from which the AVX2 kernels beat these SWAR word loops — so the crossover
// has a single source of truth in the simd package. The simd.Accelerated()
// gate is a constant false off amd64, so the accelerated branches compile
// away entirely there.

// IsASCII reports whether s contains only ASCII bytes (no byte >= 0x80).
// It ORs four words per iteration where possible (no index is needed, so
// detection can be deferred to one test per 32 bytes), then two, then
// word-wise; inputs of 8+ bytes finish with one overlapping word at n-8,
// shorter ones byte-wise.
// On amd64 CPUs with AVX2, inputs of 32+ bytes dispatch to package simd
// instead.
//
// The unrolled loops pin their group of words with a reslice and load at
// constant offsets inside it. swar.Load8's own reslice is bounds-checked
// against the backing array, which the length-based loop condition does
// not imply, so a variable index costs a compare and a branch per load;
// pinning the window pays that once per group instead. This mirrors
// simd.isASCIIGeneric — the same scan, kept in step so a retune of either
// is a retune of both.
func IsASCII[S byteSeq](s S) bool {
	n := len(s)
	if n >= simd.MinLen && simd.Accelerated() {
		return simd.IsASCII(unsafeconv.Bytes(s))
	}
	if n >= 8 && n < 16 {
		return (swar.Load8(s, 0)|swar.Load8(s, n-8))&swar.HighBits == 0
	}
	i := 0
	for ; i+32 <= n; i += 32 {
		w := s[i : i+32]
		acc := swar.Load8(w, 0) | swar.Load8(w, 8) | swar.Load8(w, 16) | swar.Load8(w, 24)
		if acc&swar.HighBits != 0 {
			return false
		}
	}
	for ; i+16 <= n; i += 16 {
		w := s[i : i+16]
		if (swar.Load8(w, 0)|swar.Load8(w, 8))&swar.HighBits != 0 {
			return false
		}
	}
	for ; i+8 <= n; i += 8 {
		if swar.Load8(s, i)&swar.HighBits != 0 {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		return swar.Load8(s, n-8)&swar.HighBits == 0
	}
	for ; i < n; i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// IndexNonQuotable returns the index of the first byte of s that cannot
// appear verbatim inside an RFC 9110 quoted-string — that is, a byte
// matching c == '\\' || c == '"' || (c < 0x20 && c != '\t') || c == 0x7f —
// or -1 if every byte is quotable. qdtext is HTAB / SP / %x21 / %x23-5B /
// %x5D-7E / obs-text, so HTAB and bytes >= 0x80 are quotable. Returning the
// index (rather than a bool) lets callers copy the clean prefix and start
// escaping exactly at the offending byte. Inputs of 8+ bytes finish with
// one overlapping word at n-8, shorter ones byte-wise.
func IndexNonQuotable[S byteSeq](s S) int {
	n := len(s)
	i := 0
	// Two words per branch: the masks compute independently, and if only the
	// second word matched, its first set lane is still exact.
	for ; i+16 <= n; i += 16 {
		// Pinned window, as in IsASCII above: the two loads land at
		// constant offsets and skip their bounds checks.
		w := s[i : i+16]
		m0 := nonQuotableMask(swar.Load8(w, 0))
		m1 := nonQuotableMask(swar.Load8(w, 8))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + 8 + swar.FirstLane(m1)
		}
	}
	for ; i+8 <= n; i += 8 {
		if m := nonQuotableMask(swar.Load8(s, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := nonQuotableMask(swar.Load8(s, n-8)); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := s[i]; c == '\\' || c == '"' || (c < 0x20 && c != '\t') || c == 0x7f {
			return i
		}
	}
	return -1
}

// nonQuotableMask marks the lanes of w holding bytes that need RFC 9110
// quoted-string escaping. Like swar.ZeroLanes, the result is exact in and below
// the first marked lane, which is all the first-match scan above consumes.
func nonQuotableMask(w uint64) uint64 {
	// Controls (< 0x20) and DEL (0x7F) share one biased range test:
	// t := ((c & 0x7F) + 1) & 0x7F maps DEL to 0 and controls to 0x01..0x20,
	// so t <= 0x20 captures exactly both; lanes with the high bit set
	// (obs-text, always quotable) are excluded by the &^ w term. HTAB is
	// quotable qdtext, so its lanes are cleared with an exact match mask —
	// an approximate one could wrongly clear a control lane above a tab.
	t := ((w & swar.LowSeven) + swar.Ones) & swar.LowSeven
	ctl := ^(t + (0x80-0x21)*swar.Ones) &^ w & swar.HighBits &^ swar.MatchByteMask(w, '\t')
	return ctl | swar.ZeroLanes(w^quoteLanes) | swar.ZeroLanes(w^backslashLanes)
}
