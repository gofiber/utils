package swar_test

import (
	"fmt"

	"github.com/gofiber/utils/v2/swar"
)

// ExampleZeroLanes shows the canonical first-match scan: broadcast the
// needle once, scan full words, and finish 8+ byte inputs with one
// overlapping word at len(s)-8. The overlap is safe because lanes before
// the loop's exit index are known non-matching, so the first set lane
// always falls into the new bytes.
func ExampleZeroLanes() {
	indexByte := func(s string, c byte) int {
		n := len(s)
		if n < swar.WordLen {
			for i := range n {
				if s[i] == c {
					return i
				}
			}
			return -1
		}
		needle := swar.Broadcast(c)
		i := 0
		for ; i+swar.WordLen <= n; i += swar.WordLen {
			if m := swar.ZeroLanes(swar.Load8(s, i) ^ needle); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		if i == n {
			return -1
		}
		if m := swar.ZeroLanes(swar.Load8(s, n-swar.WordLen) ^ needle); m != 0 {
			return n - swar.WordLen + swar.FirstLane(m)
		}
		return -1
	}

	fmt.Println(indexByte("cache-control: no-store", ':')) //nolint:errcheck // example output
	fmt.Println(indexByte("etag", ':'))                    //nolint:errcheck // example output
	// Output:
	// 13
	// -1
}

// ExampleMatchRangeMask classifies one word of input: which lanes hold an
// ASCII digit, and where the first and last digit sit.
func ExampleMatchRangeMask() {
	w := swar.Load8("Fiber123", 0)
	digits := swar.MatchRangeMask(w, '0', '9')
	fmt.Println(swar.FirstLane(digits), swar.LastLane(digits)) //nolint:errcheck // example output
	// Output: 5 7
}

// ExampleLoad8 shows the little-endian lane order: the byte at the lowest
// index becomes lane 0, the least significant byte of the word.
func ExampleLoad8() {
	w := swar.Load8([]byte("ABCDEFGH"), 0)
	fmt.Printf("lane0=%c lane7=%c\n", byte(w), byte(w>>56)) //nolint:errcheck // example output
	// Output: lane0=A lane7=H
}
