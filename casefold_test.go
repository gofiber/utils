package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	casebytes "github.com/gofiber/utils/v2/bytes"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	casestrings "github.com/gofiber/utils/v2/strings"
	"github.com/stretchr/testify/require"
)

// Test_CaseFold_Randomized cross-checks utils.EqualFold and the bytes/strings
// converters against the stdlib and a scalar ASCII fold over random inputs.
// The other folding helpers live in search_test.go and headerkey_test.go, and
// EqualFold is built on swar rather than internal/caseconv, so its assertions
// here and in byteseq_test.go are the only ones pinning it to strings.EqualFold.
func Test_CaseFold_Randomized(t *testing.T) {
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
	}
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
