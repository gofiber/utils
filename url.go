package utils

import (
	"bytes"
	"net/url"
	"slices"
	"strings"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

type escapeMode uint8

const (
	escapeQuery escapeMode = iota
	escapePath
)

const upperHexDigits = "0123456789ABCDEF"

// unhexTable maps ASCII bytes to hex values; 0xFF marks non-hex bytes.
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

// The escape tables mirror net/url's shouldEscape and are pinned to it by per-byte tests: 1 marks
// a byte that is not copied as it is.
var (
	queryEscapeTable = buildEscapeTable("")
	pathEscapeTable  = buildEscapeTable("$&+:=@")
	// "$&+:=@" are the sub-delims encodePathSegment keeps; the multi-segment table adds '/'.
	pathSegmentsEscapeTable = buildEscapeTable("$&+:=@/")
)

func buildEscapeTable(alsoSafe string) [256]uint8 {
	var t [256]uint8
	for i := range 256 {
		c := byte(i)
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
			c == '-' || c == '_' || c == '.' || c == '~' ||
			strings.IndexByte(alsoSafe, c) >= 0 {
			continue
		}
		t[i] = 1
	}
	return t
}

// AppendQueryEscape appends s to dst, percent-encoded exactly like net/url.QueryEscape
// (space becomes '+'). dst must not alias s; if it is too small it grows once, to the size of
// the result.
func AppendQueryEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &queryEscapeTable, escapeQuery)
}

// AppendPathEscape appends s to dst escaped as one path segment, exactly like net/url.PathEscape.
// dst must not alias s.
func AppendPathEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &pathEscapeTable, escapePath)
}

// AppendPathSegmentsEscape appends the multi-segment path s to dst, escaping each segment
// like AppendPathEscape but keeping the '/' separators. Empty and dot segments pass through
// unchanged, so clean or reject untrusted input before appending it under a fixed prefix.
// dst must not alias s.
func AppendPathSegmentsEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &pathSegmentsEscapeTable, escapePath)
}

// escapesInWord is nonzero if any of the eight bytes of w is to be escaped: eight lookups and
// one branch from one load, where a loop over the bytes loads and branches for each.
func escapesInWord(w uint64, escape *[256]uint8) uint8 {
	return escape[byte(w)] | escape[byte(w>>8)] | escape[byte(w>>16)] | escape[byte(w>>24)] |
		escape[byte(w>>32)] | escape[byte(w>>40)] | escape[byte(w>>48)] | escape[w>>56]
}

// indexEscape returns the index of the first byte of s, from i on, that is to be escaped, or len(s).
func indexEscape(s []byte, i int, escape *[256]uint8) int {
	for ; i <= len(s)-8; i += 8 {
		if escapesInWord(swar.Load8(s, i), escape) != 0 {
			break
		}
	}
	for i < len(s) && escape[s[i]] == 0 {
		i++
	}
	return i
}

// escapedLen returns the length of s escaped: two more bytes for each one written as %XX.
func escapedLen(s []byte, escape *[256]uint8, mode escapeMode) int {
	hex := 0
	for _, c := range s {
		hex += int(escape[c])
	}
	if mode == escapeQuery {
		hex -= bytes.Count(s, []byte{' '}) // a space becomes a '+'
	}
	return len(s) + 2*hex
}

// shortRun is how far a run of unescaped bytes is copied a byte at a time; a longer one is
// worth a word scan and one copy.
const shortRun = 8

func appendEscape(dst, s []byte, escape *[256]uint8, mode escapeMode) []byte {
	n := len(s)
	i := indexEscape(s, 0, escape)
	if i == n {
		return append(dst, s...) // nothing to escape, the common case: one scan, one copy
	}
	if cap(dst)-len(dst) < i+3*(n-i) {
		// dst may be too small: grow it once, to the exact size, where appending would grow it
		// again and again.
		dst = slices.Grow(dst, i+escapedLen(s[i:], escape, mode))
	}
	dst = append(dst, s[:i]...)
	for {
		// s[i] is to be escaped.
		if c := s[i]; mode == escapeQuery && c == ' ' {
			dst = append(dst, '+')
		} else {
			dst = append(dst, '%', upperHexDigits[c>>4], upperHexDigits[c&0x0F])
		}
		i++
		short := min(i+shortRun, n)
		for i < short && escape[s[i]] == 0 {
			dst = append(dst, s[i])
			i++
		}
		if i == short && i < n {
			j := indexEscape(s, i, escape)
			dst = append(dst, s[i:j]...)
			i = j
		}
		if i == n {
			return dst
		}
	}
}

// AppendQueryUnescape appends the decoded query component s to dst, like net/url.QueryUnescape
// ('+' becomes space), including its url.EscapeError for malformed escapes. On error it returns
// dst at its original length. dst may be s[:0] to decode in place; any other overlap is invalid,
// and a failed in-place decode has already overwritten a prefix of s, so treat s as consumed.
// A dst too small grows once, to the size of the result.
func AppendQueryUnescape[S byteSeq](dst []byte, s S) ([]byte, error) {
	return appendUnescape(dst, unsafeconv.Bytes(s), escapeQuery)
}

// AppendPathUnescape appends the decoded path component s to dst, like net/url.PathUnescape
// ('+' is kept). Error and aliasing behavior match AppendQueryUnescape.
func AppendPathUnescape[S byteSeq](dst []byte, s S) ([]byte, error) {
	return appendUnescape(dst, unsafeconv.Bytes(s), escapePath)
}

func appendUnescape(dst, s []byte, mode escapeMode) ([]byte, error) {
	off := len(dst)
	i, n := 0, len(s)
	var j int
	if mode == escapeQuery {
		j = IndexAny2(s, '%', '+')
	} else {
		j = bytes.IndexByte(s, '%')
	}
	if j < 0 {
		return append(dst, s...), nil // nothing to decode, the common case: one scan, one copy
	}
	if cap(dst)-len(dst) < n {
		dst = growUnescaped(dst, s)
	}
	if j < n-3 && s[j] == '%' && s[j+3] == '%' {
		// One escape right after another, as in percent-encoded UTF-8: a scan for each would
		// cost more than the bytes between them.
		return appendUnescapeDense(dst, s, j, mode)
	}
	for {
		j += i
		dst = append(dst, s[i:j]...)
		if s[j] == '+' {
			dst = append(dst, ' ')
			i = j + 1
		} else {
			// s[j] == '%': the error substrings match net/url's unescape.
			if j >= n-2 {
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
		if i == n {
			return dst, nil
		}
		if mode == escapeQuery {
			j = IndexAny2(s[i:], '%', '+')
		} else {
			j = bytes.IndexByte(s[i:], '%')
		}
		if j < 0 {
			return append(dst, s[i:]...), nil
		}
	}
}

// growUnescaped makes room in dst for s decoded: once and to the exact size, where appending
// would grow it again and again. Each %XX becomes one byte, and s[:0] as dst has the room.
func growUnescaped(dst, s []byte) []byte {
	return slices.Grow(dst, max(len(s)-2*bytes.Count(s, []byte{'%'}), 0))
}

// denseRun is how many bytes appendUnescapeDense copies without meeting an escape before it
// looks for the next one with the vector scan.
const denseRun = 8

// appendUnescapeDense is appendUnescape from the first escape of s on, at i, a byte at a time:
// where escapes come in clusters, a scan and a copy for each cost more than the loop.
func appendUnescapeDense(dst, s []byte, i int, mode escapeMode) ([]byte, error) {
	off := len(dst)
	dst = append(dst, s[:i]...)
	for run := 0; i < len(s); {
		c := s[i]
		switch {
		case c == '%':
			if i >= len(s)-2 {
				return dst[:off], url.EscapeError(s[i:])
			}
			hi := unhexTable[s[i+1]]
			lo := unhexTable[s[i+2]]
			if hi > 0x0F || lo > 0x0F {
				return dst[:off], url.EscapeError(s[i : i+3])
			}
			dst = append(dst, hi<<4|lo)
			i += 3
			run = 0
		case c == '+' && mode == escapeQuery:
			dst = append(dst, ' ')
			i++
			run = 0
		default:
			dst = append(dst, c)
			i++
			if run++; run == denseRun {
				// The cluster is behind: on to the next escape, and one copy up to it.
				var j int
				if mode == escapeQuery {
					j = IndexAny2(s[i:], '%', '+')
				} else {
					j = bytes.IndexByte(s[i:], '%')
				}
				if j < 0 {
					return append(dst, s[i:]...), nil
				}
				dst = append(dst, s[i:i+j]...)
				i += j
				run = 0
			}
		}
	}
	return dst, nil
}
