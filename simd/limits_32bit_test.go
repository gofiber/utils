//go:build 386 || arm

package simd

import (
	"math"
	"testing"

	"github.com/gofiber/utils/v2/internal/hugeslice"
)

// Portable scans loop on i+K <= n, which wraps on 32-bit within K of math.MaxInt32
// and then slices with a negative bound; MemchrPair/far wraps through offset+WordLen
// too. Opt in with GOFIBER_TEST_HUGE_SLICES=1 (2 GiB).
func Test_HugeSlice_Scans(t *testing.T) {
	b := hugeslice.Get(t)
	for name, got := range map[string]int{
		"Memchr2":        Memchr2(b, 1, 2),
		"Memchr3":        Memchr3(b, 1, 2, 3),
		"MemchrPair":     MemchrPair(b, 1, 2, 1),
		"MemchrPair/far": MemchrPair(b, 1, 2, math.MaxInt32-1),
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
