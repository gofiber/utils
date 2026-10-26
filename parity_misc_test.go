package utils

import (
	"encoding/base64"
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var urlSafeToken = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// Test_GenerateSecureToken_EncodingParity holds the tokens to the encoding they
// are documented to use: base64.RawURLEncoding of the random bytes. Every length
// is checked, so that both the 32-byte path and the general one, and all three
// remainders of the encoding, are decoded back to the bytes that were asked for.
func Test_GenerateSecureToken_EncodingParity(t *testing.T) {
	t.Parallel()
	strict := base64.RawURLEncoding.Strict() // refuses padding bits that are not zero
	check := func(token string, n int) {
		t.Helper()
		require.Len(t, token, base64.RawURLEncoding.EncodedLen(n), "length %d", n)
		require.Regexp(t, urlSafeToken, token, "length %d", n)
		decoded, err := strict.DecodeString(token)
		require.NoError(t, err, "length %d", n)
		require.Len(t, decoded, n, "length %d", n)
		if n >= 16 {
			require.NotEqual(t, make([]byte, n), decoded, "length %d is all zero bytes", n)
		}
	}
	for n := 1; n <= 130; n++ {
		check(GenerateSecureToken(n), n)
	}
	check(GenerateSecureToken(1<<16), 1<<16)
	check(SecureToken(), 32)
	for _, n := range []int{0, -1, -1000, math.MinInt32} {
		check(GenerateSecureToken(n), 32) // anything below one asks for the default
	}

	// A token is fresh every time, on both paths.
	seen := make(map[string]struct{})
	for range 500 {
		seen[GenerateSecureToken(32)] = struct{}{}
		seen[GenerateSecureToken(33)] = struct{}{}
	}
	require.Len(t, seen, 1000)
}

// Test_Fold_StdlibParity_Exhaustive compares the case-insensitive searches with
// the stdlib searches on lower-cased copies, which is what ASCII folding means,
// over every haystack and needle of a small alphabet that holds both cases and
// the bytes next to the letter ranges.
func Test_Fold_StdlibParity_Exhaustive(t *testing.T) {
	t.Parallel()
	check := func(s, needle string) {
		t.Helper()
		ls, ln := strings.ToLower(s), strings.ToLower(needle)
		want := strings.Index(ls, ln)
		if got := IndexFold(s, needle); got != want {
			t.Fatalf("IndexFold(%q, %q) = %d, want %d", s, needle, got, want)
		}
		if got := IndexFold([]byte(s), needle); got != want {
			t.Fatalf("IndexFold(bytes %q, %q) = %d, want %d", s, needle, got, want)
		}
		if got := ContainsFold(s, needle); got != (want >= 0) {
			t.Fatalf("ContainsFold(%q, %q) = %v, want %v", s, needle, got, want >= 0)
		}
		if got, want := HasPrefixFold(s, needle), strings.HasPrefix(ls, ln); got != want {
			t.Fatalf("HasPrefixFold(%q, %q) = %v, want %v", s, needle, got, want)
		}
		if got, want := HasSuffixFold(s, needle), strings.HasSuffix(ls, ln); got != want {
			t.Fatalf("HasSuffixFold(%q, %q) = %v, want %v", s, needle, got, want)
		}
		if got, want := EqualFold(s, needle), strings.EqualFold(s, needle); got != want {
			t.Fatalf("EqualFold(%q, %q) = %v, want %v", s, needle, got, want)
		}
	}

	// Short haystacks, which take the byte loops.
	forEachString("aA@[-", 5, func(s string) {
		forEachString("aA@[-", 3, func(needle string) { check(s, needle) })
	})
	// Haystacks of a word and more, where the word loops and the tails are.
	forEachString("aA", 12, func(s string) {
		forEachString("aA", 6, func(needle string) { check(s, needle) })
	})
	// Needles that are longer than a word.
	long := "abcdefghijklmnopqrstuvwxyz"
	for n := 0; n <= len(long); n++ {
		for start := 0; start+n <= len(long); start++ {
			needle := long[start : start+n]
			check(strings.ToUpper(long)+"-"+long, needle)
			check(long[:start]+strings.ToUpper(needle)+long[start+n:], needle)
			check(long, strings.ToUpper(needle))
		}
	}
}
