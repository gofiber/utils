package utils

import (
	"unicode/utf8"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

// ltLanes, gtLanes, and ampLanes broadcast the HTML-escaped characters to
// every lane; quoteLanes and backslashLanes live in ascii.go.
const (
	ltLanes  = uint64('<') * swar.Ones
	gtLanes  = uint64('>') * swar.Ones
	ampLanes = uint64('&') * swar.Ones
)

// jsonHexDigits is the escape alphabet for \u00XX sequences, matching
// encoding/json.
const jsonHexDigits = "0123456789abcdef"

// jsonSafeTable flags the ASCII bytes that pass through a JSON string
// encoder verbatim: 0x20..0x7E minus '"', '\\', and the HTML-escaped
// '<', '>', '&'. Bytes >= 0x80 are false so the scalar loop routes them
// to the UTF-8 path.
var jsonSafeTable = buildJSONSafeTable()

func buildJSONSafeTable() [256]bool {
	var t [256]bool
	for c := byte(0x20); c < utf8.RuneSelf; c++ {
		t[c] = c != '"' && c != '\\' && c != '<' && c != '>' && c != '&'
	}
	return t
}

// AppendJSONString appends s to dst as a double-quoted JSON string and
// returns the extended slice. The output is byte-identical to
// encoding/json.Marshal of the same string value, including its default
// HTML escaping: '"' and '\\' get backslash escapes; control bytes below
// 0x20 become \b, \f, \n, \r, \t, or \u00XX; '<', '>', '&' become \u00XX; each
// invalid UTF-8 byte becomes the six literal characters \ufffd; and the
// line separators U+2028/U+2029 become the \u2028 and \u2029 escapes.
// All other bytes, including multi-byte UTF-8 sequences, are copied
// verbatim. dst must not alias s: the output is longer than the input (at
// minimum by the surrounding quotes), so in-place encoding is impossible
// and an aliased dst would overwrite bytes before they are read. Clean
// spans are located with a SWAR scan and copied wholesale, so typical log
// or header values cost one scan and one copy — with none of
// encoding/json.Marshal's reflection or allocation.
func AppendJSONString[S byteSeq](dst []byte, s S) []byte {
	bs := unsafeconv.Bytes(s)
	n := len(bs)
	dst = append(dst, '"')
	start, i := 0, 0
	for i < n {
		// Skip to the next byte needing attention. Like nonQuotableMask,
		// the combined mask is exact in and below its first set lane, which
		// is all this first-match scan consumes.
		if i+8 <= n {
			m := jsonUnsafeMask(swar.Load8(bs, i))
			if m == 0 {
				i += 8
				continue
			}
			i += swar.FirstLane(m)
		} else if c := bs[i]; jsonSafeTable[c] {
			i++
			continue
		}
		c := bs[i]
		if c < utf8.RuneSelf {
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\')
			switch c {
			case '"', '\\':
				dst = append(dst, c)
			case '\b':
				dst = append(dst, 'b')
			case '\f':
				dst = append(dst, 'f')
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				// Other control bytes and the HTML characters '<', '>', '&'.
				dst = append(dst, 'u', '0', '0', jsonHexDigits[c>>4], jsonHexDigits[c&0x0F])
			}
			i++
			start = i
			continue
		}
		r, size := utf8.DecodeRune(bs[i:])
		switch {
		case r == utf8.RuneError && size == 1:
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\', 'u', 'f', 'f', 'f', 'd')
			i++
			start = i
		case r == '\u2028' || r == '\u2029':
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\', 'u', '2', '0', '2', byte('8'+r-'\u2028'))
			i += size
			start = i
		default:
			i += size
		}
	}
	dst = append(dst, bs[start:]...)
	return append(dst, '"')
}

// jsonUnsafeMask flags the lanes of w holding bytes a JSON string encoder
// cannot copy verbatim: controls below 0x20 (exact range mask), the five
// escaped ASCII characters (ZeroLanes needle matches, exact in and below
// their first hit), and all bytes >= 0x80, which the caller inspects as
// UTF-8.
func jsonUnsafeMask(w uint64) uint64 {
	return swar.MatchRangeMask(w, 0x00, 0x1F) | w&swar.HighBits |
		swar.ZeroLanes(w^quoteLanes) | swar.ZeroLanes(w^backslashLanes) |
		swar.ZeroLanes(w^ltLanes) | swar.ZeroLanes(w^gtLanes) | swar.ZeroLanes(w^ampLanes)
}
