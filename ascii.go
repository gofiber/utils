package utils

import (
	"github.com/gofiber/utils/v2/swar"
)

const (
	// asciiHighBits has 0x80 in every byte lane; a word ANDed with it is
	// nonzero iff some byte has its high bit set (i.e. is not ASCII).
	asciiHighBits = 0x8080808080808080
	// asciiOnes has 0x01 in every byte lane.
	asciiOnes = 0x0101010101010101
	// asciiLowSeven has 0x7F in every byte lane.
	asciiLowSeven = 0x7f7f7f7f7f7f7f7f
	// quoteLanes and backslashLanes broadcast '"' and '\\' to every lane.
	quoteLanes     = 0x2222222222222222
	backslashLanes = 0x5c5c5c5c5c5c5c5c
)

// IsASCII reports whether s contains only ASCII bytes (no byte >= 0x80).
// It tests 8 bytes per iteration; inputs of 8+ bytes finish with one
// overlapping word at n-8, shorter ones byte-wise.
func IsASCII[S byteSeq](s S) bool {
	n := len(s)
	i := 0
	for ; i+8 <= n; i += 8 {
		if swar.Load8(s, i)&asciiHighBits != 0 {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		return swar.Load8(s, n-8)&asciiHighBits == 0
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
// matching c == '\\' || c == '"' || c < 0x20 || c == 0x7f — or -1 if every
// byte is quotable. Returning the index (rather than a bool) lets callers
// copy the clean prefix and start escaping exactly at the offending byte.
// Bytes >= 0x80 (obs-text) are quotable. Inputs of 8+ bytes finish with one
// overlapping word at n-8, shorter ones byte-wise.
func IndexNonQuotable[S byteSeq](s S) int {
	n := len(s)
	i := 0
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
		if c := s[i]; c == '\\' || c == '"' || c < 0x20 || c == 0x7f {
			return i
		}
	}
	return -1
}

// nonQuotableMask marks the lanes of w holding bytes that need RFC 9110
// quoted-string escaping. The result is exact in and below the first marked
// lane; higher lanes may carry false positives (the quote/backslash tests use
// the cheaper zero-detect form, whose borrows only corrupt lanes above a true
// match), which is fine for the first-match scans in IndexNonQuotable.
func nonQuotableMask(w uint64) uint64 {
	// Controls (< 0x20) and DEL (0x7F) share one biased range test:
	// t := ((c & 0x7F) + 1) & 0x7F maps DEL to 0 and controls to 0x01..0x20,
	// so t <= 0x20 captures exactly both; lanes with the high bit set
	// (obs-text, always quotable) are excluded by the &^ w term.
	t := ((w & asciiLowSeven) + asciiOnes) & asciiLowSeven
	ctl := ^(t + (0x80-0x21)*asciiOnes) &^ w & asciiHighBits
	// '"' and '\\' via zero-byte detection on the XORed word.
	q := w ^ quoteLanes
	bs := w ^ backslashLanes
	quotes := ((q-asciiOnes)&^q | (bs-asciiOnes)&^bs) & asciiHighBits
	return ctl | quotes
}
