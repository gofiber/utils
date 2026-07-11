package utils

// Tests and benchmarks for the 8-digit SWAR parse step in ParseUint/ParseInt.
//
// Historical gate notes: the IsIPv4/IsIPv6 SWAR pre-validation (handoff P6)
// was prototyped here and DROPPED — it lost on every shape (+10% to +49%)
// because the scalar octet parse still has to run after the word checks on
// <= 15-byte inputs. The scalar-baseline comparison that gated P7 lives in
// the introducing commit's message; ParseUint won -16% to -27% at 8+ digits
// with a one-branch cost on 1-digit inputs.

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/gofiber/utils/v2/swar"
	"github.com/stretchr/testify/require"
)

func Test_Parse8Digits_And_IsEightDigits(t *testing.T) {
	t.Parallel()
	for _, v := range []uint64{0, 1, 9, 10, 12345678, 99999999, 10000000, 90000009, 9999990} {
		s := FormatUint(v)
		for len(s) < 8 {
			s = "0" + s
		}
		w := swar.Load8(s, 0)
		require.True(t, isEightDigits(w), s)
		require.Equal(t, v, parse8Digits(w), s)
	}
	// isEightDigits must reject every byte value outside '0'..'9' in every
	// lane, including bytes >= 0xFA whose +6 carries into the next lane.
	for c := range 256 {
		lanes := []byte("11111111")
		for lane := range 8 {
			lanes[lane] = byte(c)
			got := isEightDigits(swar.Load8(lanes, 0))
			require.Equal(t, c >= '0' && c <= '9', got, "byte 0x%02x lane %d", c, lane)
			lanes[lane] = '1'
		}
	}
}

func Test_ParseSWAR_Parity(t *testing.T) {
	t.Parallel()
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

	// Error parity with strconv on syntax/overflow shapes that exercise the
	// SWAR word steps, including the scalar handoff around 16-19 digits.
	for _, s := range []string{
		"", "a", "12345678a", "123456781234567812345", "18446744073709551615",
		"18446744073709551616", "99999999999999999999", "00000000000000000001",
		"1234567x", "12345678\n2345678", "184467440737095516150",
		"0000000000000000000000042", "999999999999999999999999",
	} {
		want, wantErr := strconv.ParseUint(s, 10, 64)
		got, gotErr := ParseUint(s)
		require.Equal(t, wantErr == nil, gotErr == nil, "%q", s)
		if wantErr == nil {
			require.Equal(t, want, got, "%q", s)
		} else {
			require.Zero(t, got, "%q", s)
			// The error must keep its strconv shape: *NumError carrying the
			// function name and the ErrSyntax/ErrRange cause.
			var numErr *strconv.NumError
			require.ErrorAs(t, gotErr, &numErr, "%q", s)
			require.Equal(t, "ParseUint", numErr.Func, "%q", s)
			var wantNum *strconv.NumError
			require.ErrorAs(t, wantErr, &wantNum)
			require.ErrorIs(t, numErr.Err, wantNum.Err, "%q", s)
		}
	}
}

func Test_ParseInt_SWARPaths(t *testing.T) {
	t.Parallel()
	// Fast path: a non-digit inside the first word falls back to the scalar
	// loop, which reports the syntax error.
	for _, s := range []string{"1234x678", "12345678\x00", "123456789012345678x"} {
		_, err := ParseInt(s)
		require.Error(t, err, s)
		want, wantErr := strconv.ParseInt(s, 10, 64)
		require.Error(t, wantErr, s)
		require.Equal(t, want, int64(0), s) // strconv returns 0 on syntax errors too
	}

	// Signed inputs with 8+ digit runs route through parseSigned's word path.
	for _, s := range []string{
		"-123456789", "+123456789012", "-9223372036854775808", "+9223372036854775807",
		"-9223372036854775809", "-99999999999999999999", "-12345678a", "-000000001",
	} {
		got, gotErr := ParseInt(s)
		want, wantErr := strconv.ParseInt(s, 10, 64)
		require.Equal(t, wantErr == nil, gotErr == nil, s)
		if wantErr == nil {
			require.Equal(t, want, got, s)
		} else {
			require.Zero(t, got, s)
		}
	}

	// The same dispatch feeds the narrower signed parsers.
	v32, err := ParseInt32("-214748364")
	require.NoError(t, err)
	require.Equal(t, int32(-214748364), v32)
	_, err = ParseInt32("-21474836480")
	require.Error(t, err)
}

func Benchmark_ParseUint_DigitRuns(b *testing.B) {
	nums := map[string]string{
		"1digit":  "8",
		"4digit":  "8080",
		"6digit":  "123456",
		"8digit":  "12345678",
		"12digit": "123456789012",
		"19digit": "1234567890123456789",
	}
	for name, s := range nums {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseUint(s); err != nil {
					b.Fatal("failed to parse uint")
				}
			}
		})
	}
}
