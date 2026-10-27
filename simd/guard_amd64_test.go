//go:build amd64 && linux

package simd

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

// Go rounds allocations up to a size class, so an over-read of an ordinary slice lands on
// zero-filled slack and passes. These tests call the kernels directly with the input flush
// against an unmapped page, so an over-read faults.
//
// countNonASCIIAVX2 is why this file exists: an earlier version read whole vectors past the slice.

func guardedTail(t *testing.T, n int) []byte {
	t.Helper()
	pageSize := syscall.Getpagesize()
	require.LessOrEqual(t, n, pageSize, "input must fit in one page")

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
	for _, n := range []int{
		1, 2, 7, 8, 15, 16, 31, 32, 33, 63, 64, 65, 95,
		96, 127, 128, 129, 159, 160, 161, 200, 255, 256, 257, 384, 385,
	} {
		buf := guardedTail(t, n)
		for i := range buf {
			buf[i] = byte('a' + i%23)
		}
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

		// The word scan needs a word-free haystack to walk to the end.
		for i := range buf {
			buf[i] = '.'
		}
		require.Equal(t, -1, memchrWordAVX2(buf), "memchrWordAVX2 len %d", n)
	}
}

// Test_CountNonASCIIAVX2_CountsExactly pins the tail count with non-ASCII data: an
// over-read would inflate it, which ASCII-only input cannot show.
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
