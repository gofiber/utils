package swar_test

import (
	"fmt"

	"github.com/gofiber/utils/v2/swar"
)

// ExampleZeroLanes shows the canonical first-match scan. Its overlapping last word is
// safe: lanes before the loop's exit index are known non-matching, so the first set
// lane falls in the new bytes.
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

// ExampleZeroLanes_blockScan ORs the masks of several words and tests once per block, resolving
// the lower word first. Safe because a ZeroLanes mask is zero iff its word has no match and its
// lowest set lane is exact; MatchByteMask and MatchRangeMask combine the same way. Wider blocks
// add an `if` per word to the ladder.
func ExampleZeroLanes_blockScan() {
	indexByte := func(s string, c byte) int {
		needle := swar.Broadcast(c)
		i := 0
		for ; i+16 <= len(s); i += 16 {
			m0 := swar.ZeroLanes(swar.Load8(s, i) ^ needle)
			m1 := swar.ZeroLanes(swar.Load8(s, i+8) ^ needle)
			if m0|m1 != 0 {
				if m0 != 0 {
					return i + swar.FirstLane(m0)
				}
				return i + 8 + swar.FirstLane(m1)
			}
		}
		for ; i+8 <= len(s); i += 8 {
			if m := swar.ZeroLanes(swar.Load8(s, i) ^ needle); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		for ; i < len(s); i++ {
			if s[i] == c {
				return i
			}
		}
		return -1
	}

	fmt.Println(indexByte("cache-control: public, max-age=604800, immutable", '=')) //nolint:errcheck // example output
	fmt.Println(indexByte("cache-control: public", '='))                            //nolint:errcheck // example output
	// Output:
	// 30
	// -1
}

// ExampleMatchRangeMask finds the first and last ASCII digit of one word.
func ExampleMatchRangeMask() {
	w := swar.Load8("Fiber123", 0)
	digits := swar.MatchRangeMask(w, '0', '9')
	fmt.Println(swar.FirstLane(digits), swar.LastLane(digits)) //nolint:errcheck // example output
	// Output: 5 7
}

// ExampleLoad8 shows the little-endian lane order: the lowest index is lane 0.
func ExampleLoad8() {
	w := swar.Load8([]byte("ABCDEFGH"), 0)
	fmt.Printf("lane0=%c lane7=%c\n", byte(w), byte(w>>56)) //nolint:errcheck // example output
	// Output: lane0=A lane7=H
}
