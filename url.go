package utils

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
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

// The no-escape tables mirror net/url's shouldEscape and are pinned to it by per-byte tests.
var (
	queryNoEscapeTable = buildNoEscapeTable("")
	pathNoEscapeTable  = buildNoEscapeTable("$&+:=@")
	// "$&+:=@" are the sub-delims encodePathSegment keeps; the multi-segment table adds '/'.
	pathSegmentsNoEscapeTable = buildNoEscapeTable("$&+:=@/")
)

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

// AppendQueryEscape appends s to dst, percent-encoded exactly like net/url.QueryEscape
// (space becomes '+'). dst must not alias s.
func AppendQueryEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &queryNoEscapeTable, escapeQuery)
}

// AppendPathEscape appends s to dst escaped as one path segment, exactly like net/url.PathEscape.
// dst must not alias s.
func AppendPathEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &pathNoEscapeTable, escapePath)
}

// AppendPathSegmentsEscape appends the multi-segment path s to dst, escaping each segment
// like AppendPathEscape but keeping the '/' separators. Empty and dot segments pass through
// unchanged, so clean or reject untrusted input before appending it under a fixed prefix.
// dst must not alias s.
func AppendPathSegmentsEscape[S byteSeq](dst []byte, s S) []byte {
	return appendEscape(dst, unsafeconv.Bytes(s), &pathSegmentsNoEscapeTable, escapePath)
}

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

// AppendQueryUnescape appends the decoded query component s to dst, like net/url.QueryUnescape
// ('+' becomes space), including its url.EscapeError for malformed escapes. On error it returns
// dst at its original length. dst may be s[:0] to decode in place; any other overlap is invalid,
// and a failed in-place decode has already overwritten a prefix of s, so treat s as consumed.
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
	return dst, nil
}
