package utils

import (
	"github.com/gofiber/utils/v2/swar"
)

// quoteLanes and backslashLanes broadcast '"' and '\\' to every lane.
const (
	quoteLanes     = uint64('"') * swar.Ones
	backslashLanes = uint64('\\') * swar.Ones
)

// IsASCII reports whether s contains only ASCII bytes (no byte >= 0x80).
// It ORs four words per iteration where possible (no index is needed, so
// detection can be deferred to one test per 32 bytes), then tests word-wise;
// inputs of 8+ bytes finish with one overlapping word at n-8, shorter ones
// byte-wise.
func IsASCII[S byteSeq](s S) bool {
	n := len(s)
	i := 0
	for ; i+32 <= n; i += 32 {
		acc := swar.Load8(s, i) | swar.Load8(s, i+8) | swar.Load8(s, i+16) | swar.Load8(s, i+24)
		if acc&swar.HighBits != 0 {
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
// quoted-string escaping. Like swar.ZeroLanes, the result is exact in and below
// the first marked lane, which is all the first-match scan above consumes.
func nonQuotableMask(w uint64) uint64 {
	// Controls (< 0x20) and DEL (0x7F) share one biased range test:
	// t := ((c & 0x7F) + 1) & 0x7F maps DEL to 0 and controls to 0x01..0x20,
	// so t <= 0x20 captures exactly both; lanes with the high bit set
	// (obs-text, always quotable) are excluded by the &^ w term.
	t := ((w & swar.LowSeven) + swar.Ones) & swar.LowSeven
	ctl := ^(t + (0x80-0x21)*swar.Ones) &^ w & swar.HighBits
	return ctl | swar.ZeroLanes(w^quoteLanes) | swar.ZeroLanes(w^backslashLanes)
}
