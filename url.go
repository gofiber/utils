package utils

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// escapeMode selects between net/url's query-component and path-segment
// escaping dialects: only the query dialect maps '+' to space.
type escapeMode uint8

const (
	escapeQuery escapeMode = iota
	escapePath
)

// upperHexDigits is the escape alphabet used by the Append*Escape helpers,
// matching net/url.
const upperHexDigits = "0123456789ABCDEF"

// unhexTable maps an ASCII byte to its hexadecimal digit value, or 0xFF for
// bytes that are not hex digits.
var unhexTable = buildUnhexTable()

func buildUnhexTable() [256]byte {
	var t [256]byte
	for i := range t {
		t[i] = 0xFF
	}
	for c := byte('0'); c <= '9'; c++ {
		t[c] = c - '0'
	}
	for c := byte('a'); c <= 'f'; c++ {
		t[c] = c - 'a' + 10
	}
	for c := byte('A'); c <= 'F'; c++ {
		t[c] = c - 'A' + 10
	}
	return t
}

// queryNoEscapeTable and pathNoEscapeTable flag the bytes net/url leaves
// verbatim in query components and path segments respectively. They are
// derived from net/url's shouldEscape rules and pinned to them by
// exhaustive per-byte tests. Path segments additionally keep the sub-delims
// net/url's encodePathSegment mode does not escape.
var (
	queryNoEscapeTable = buildNoEscapeTable("")
	pathNoEscapeTable  = buildNoEscapeTable("$&+:=@")
)

// buildNoEscapeTable flags the RFC 3986 unreserved bytes plus alsoSafe as
// not needing escape.
func buildNoEscapeTable(alsoSafe string) [256]bool {
	var t [256]bool
	for i := range 256 {
		c := byte(i)
		t[i] = c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
			c == '-' || c == '_' || c == '.' || c == '~' ||
			strings.IndexByte(alsoSafe, c) >= 0
	}
	return t
}

// AppendQueryEscape appends the percent-encoded form of s to dst and returns
// the extended slice. The output is byte-identical to net/url.QueryEscape:
// unreserved bytes pass through, space becomes '+', everything else becomes
// %XX with uppercase hex. dst must not alias s: escaping can grow the input,
// so in-place operation is impossible.
func AppendQueryEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &queryNoEscapeTable, escapeQuery)
}

// AppendPathEscape appends the percent-encoded path-segment form of s to dst
// and returns the extended slice. The output is byte-identical to
// net/url.PathEscape; the aliasing rule matches AppendQueryEscape.
func AppendPathEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &pathNoEscapeTable, escapePath)
}

// appendEscape copies runs of unescaped bytes wholesale and expands the rest,
// in one pass with no intermediate allocation — unlike net/url, which counts
// in a first pass and then rebuilds the string byte by byte.
func appendEscape(dst, s []byte, noEscape *[256]bool, mode escapeMode) []byte {
	i, n := 0, len(s)
	for i < n {
		j := i
		for j < n && noEscape[s[j]] {
			j++
		}
		dst = append(dst, s[i:j]...)
		if j == n {
			break
		}
		if c := s[j]; mode == escapeQuery && c == ' ' {
			dst = append(dst, '+')
		} else {
			dst = append(dst, '%', upperHexDigits[c>>4], upperHexDigits[c&0x0F])
		}
		i = j + 1
	}
	return dst
}

// AppendQueryUnescape appends the decoded form of the query component s to
// dst and returns the extended slice, converting '+' to space and %XX to the
// byte it encodes, exactly like net/url.QueryUnescape including its
// url.EscapeError values for malformed escapes. On error the returned slice
// is dst with its original length (its backing array may still have been
// reallocated by growth). Decoding never grows the input, so dst may be
// s[:0] on a common backing array to decode in place; any other overlap is
// invalid. Note that an in-place decode that fails has already overwritten
// the prefix of s before the malformed escape with decoded bytes — treat s
// as consumed once an in-place decode starts, error or not.
func AppendQueryUnescape[S byteSeq](dst []byte, s S) ([]byte, error) {
	return appendUnescape(dst, unsafeconv.Bytes(s), escapeQuery)
}

// AppendPathUnescape appends the decoded form of the path component s to dst
// and returns the extended slice, decoding %XX and leaving '+' verbatim,
// exactly like net/url.PathUnescape. Error and aliasing behavior match
// AppendQueryUnescape.
func AppendPathUnescape[S byteSeq](dst []byte, s S) ([]byte, error) {
	return appendUnescape(dst, unsafeconv.Bytes(s), escapePath)
}

// appendUnescape jumps between escape sites with vectorized scans —
// IndexAny2 (SWAR/AVX2) when '+' also needs rewriting, bytes.IndexByte
// otherwise — and copies the clean spans between them wholesale, so input
// without escapes costs one scan and one copy.
func appendUnescape(dst, s []byte, mode escapeMode) ([]byte, error) {
	off := len(dst)
	i, n := 0, len(s)
	for i < n {
		var j int
		if mode == escapeQuery {
			j = IndexAny2(s[i:], '%', '+')
		} else {
			j = bytes.IndexByte(s[i:], '%')
		}
		if j < 0 {
			return append(dst, s[i:]...), nil
		}
		j += i
		dst = append(dst, s[i:j]...)
		if s[j] == '+' {
			dst = append(dst, ' ')
			i = j + 1
			continue
		}
		// s[j] == '%': the error substrings below match net/url's unescape.
		if j+2 >= n {
			return dst[:off], url.EscapeError(s[j:])
		}
		hi := unhexTable[s[j+1]]
		lo := unhexTable[s[j+2]]
		if hi > 0x0F || lo > 0x0F {
			return dst[:off], url.EscapeError(s[j : j+3])
		}
		dst = append(dst, hi<<4|lo)
		i = j + 3
	}
	return dst, nil
}
