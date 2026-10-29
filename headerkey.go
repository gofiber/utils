package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// headerTokenTable flags the bytes valid in an HTTP header field name — the
// RFC 9110 token alphabet, matching net/http's validHeaderFieldByte.
var headerTokenTable = buildHeaderTokenTable()

func buildHeaderTokenTable() [256]bool {
	var t [256]bool
	for c := byte('0'); c <= '9'; c++ {
		t[c] = true
	}
	for c := byte('a'); c <= 'z'; c++ {
		t[c] = true
	}
	for c := byte('A'); c <= 'Z'; c++ {
		t[c] = true
	}
	for _, c := range []byte("!#$%&'*+-.^_`|~") {
		t[c] = true
	}
	return t
}

// CanonicalHeaderKey returns the canonical form of the HTTP header key:
// the first letter and any letter following a hyphen upper-cased, the rest
// lower-cased. The output matches net/http.CanonicalHeaderKey byte for
// byte, including its guard: a key containing any byte outside the header
// token alphabet is returned unchanged. Unlike the stdlib it accepts byte
// slices as well as strings and allocates only when the key actually needs
// rewriting — already-canonical keys, the overwhelmingly common case on
// receive paths, are validated in one table pass and returned as-is
// without the stdlib's canonical-key cache lookup.
func CanonicalHeaderKey[S byteSeq](key S) S {
	upper := true
	fix := -1
	for i := 0; i < len(key); i++ {
		c := key[i]
		if !headerTokenTable[c] {
			return key
		}
		if fix < 0 && (upper && 'a' <= c && c <= 'z' || !upper && 'A' <= c && c <= 'Z') {
			fix = i
		}
		upper = c == '-'
	}
	if fix < 0 {
		return key
	}

	buf := make([]byte, len(key))
	copy(buf, unsafeconv.Bytes(key))
	upper = fix == 0 || buf[fix-1] == '-'
	for i := fix; i < len(buf); i++ {
		c := buf[i]
		if upper {
			c = caseconv.ToUpperTable[c]
		} else {
			c = caseconv.ToLowerTable[c]
		}
		buf[i] = c
		upper = c == '-'
	}
	// The zero-copy view is sound: buf is uniquely owned here and is never
	// written again, which is what the string instantiation requires.
	return unsafeconv.Seq[S](buf)
}
