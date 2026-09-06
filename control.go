package utils

import (
	"github.com/gofiber/utils/v2/swar"
)

// noControlExemption is never a control lane, so exempting it is a no-op.
const noControlExemption = 0x80

// IndexControl returns the index of the first ASCII control byte in s (below
// 0x20 or DEL, the RFC 5234 CTL set), or -1. HTAB counts as a control byte;
// bytes >= 0x80 never match, unlike unicode.IsControl's C1 range.
func IndexControl[S byteSeq](s S) int {
	return scanControl(s, noControlExemption)
}

// IndexControlExceptTab is IndexControl with HTAB permitted: the bytes an
// RFC 9110 field value may not contain.
func IndexControlExceptTab[S byteSeq](s S) int {
	return scanControl(s, '\t')
}

// scanControl is the word scan behind both, with exempt masked out per word.
func scanControl[S byteSeq](s S, exempt byte) int {
	n := len(s)
	i := 0
	for ; i+16 <= n; i += 16 {
		w := s[i : i+16]
		m0 := controlLanes(swar.Load8(w, 0), exempt)
		m1 := controlLanes(swar.Load8(w, 8), exempt)
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + 8 + swar.FirstLane(m1)
		}
	}
	for ; i+8 <= n; i += 8 {
		if m := controlLanes(swar.Load8(s, i), exempt); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := controlLanes(swar.Load8(s, n-8), exempt); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := s[i]; (c < 0x20 || c == 0x7f) && c != exempt {
			return i
		}
	}
	return -1
}

// controlLanes flags the lanes of w below 0x20 or equal to DEL, other than
// exempt, exactly per lane; the biased lanes stay below 0xE0, so no carries.
func controlLanes(w uint64, exempt byte) uint64 {
	b := w & swar.LowSeven
	ctl := (^(b + (0x80-0x20)*swar.Ones) | (b + swar.Ones)) &^ w & swar.HighBits
	return ctl &^ swar.MatchByteMask(w, exempt)
}
