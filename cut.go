package utils

import (
	"bytes"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// CutByte slices s around the first occurrence of sep, returning the text
// before sep, the text after it, and whether sep was found — strings.Cut
// and bytes.Cut for a single-byte separator, generic over strings and byte
// slices. When sep is absent the results are s, the zero value of S ("" or
// nil), and false, matching the stdlib pair. The scan is a single
// IndexByte, so byte-slice callers skip both the []byte{sep} needle
// bytes.Cut wants and the length dispatch inside bytes.Index.
func CutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.IndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}

// LastCutByte is CutByte around the last occurrence of sep: the text before
// the last sep, the text after it, and whether sep was found. It is the
// natural spelling of host:port and name.ext splits, which strings.Cut
// cannot express.
func LastCutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.LastIndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}
