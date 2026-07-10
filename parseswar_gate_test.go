package utils

// Benchmark gate for P7 (8-digit SWAR parse step), comparing the integrated
// ParseUint against a same-process copy of the previous scalar
// implementation so VM noise cannot skew the verdict.
//
// P6 (IsIPv4 SWAR pre-validation) was prototyped here and DROPPED: it lost
// on every shape (valid-short +10%, valid-long +19%, invalid +21%, garbage
// +49%) because the scalar octet parse still has to run after the word
// checks — the assist is purely additive on <= 15-byte inputs.

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/gofiber/utils/v2/swar"
	"github.com/stretchr/testify/require"
)

// refParseUintScalar mirrors ParseUint exactly as it was before the 8-digit
// SWAR step landed in parseDigits.
func refParseUintScalar(s string) (uint64, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: "ParseUint", Num: "", Err: strconv.ErrSyntax}
	}
	var n uint64
	digits := 0
	for i := 0; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, &strconv.NumError{Func: "ParseUint", Num: s, Err: strconv.ErrSyntax}
		}
		d := uint64(c)
		if digits >= 19 && (n > maxUint64Cutoff || (n == maxUint64Cutoff && d > maxUint64Cutlim)) {
			return 0, &strconv.NumError{Func: "ParseUint", Num: s, Err: strconv.ErrRange}
		}
		n = n*10 + d
		digits++
	}
	return n, nil
}

func Test_ParseSWAR_Gate(t *testing.T) {
	t.Parallel()
	// parse8Digits against formatted boundaries.
	for _, v := range []uint64{0, 1, 9, 10, 12345678, 99999999, 10000000, 90000009} {
		s := FormatUint(v)
		for len(s) < 8 {
			s = "0" + s
		}
		require.Equal(t, v, parse8Digits(swar.Load8(s, 0)), s)
	}

	rng := rand.New(rand.NewSource(3)) //nolint:gosec // deterministic test data
	for range 200000 {
		v := rng.Uint64() >> rng.Intn(60) // cover all digit counts
		s := strconv.FormatUint(v, 10)
		got, err := ParseUint(s)
		require.NoError(t, err, s)
		require.Equal(t, v, got, s)

		gotI, err := ParseInt(s)
		wantI, wantErr := strconv.ParseInt(s, 10, 64)
		require.Equal(t, wantErr == nil, err == nil, s)
		if wantErr == nil {
			// On range errors strconv returns the clamped value while utils
			// returns 0, so values are only comparable on success.
			require.Equal(t, wantI, gotI, s)
		}
	}

	// Error and overflow parity with the previous implementation.
	for _, s := range []string{
		"", "a", "12345678a", "123456781234567812345", "18446744073709551615",
		"18446744073709551616", "99999999999999999999", "00000000000000000001",
		"1234567x", "12345678\n2345678", "184467440737095516150", "0000000000000000000000042",
	} {
		want, wantErr := refParseUintScalar(s)
		got, gotErr := ParseUint(s)
		require.Equal(t, wantErr == nil, gotErr == nil, "%q", s)
		require.Equal(t, want, got, "%q", s)
	}
}

func Benchmark_ParseSWAR_Gate(b *testing.B) {
	nums := map[string]string{
		"1digit":  "8",
		"4digit":  "8080",
		"6digit":  "123456",
		"8digit":  "12345678",
		"12digit": "123456789012",
		"19digit": "1234567890123456789",
	}
	for name, s := range nums {
		b.Run("ParseUint/"+name+"/old-scalar", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := refParseUintScalar(s); err != nil {
					b.Fatal("failed to parse uint")
				}
			}
		})
		b.Run("ParseUint/"+name+"/integrated", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseUint(s); err != nil {
					b.Fatal("failed to parse uint")
				}
			}
		})
	}
}
