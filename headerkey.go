package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

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

// CanonicalHeaderKey returns key in canonical HTTP header form, as
// net/http.CanonicalHeaderKey does: the first letter and any letter after a
// hyphen upper-cased, the rest lower-cased. A key with a byte outside the
// token alphabet is returned unchanged. It allocates only if the key changes.
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
	// buf is never written again, so the zero-copy string view is sound.
	return unsafeconv.Seq[S](buf)
}
