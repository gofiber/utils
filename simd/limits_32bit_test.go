//go:build 386 || arm

package simd

import (
	"math"
	"os"
	"sync"
	"testing"
)

// The portable scans test whether a word is left before loading it. Written as
// i+K <= n, the sum wraps on a 32-bit platform for a length within K of
// math.MaxInt32, the longest slice there is, and the loop then slices with a
// negative bound. These tests scan such a slice end to end. It takes 2 GiB of
// address space and several seconds, so they run on request only:
// GOFIBER_TEST_HUGE_SLICES=1 go test.
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
	for name, got := range map[string]int{
		"Memchr2":        Memchr2(b, 1, 2),
		"Memchr3":        Memchr3(b, 1, 2, 3),
		"MemchrPair":     MemchrPair(b, 1, 2, 1),
		"MemchrPair/far": MemchrPair(b, 1, 2, math.MaxInt32-1), // offset+WordLen wraps too
		"MemchrDigit":    MemchrDigit(b),
		"MemchrDigitAt":  MemchrDigitAt(b, 1),
		"MemchrWord":     MemchrWord(b),
		"FirstNonASCII":  FirstNonASCII(b),
		"CountNonASCII":  CountNonASCII(b),
	} {
		want := -1
		if name == "CountNonASCII" {
			want = 0
		}
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}
	if !IsASCII(b) {
		t.Error("IsASCII = false, want true")
	}
}
