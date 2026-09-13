package caseconv

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

// refLower and refUpper are the scalar ASCII folds every test measures against;
// bytes >= 0x80 pass through, as both the tables and the kernels promise.
func refLower(b []byte) []byte {
	out := make([]byte, len(b))
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		out[i] = c
	}
	return out
}

func refUpper(b []byte) []byte {
	out := make([]byte, len(b))
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out[i] = c
	}
	return out
}

func refFirstIndex(b []byte, lo, hi byte) int {
	for i, c := range b {
		if c >= lo && c <= hi {
			return i
		}
	}
	return -1
}

// assertCaseParity drives the whole interface over one input; the randomized and
// fuzz tests share it.
//
// Expectations are computed before each call and src is compared against an
// untouched copy after: the converters promise to read src only, and an illegal
// write there is itself a case-fold, so a reference read afterwards would match
// the corruption rather than catch it.
func assertCaseParity(t *testing.T, src []byte) {
	t.Helper()

	untouched := append([]byte(nil), src...)
	wantLower, wantUpper := refLower(src), refUpper(src)
	wantUp, wantLow := refFirstIndex(src, 'A', 'Z'), refFirstIndex(src, 'a', 'z')

	iUp, iLow := FirstUpperIndex(src), FirstLowerIndex(src)
	require.Equal(t, wantUp, iUp, "FirstUpperIndex(%q)", src)
	require.Equal(t, wantLow, iLow, "FirstLowerIndex(%q)", src)

	// The From converters require a word of input and the scanner's own index.
	if len(src) >= WordLen {
		if iUp >= 0 {
			require.Equal(t, string(wantLower), string(ToLowerFrom(src, iUp)), "ToLowerFrom(%q, %d)", src, iUp)
			require.Equal(t, string(untouched), string(src), "ToLowerFrom wrote through src")
		}
		if iLow >= 0 {
			require.Equal(t, string(wantUpper), string(ToUpperFrom(src, iLow)), "ToUpperFrom(%q, %d)", src, iLow)
			require.Equal(t, string(untouched), string(src), "ToUpperFrom wrote through src")
		}
	}

	lowered := append([]byte(nil), src...)
	ToLowerInPlace(lowered)
	require.Equal(t, string(wantLower), string(lowered), "ToLowerInPlace(%q)", untouched)

	uppered := append([]byte(nil), src...)
	ToUpperInPlace(uppered)
	require.Equal(t, string(wantUpper), string(uppered), "ToUpperInPlace(%q)", untouched)
}

// Test_CaseConv_Randomized covers sub-word inputs, every %WordLen tail, the
// block loop and its overlapping tail, with non-ASCII bytes mixed in.
func Test_CaseConv_Randomized(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic test data

	for range 4000 {
		src := make([]byte, rng.Intn(101))
		for i := range src {
			src[i] = byte(rng.Intn(256))
		}
		assertCaseParity(t, src)
	}
}

// Test_CaseConv_Lengths walks every length with the changing byte at every
// position, so each word and block offset is the first change once and the
// overlapping tail store meets every residue. The letter-free filler carries a
// byte >= 0x80 every fourth offset: without MatchRangeMask's high-bit clear,
// 0xE9 falsely matches 'a'..'z' and this fails.
func Test_CaseConv_Lengths(t *testing.T) {
	t.Parallel()

	for n := range 71 {
		base := make([]byte, n)
		for i := range base {
			base[i] = "40-\xe9/~:9"[i%8] // no letters at all, incl. >= 0x80
		}
		// The untouched filler drives both scanners to -1.
		assertCaseParity(t, base)

		for pos := range n {
			lo := append([]byte(nil), base...)
			lo[pos] = 'A' + byte(pos%26)
			if pos+1 < n {
				lo[pos+1] = 0xE9
			}
			assertCaseParity(t, lo)

			up := append([]byte(nil), base...)
			up[pos] = 'a' + byte(pos%26)
			if pos+1 < n {
				up[pos+1] = 0xE9
			}
			assertCaseParity(t, up)
		}
	}
}

// Test_ConvertFrom_PrefixIsTakenOnTrust pins the prefix skip, which no
// in-contract equality oracle can see: there the skipped bytes are unchanged
// either way. So these pass a first above the true first change -- the case the
// doc warns mis-converts -- and fail if the alignment or prefix copy is dropped.
// The results are not a clean unconverted prefix, because the overlapping tail
// store folds src back over part of it.
func Test_ConvertFrom_PrefixIsTakenOnTrust(t *testing.T) {
	t.Parallel()

	// The tail store folds src[n-WordLen:n] over dst, so only byte 0 survives.
	require.Equal(t, "Aaaaaaaab", string(ToLowerFrom([]byte("AAAAAAAAB"), WordLen)))
	require.Equal(t, "aAAAAAAAB", string(ToUpperFrom([]byte("aaaaaaaab"), WordLen)))

	// An unaligned first aligns down, so its whole word folds after all.
	require.Equal(t, "aaaaaaaab", string(ToLowerFrom([]byte("AAAAAAAAB"), WordLen-1)))
	require.Equal(t, "AAAAAAAAB", string(ToUpperFrom([]byte("aaaaaaaab"), WordLen-1)))
}

// FuzzCaseConv is the property test AGENTS.md asks of bit-twiddling code; the
// seeds pin the word and block boundaries and carry non-ASCII at word lengths
// and above, where the kernels rather than the tables fold.
func FuzzCaseConv(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("a"))
	f.Add([]byte("Content-Type"))
	f.Add([]byte("0123456A"))
	f.Add([]byte("0123456789abcdefghijklmnopqrstuvW"))
	f.Add([]byte("\xe9A\xe9a"))
	f.Add([]byte("0123456789\xe9\xc3\xa9ABCDEFGHIJKLMNOPQRST"))
	f.Add([]byte("\xffA\xff\xffBCDEFGH\x80ijklmnop\x7f"))
	f.Add([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	f.Fuzz(func(t *testing.T, src []byte) {
		assertCaseParity(t, src)
	})
}
