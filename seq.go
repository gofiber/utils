package utils

import (
	"bytes"
	"iter"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// SplitTrimSeq iterates over the sep-separated elements of s with ASCII
// whitespace trimmed and empty elements skipped, as HTTP list fields are read
// (RFC 9110 Section 5.6.1); the elements alias s.
func SplitTrimSeq[S byteSeq](s S, sep byte) iter.Seq[S] {
	return func(yield func(S) bool) {
		for {
			i := bytes.IndexByte(unsafeconv.Bytes(s), sep)
			elem := s
			if i >= 0 {
				elem = s[:i]
			}
			if elem = TrimSpace(elem); len(elem) > 0 && !yield(elem) {
				return
			}
			if i < 0 {
				return
			}
			s = s[i+1:]
		}
	}
}
