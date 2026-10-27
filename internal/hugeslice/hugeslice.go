//go:build 386 || arm

// Package hugeslice hands the 32-bit tests of this module the longest slice there is.
package hugeslice

import (
	"math"
	"os"
	"sync"
	"testing"
)

var (
	once sync.Once
	buf  []byte
)

// Get returns the one slice of math.MaxInt32 zero bytes of the test binary, allocated on first
// use, and skips the test unless GOFIBER_TEST_HUGE_SLICES is set: the slice takes 2 GiB of address
// space and the scans several seconds. It is never written, so its pages stay the shared zero
// page, and a second slice would not fit the address space.
func Get(t testing.TB) []byte {
	t.Helper()
	if os.Getenv("GOFIBER_TEST_HUGE_SLICES") == "" {
		t.Skip("set GOFIBER_TEST_HUGE_SLICES=1 to scan a slice of math.MaxInt32 bytes")
	}
	once.Do(func() { buf = make([]byte, math.MaxInt32) })
	return buf
}
