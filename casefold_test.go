package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"

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

// deprecated holds the wrappers under test, called through it because their `//go:fix inline`
// makes vet flag direct calls.
var deprecated = struct {
	toLower, toUpper           func(string) string
	toLowerBytes, toUpperBytes func([]byte) []byte
}{ToLower, ToUpper, ToLowerBytes, ToUpperBytes}

func Test_Case_StdlibParity(t *testing.T) {
	t.Parallel()
	check := func(s string) {
		t.Helper()
		ascii := utf8.ValidString(s) && IsASCII(s)
		wantLower, wantUpper := asciiFoldLower(s), asciiFoldUpper(s)
		if ascii && (strings.ToLower(s) != wantLower || strings.ToUpper(s) != wantUpper) {
			t.Fatalf("the ASCII fold disagrees with the stdlib on ASCII input %q", s)
		}
		if got := casestrings.ToLower(s); got != wantLower {
			t.Fatalf("strings.ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := casestrings.ToUpper(s); got != wantUpper {
			t.Fatalf("strings.ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}
		if got := deprecated.toLower(s); got != wantLower {
			t.Fatalf("ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := deprecated.toUpper(s); got != wantUpper {
			t.Fatalf("ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}

		b := []byte(s)
		if got := casebytes.ToLower(b); string(got) != wantLower {
			t.Fatalf("bytes.ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := casebytes.ToUpper(b); string(got) != wantUpper {
			t.Fatalf("bytes.ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}
		if string(b) != s {
			t.Fatalf("bytes.ToLower/ToUpper modified their input: %q became %q", s, b)
		}

		lower := []byte(s)
		if got := casebytes.UnsafeToLower(lower); string(got) != wantLower || string(lower) != wantLower || (len(got) > 0 && &got[0] != &lower[0]) {
			t.Fatalf("bytes.UnsafeToLower(%q) = %q, input %q, want %q in place", s, got, lower, wantLower)
		}
		upper := []byte(s)
		if got := casebytes.UnsafeToUpper(upper); string(got) != wantUpper || string(upper) != wantUpper || (len(got) > 0 && &got[0] != &upper[0]) {
			t.Fatalf("bytes.UnsafeToUpper(%q) = %q, input %q, want %q in place", s, got, upper, wantUpper)
		}
		lowerB, upperB := []byte(s), []byte(s)
		if got := deprecated.toLowerBytes(lowerB); string(got) != wantLower || string(lowerB) != wantLower {
			t.Fatalf("ToLowerBytes(%q) = %q, input %q, want %q in place", s, got, lowerB, wantLower)
		}
		if got := deprecated.toUpperBytes(upperB); string(got) != wantUpper || string(upperB) != wantUpper {
			t.Fatalf("ToUpperBytes(%q) = %q, input %q, want %q in place", s, got, upperB, wantUpper)
		}
		lowerS, upperS := []byte(s), []byte(s)
		if got := casestrings.UnsafeToLower(unsafeconv.UnsafeString(lowerS)); got != wantLower || string(lowerS) != wantLower {
			t.Fatalf("strings.UnsafeToLower(%q) = %q, backing %q, want %q", s, got, lowerS, wantLower)
		}
		if got := casestrings.UnsafeToUpper(unsafeconv.UnsafeString(upperS)); got != wantUpper || string(upperS) != wantUpper {
			t.Fatalf("strings.UnsafeToUpper(%q) = %q, backing %q, want %q", s, got, upperS, wantUpper)
		}
	}

	for c := range 256 {
		cs := string([]byte{byte(c)})
		check(cs)
		check("aB" + cs + "Cd")
		check(strings.Repeat("xY", 9) + cs + strings.Repeat("Zz", 9))
	}
	forEachString("aAzZ@[`{\xe9\xc9", 4, check)
	// Lengths around the word and 32-byte block sizes, the one letter to convert at every position.
	for n := range 100 {
		for pos := range n {
			plain := []byte(strings.Repeat("-4/", 40)[:n])
			up, low := bytes.Clone(plain), bytes.Clone(plain)
			up[pos], low[pos] = 'Q', 'q'
			check(string(up))
			check(string(low))
		}
		check(strings.Repeat("-4/", 40)[:n])
	}
}

// Test_Case_ASCIIOnly pins that case conversion stays in ASCII, where the stdlib maps every Unicode letter.
func Test_Case_ASCIIOnly(t *testing.T) {
	t.Parallel()
	for _, in := range []string{"\u00e9", "\u00c9", "\u0131", "\u212a", "\u00df", "\u01c5", "\u00e9 \u00c9 123 !", "\xff\xfe"} {
		if casestrings.ToLower(in) != in || casestrings.ToUpper(in) != in {
			t.Fatalf("case conversion changed the non-ASCII %q", in)
		}
	}
	if strings.ToUpper("\u00e9") == "\u00e9" {
		t.Fatal("the stdlib no longer upper-cases U+00E9: the documented difference is gone")
	}
	if got := casestrings.ToUpper("caf\u00e9 ok"); got != "CAF\u00e9 OK" {
		t.Fatalf("ToUpper left %q", got)
	}
}

// Test_CaseBytes_Aliasing pins that nothing to convert returns the argument itself, where bytes.ToLower copies.
func Test_CaseBytes_Aliasing(t *testing.T) {
	t.Parallel()
	for _, n := range []int{0, 1, 7, 8, 9, 31, 32, 33, 100} {
		plain := bytes.Repeat([]byte{'a'}, n)
		got := casebytes.ToLower(plain)
		if n > 0 && &got[0] != &plain[0] {
			t.Fatalf("ToLower(%d lower-case bytes) copied", n)
		}
		if std := bytes.ToLower(plain); n > 0 && &std[0] == &plain[0] {
			t.Fatalf("bytes.ToLower(%d lower-case bytes) no longer copies: the documented difference is gone", n)
		}

		upper := bytes.Repeat([]byte{'A'}, n)
		got = casebytes.ToUpper(upper)
		if n > 0 && &got[0] != &upper[0] {
			t.Fatalf("ToUpper(%d upper-case bytes) copied", n)
		}

		mixed := bytes.Repeat([]byte{'a'}, n)
		if n > 0 {
			mixed[n-1] = 'Q'
			if res := casebytes.ToLower(mixed); &res[0] == &mixed[0] || mixed[n-1] != 'Q' {
				t.Fatalf("ToLower(%d bytes ending in Q) did not return a fresh array", n)
			}
		}
	}
}
