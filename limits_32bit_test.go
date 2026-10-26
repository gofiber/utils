//go:build 386 || arm

package utils

import (
	"math"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// The word loops of this package test whether a word is left before loading
// it. Written as i+K <= n, the sum wraps on a 32-bit platform for a length
// within K of math.MaxInt32, the longest slice there is, and the loop then
// slices with a negative bound. These tests scan such a slice end to end. It
// takes 2 GiB of address space and several seconds, so they run on request
// only: GOFIBER_TEST_HUGE_SLICES=1 go test.
var (
	hugeOnce sync.Once
	hugeBuf  []byte
)

// hugeSlice returns the one slice of math.MaxInt32 zero bytes of the test
// binary, allocated on first use. It is never written, so its pages stay the
// shared zero page, and a second slice would not fit the address space.
func hugeSlice(t *testing.T) []byte {
	t.Helper()
	if os.Getenv("GOFIBER_TEST_HUGE_SLICES") == "" {
		t.Skip("set GOFIBER_TEST_HUGE_SLICES=1 to scan a slice of math.MaxInt32 bytes")
	}
	hugeOnce.Do(func() { hugeBuf = make([]byte, math.MaxInt32) })
	return hugeBuf
}

func Test_HugeSlice_Scans(t *testing.T) {
	b := hugeSlice(t)

	require.True(t, IsASCII(b))
	require.Equal(t, -1, IndexAny2(b, 1, 2))
	require.Equal(t, -1, IndexAny3(b, 1, 2, 3))
	require.Equal(t, -1, IndexFold(b, "abcdefghijklmnop"))
	require.True(t, EqualFold(b, b))
	require.NotPanics(t, func() { _ = HashFold(b) })
}
