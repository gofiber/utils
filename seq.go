package utils

import (
	"bytes"
	"iter"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// SplitTrimSeq returns an iterator over the elements of s separated by
// sep, with ASCII whitespace trimmed from each element and empty elements
// skipped — the way HTTP list fields are read (RFC 9110 Section 5.6.1:
// optional whitespace around the commas, empty elements ignored). It is
// strings.SplitSeq followed by TrimSpace and an emptiness check, with
// the separator located by a single IndexByte per element, generic over
// strings and byte slices, and yielding subslices of s without copying.
// Because empty elements are dropped it is not the right tool where
// they are significant, such as validating dotted host labels.
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
