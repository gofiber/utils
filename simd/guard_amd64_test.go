//go:build amd64 && linux

package simd

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

// Every kernel in this package is self-bounding: it clamps its own vector
// loops and finishes short remainders itself, so it may be called directly
// with any length. That is not testable with ordinary slices — Go rounds
// allocations up to a size class, so a read past the end lands on
// zero-filled slack and the wrong answer never materializes. These tests
// place the last byte of the input flush against an unmapped page, turning
// any over-read into a fault instead of a silent success.
//
// countNonASCIIAVX2 is the reason this file exists: it is the one kernel
// whose result depends on the trailing bytes, so it cannot use the
// overlapping-window epilogue its siblings use, and an earlier version
// looped while base < end and read whole vectors past the slice.

// guardedTail returns a slice of n bytes whose final byte is the last
// readable byte before an unmapped page, plus a cleanup function. Reading
// even one byte past the end faults.
func guardedTail(t *testing.T, n int) []byte {
	t.Helper()
	pageSize := syscall.Getpagesize()
	require.LessOrEqual(t, n, pageSize, "input must fit in one page")

	// Two pages, then drop read access to the second one.
	region, err := syscall.Mmap(-1, 0, 2*pageSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANON)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, syscall.Munmap(region))
	})
	require.NoError(t, syscall.Mprotect(region[pageSize:], syscall.PROT_NONE))

	return region[pageSize-n : pageSize : pageSize]
}

func Test_KernelsAreSelfBounding(t *testing.T) {
	if !Accelerated() {
		t.Skip("AVX2 kernels not active on this CPU")
	}
	// Lengths that leave every possible sub-vector and sub-block
	// remainder, including the 1..31 that has no whole vector at all.
	for _, n := range []int{
		1, 2, 7, 8, 15, 16, 31, 32, 33, 63, 64, 65, 95,
		96, 127, 128, 129, 159, 160, 161, 200, 255, 256, 257, 384, 385,
	} {
		buf := guardedTail(t, n)
		for i := range buf {
			buf[i] = byte('a' + i%23)
		}
		// Each kernel is called directly, bypassing the *Impl dispatch, so
		// the guard is on the kernel itself and not on a wrapper.
		require.True(t, isASCIIAVX2(buf), "isASCIIAVX2 len %d", n)
		require.Equal(t, -1, firstNonASCIIAVX2(buf), "firstNonASCIIAVX2 len %d", n)
		if hasPOPCNT {
			require.Equal(t, 0, countNonASCIIAVX2(buf), "countNonASCIIAVX2 len %d", n)
		}
		require.Equal(t, -1, memchr2AVX2(buf, 'z', 'y'), "memchr2AVX2 len %d", n)
		require.Equal(t, -1, memchr3AVX2(buf, 'z', 'y', 'x'), "memchr3AVX2 len %d", n)
		require.Equal(t, -1, memchrDigitAVX2(buf), "memchrDigitAVX2 len %d", n)
		require.Equal(t, -1, memchrNotWordAVX2(buf), "memchrNotWordAVX2 len %d", n)
		for _, offset := range []int{1, 2, 8, 31, 32, 33} {
			if offset >= n {
				continue
			}
			require.Equal(t, -1, memchrPairAVX2(buf, 'z', 'y', offset),
				"memchrPairAVX2 len %d offset %d", n, offset)
		}

		// The word scan needs a haystack with no word characters, so it
		// too has to walk to the end.
		for i := range buf {
			buf[i] = '.'
		}
		require.Equal(t, -1, memchrWordAVX2(buf), "memchrWordAVX2 len %d", n)
	}
}

// Test_CountNonASCIIAVX2_CountsExactly pins the counting kernel's tail
// handling: an over-read that happened to land on non-ASCII bytes would
// inflate the count rather than fault, which the guard test alone cannot
// distinguish from a correct scan of zero-filled slack.
func Test_CountNonASCIIAVX2_CountsExactly(t *testing.T) {
	if !Accelerated() || !hasPOPCNT {
		t.Skip("AVX2+POPCNT kernels not active on this CPU")
	}
	for n := range 400 {
		buf := guardedTail(t, max(n, 1))[:n]
		want := 0
		for i := range buf {
			buf[i] = byte(i * 7)
			if buf[i] >= 0x80 {
				want++
			}
		}
		require.Equal(t, want, countNonASCIIAVX2(buf), "len %d", n)
		require.Equal(t, want, CountNonASCII(buf), "dispatched, len %d", n)
	}
}
