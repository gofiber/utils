//go:build 386 || arm

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/utils/v2/internal/hugeslice"
)

// The word loops of this package test whether a word is left before loading
// it. Written as i+K <= n, the sum wraps on a 32-bit platform for a length
// within K of math.MaxInt32, the longest slice there is, and the loop then
// slices with a negative bound. These tests scan such a slice end to end. It
// takes 2 GiB of address space and several seconds, so they run on request
// only: GOFIBER_TEST_HUGE_SLICES=1 go test.
func Test_HugeSlice_Scans(t *testing.T) {
	b := hugeslice.Get(t)

	require.True(t, IsASCII(b))
	require.Equal(t, -1, IndexAny2(b, 1, 2))
	require.Equal(t, -1, IndexAny3(b, 1, 2, 3))
	require.Equal(t, -1, IndexFold(b, "abcdefghijklmnop"))
	require.True(t, EqualFold(b, b))
	require.NotPanics(t, func() { _ = HashFold(b) })
}
