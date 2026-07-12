package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	casebytes "github.com/gofiber/utils/v2/bytes"
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	casestrings "github.com/gofiber/utils/v2/strings"
	"github.com/stretchr/testify/require"
)

func Test_SWAR_Randomized(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic test data

	for range 20000 {
		// Lengths 0..100 cover sub-word inputs, all %8 tails, and several
		// full word iterations.
		n := rng.Intn(101)
		a := make([]byte, n)
		b := make([]byte, n)
		for i := range n {
			a[i] = byte(rng.Intn(256))
			// Bias b towards being a case-variant of a so EqualFold hits
			// both outcomes.
			if rng.Intn(2) == 0 {
				b[i] = a[i] ^ 0x20
			} else {
				b[i] = byte(rng.Intn(256))
			}
		}
		as, bs := string(a), string(b)

		require.Equal(t, strings.ToLower(asciiOnly(as)), casestrings.ToLower(asciiOnly(as)))
		require.Equal(t, strings.ToUpper(asciiOnly(as)), casestrings.ToUpper(asciiOnly(as)))
		require.Equal(t, strings.EqualFold(asciiOnly(as), asciiOnly(bs)), EqualFold(asciiOnly(as), asciiOnly(bs)))

		// Non-ASCII bytes must pass through unchanged (table semantics).
		require.Equal(t, asciiFoldLower(as), casestrings.ToLower(as))
		require.Equal(t, asciiFoldUpper(as), casestrings.ToUpper(as))
		require.Equal(t, bytes.Equal(asciiFoldUpperB(a), asciiFoldUpperB(b)), EqualFold(a, b))

		require.Equal(t, asciiFoldLower(as), string(casebytes.ToLower(append([]byte(nil), a...))))
		require.Equal(t, asciiFoldUpper(as), string(casebytes.ToUpper(append([]byte(nil), a...))))
		require.Equal(t, asciiFoldLower(as), string(casebytes.UnsafeToLower(append([]byte(nil), a...))))
		require.Equal(t, asciiFoldUpper(as), string(casebytes.UnsafeToUpper(append([]byte(nil), a...))))

		// The strings in-place wrappers mutate their backing bytes, so hand
		// each a string view over a private mutable copy — string(a) may
		// reference read-only interned memory (e.g. 1-byte strings).
		la := append([]byte(nil), a...)
		require.Equal(t, asciiFoldLower(as), casestrings.UnsafeToLower(unsafeconv.UnsafeString(la)))
		ua := append([]byte(nil), a...)
		require.Equal(t, asciiFoldUpper(as), casestrings.UnsafeToUpper(unsafeconv.UnsafeString(ua)))

		// Exercise the caseconv helpers directly so their sub-word tails are
		// covered even though the public wrappers route short inputs away.
		lowered := append([]byte(nil), a...)
		caseconv.ToLowerInPlace(lowered)
		require.Equal(t, asciiFoldLower(as), string(lowered))
		uppered := append([]byte(nil), a...)
		caseconv.ToUpperInPlace(uppered)
		require.Equal(t, asciiFoldUpper(as), string(uppered))

		require.Equal(t, stdFirstIndex(a, 'A', 'Z'), caseconv.FirstUpperIndex(a))
		require.Equal(t, stdFirstIndex(a, 'a', 'z'), caseconv.FirstLowerIndex(a))

		loweredCopy := make([]byte, n)
		caseconv.ToLowerCopy(loweredCopy, a, 0)
		require.Equal(t, asciiFoldLower(as), string(loweredCopy))
		upperedCopy := make([]byte, n)
		caseconv.ToUpperCopy(upperedCopy, a, 0)
		require.Equal(t, asciiFoldUpper(as), string(upperedCopy))
	}
}

func stdFirstIndex(b []byte, lo, hi byte) int {
	for i, c := range b {
		if c >= lo && c <= hi {
			return i
		}
	}
	return -1
}

func asciiOnly(s string) string {
	b := []byte(s)
	for i := range b {
		b[i] &= 0x7f
	}
	return string(b)
}

func asciiFoldLower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func asciiFoldUpper(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 'a' - 'A'
		}
	}
	return string(b)
}

func asciiFoldUpperB(b []byte) []byte {
	return []byte(asciiFoldUpper(string(b)))
}

// Test_CaseCopy_FromOffsets pins the from > 0 contract of ToLowerCopy and
// ToUpperCopy directly: from is word-aligned and at or below the first byte
// the conversion changes, so the overlapping tail store may rewrite
// dst[n-8:from) only as a no-op. The public wrappers exercise this path
// indirectly; this test does it without them.
func Test_CaseCopy_FromOffsets(t *testing.T) {
	t.Parallel()
	for n := 8; n <= 26; n++ {
		for from := 0; from < n; from += 8 {
			// Case-stable filler (digits): bytes before firstChange must not
			// change under either conversion, or from would be invalid.
			base := make([]byte, n)
			for i := range base {
				base[i] = byte('0' + i%10)
			}
			for firstChange := from; firstChange < n; firstChange++ {
				lsrc := append([]byte(nil), base...)
				lsrc[firstChange] = 'A' + byte(firstChange%26)
				if firstChange+1 < n {
					lsrc[firstChange+1] = 0xE9 // non-ASCII passthrough
				}
				dst := make([]byte, n)
				copy(dst[:from], lsrc[:from])
				caseconv.ToLowerCopy(dst, lsrc, from)
				require.Equal(t, asciiFoldLower(string(lsrc)), string(dst),
					"lower n=%d from=%d change=%d", n, from, firstChange)

				usrc := append([]byte(nil), base...)
				usrc[firstChange] = 'a' + byte(firstChange%26)
				if firstChange+1 < n {
					usrc[firstChange+1] = 0xE9
				}
				dst2 := make([]byte, n)
				copy(dst2[:from], usrc[:from])
				caseconv.ToUpperCopy(dst2, usrc, from)
				require.Equal(t, asciiFoldUpper(string(usrc)), string(dst2),
					"upper n=%d from=%d change=%d", n, from, firstChange)
			}
		}
	}
}
