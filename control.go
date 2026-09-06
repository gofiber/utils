package utils

import (
	"github.com/gofiber/utils/v2/swar"
)

// noControlExemption is the byte passed to indexControl when no control
// byte is exempt: it is never a control lane, so masking it out is a no-op.
const noControlExemption = 0x80

// IndexControl returns the index of the first ASCII control byte in s — a
// byte below 0x20 or equal to 0x7F (DEL), the CTL set of RFC 5234 — or -1
// if s contains none. HTAB (0x09) is a control byte here; IndexControlExceptTab
// exempts it for HTTP field values, where HTAB is legal whitespace. Bytes
// >= 0x80 never match, so unlike unicode.IsControl the C1 range hidden in
// UTF-8 sequences is not flagged; this is the byte-level check that header
// values, request IDs, and log fields need before they are echoed.
func IndexControl[S byteSeq](s S) int {
	return indexControl(s, noControlExemption)
}

// IndexControlExceptTab is IndexControl with HTAB permitted: it returns the
// index of the first byte below 0x20 other than 0x09, or equal to 0x7F, or
// -1 if there is none. That is the byte set an RFC 9110 field value may
// not contain (field-content is VCHAR, SP, HTAB, and obs-text).
func IndexControlExceptTab[S byteSeq](s S) int {
	return indexControl(s, '\t')
}

// indexControl scans for control bytes with exempt masked out of every
// word: two words per branch, then one, then one overlapping word at n-8;
// inputs shorter than a word are checked byte-wise.
func indexControl[S byteSeq](s S, exempt byte) int {
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

// controlLanes flags the lanes of w holding control bytes (below 0x20 or
// DEL) other than exempt, exactly per lane. Lanes below 0x20 are the ones
// a bias of 0x60 does not carry into bit 7, DEL is the one 0x7F lane that
// a bias of 1 does, and the &^ w term drops lanes with their own high bit
// set. No lane can carry into its neighbor: the biased lanes stay below
// 0xE0.
func controlLanes(w uint64, exempt byte) uint64 {
	b := w & swar.LowSeven
	ctl := (^(b + (0x80-0x20)*swar.Ones) | (b + swar.Ones)) &^ w & swar.HighBits
	return ctl &^ swar.MatchByteMask(w, exempt)
}
