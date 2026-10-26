package utils

import (
	"bytes"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// CutByte is strings.Cut/bytes.Cut for a single-byte sep; the parts alias s.
func CutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.IndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}

// LastCutByte is CutByte around the last sep.
func LastCutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.LastIndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}
