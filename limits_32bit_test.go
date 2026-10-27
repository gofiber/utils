//go:build 386 || arm

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/utils/v2/internal/hugeslice"
)

// A word loop testing i+K <= n wraps near math.MaxInt32 and slices with a negative bound. This scans
// such a slice (2 GiB), so it runs only with GOFIBER_TEST_HUGE_SLICES=1.
func Test_HugeSlice_Scans(t *testing.T) {
	b := hugeslice.Get(t)

	require.True(t, IsASCII(b))
	require.Equal(t, -1, IndexAny2(b, 1, 2))
	require.Equal(t, -1, IndexAny3(b, 1, 2, 3))
	require.Equal(t, -1, IndexFold(b, "abcdefghijklmnop"))
	require.True(t, EqualFold(b, b))
	require.NotPanics(t, func() { _ = HashFold(b) })
}
