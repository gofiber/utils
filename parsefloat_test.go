package utils

import (
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// floatGrammar is the syntax the parsers accept; strconv also takes inf, nan and hex floats.
var floatGrammar = regexp.MustCompile(`^[+-]?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?$`)

// requireFloatMatchesStrconv compares both parsers and input types with strconv, bit for bit.
func requireFloatMatchesStrconv(tb testing.TB, s string) {
	tb.Helper()

	if !floatGrammar.MatchString(s) {
		_, err := ParseFloat64(s)
		require.ErrorIs(tb, err, strconv.ErrSyntax, "ParseFloat64(%q)", s)
		_, err = ParseFloat32(s)
		require.ErrorIs(tb, err, strconv.ErrSyntax, "ParseFloat32(%q)", s)
		return
	}

	want64, wantErr64 := strconv.ParseFloat(s, 64)
	got64, err64 := ParseFloat64(s)
	gotB64, errB64 := ParseFloat64([]byte(s))
	requireFloatResult(tb, "ParseFloat64", s, want64, wantErr64, got64, err64)
	requireFloatResult(tb, "ParseFloat64(bytes)", s, want64, wantErr64, gotB64, errB64)

	want32, wantErr32 := strconv.ParseFloat(s, 32)
	got32, err32 := ParseFloat32(s)
	gotB32, errB32 := ParseFloat32([]byte(s))
	requireFloatResult(tb, "ParseFloat32", s, want32, wantErr32, float64(got32), err32)
	requireFloatResult(tb, "ParseFloat32(bytes)", s, want32, wantErr32, float64(gotB32), errB32)
}

func requireFloatResult(tb testing.TB, fn, s string, want float64, wantErr error, got float64, err error) {
	tb.Helper()

	if wantErr != nil {
		// Inside the grammar only a range error is possible; strconv returns Inf, these 0.
		require.ErrorIs(tb, wantErr, strconv.ErrRange, "%s(%q): strconv", fn, s)
		require.ErrorIs(tb, err, strconv.ErrRange, "%s(%q)", fn, s)
		require.Zero(tb, got, "%s(%q) returned a value with an error", fn, s)
		return
	}
	require.NoError(tb, err, "%s(%q)", fn, s)
	require.Equal(tb, math.Float64bits(want), math.Float64bits(got),
		"%s(%q) = %v (%#x), strconv says %v (%#x)", fn, s, got, math.Float64bits(got), want, math.Float64bits(want))
}

func Test_ParseFloat_MatchesStrconv_Corpus(t *testing.T) {
	t.Parallel()

	corpus := []string{
		// Classic 1-ulp traps of mantissa * 10^-k.
		"0.1", "0.2", "0.3", "0.7", "0.9631", "1.1", "2.675", "9.95", "123.456", "1328.3",
		"88.9438", "1024.62342", "4740156.50398", "8544.4627727", "331.075823", "0.30000000000000004",
		// Zero and its spellings.
		"0", "-0", "+0", "0.0", "-0.0", "0e0", "-0e0", "0e400", "-0e400", "0.000e-400", "00000", ".0", "-.0e5",
		// Shapes of the grammar.
		"1.", ".5", "-.5", "+.5e1", "5.e1", "1e2", "1E2", "1e+2", "1e-2", "-1.2e3", "+123.45", "0001.50", "1.0000000000",
		// Integers across the exact range of float64 and float32.
		"16777215", "16777216", "16777217", "16777218", "33554433",
		"9007199254740991", "9007199254740992", "9007199254740993", "9007199254740994", "9007199254740995",
		"18446744073709551615", "18446744073709551616", "18446744073709551617", "99999999999999999999",
		// Edges of the exact path (10^22 for float64, 10^10 for float32).
		"1e10", "1e11", "1e22", "1e23", "123e-22", "123e-23", "1e-10", "1e-11", "1e-22", "1e-23",
		"9007199254740992e22", "9007199254740992e-22", "9007199254740993e-22",
		// Halfway and near-halfway cases.
		"8.41e21", "5e-324", "4e-324", "2e-324", "2.4703282292062327e-324", "2.4703282292062328e-324",
		"4.9406564584124654e-324", "2.2250738585072014e-308", "2.2250738585072011e-308", "2.2250738585072012e-308",
		"1.7976931348623157e308", "1.7976931348623158e308", "1.7976931348623159e308", "1.797693134862315807e308",
		"3.4028234e38", "3.4028235e38", "3.4028236e38", "3.40282356779733661637539395458142568447e38",
		"3.4028235677973366e38", "3.4028235677973367e38", "1.1754943508222875e-38", "1.1754942807573643e-38",
		"1e-45", "7e-46", "7.006492321624085e-46", "7.006492321624086e-46", "1.401298464324817e-45",
		"1.00000017881393432617187499", "1.00000017881393432617187500", "1.00000017881393432617187501",
		// Long digit strings: dropped digits decide the rounding.
		"0.123456789012345678901234567890", "123456789012345678901234567890", "1234567890123456789.0123456789",
		"0.30000000000000004440892098500626", "0.1000000000000000055511151231257827",
		"1." + strings.Repeat("0", 40) + "1", "1." + strings.Repeat("0", 40),
		"9" + strings.Repeat("0", 30) + "e-30", strings.Repeat("9", 400), strings.Repeat("9", 400) + "e-390",
		"0." + strings.Repeat("0", 400) + "1", "1" + strings.Repeat("0", 400) + "e-400",
		// Huge exponents.
		"1e309", "1e308", "1e-324", "1e-400", "1e400", "1e-999999999999", "1e999999999999", "0e999999999999",
		"1e99999999999999999999", "1e-99999999999999999999",
		// Not in the grammar: every one a syntax error.
		"", "+", "-", ".", "+.", "-.", "e5", ".e5", "1e", "1e+", "1e-", "1.2.3", "1e1.0", "1e5x", "abc", "1x", " 1", "1 ",
		"inf", "-inf", "+Inf", "infinity", "nan", "NaN", "0x1p-2", "0x10", "1_000", "1,5", "--1", "+-1", "1e--1", "１２",
	}
	for _, s := range corpus {
		requireFloatMatchesStrconv(t, s)
	}
}

// genFloatString draws a string of the grammar with digit counts around the
// 15, 17 and 19 digit limits and exponents that reach the fallback.
func genFloatString(r *rand.Rand) string {
	var b strings.Builder
	switch r.IntN(4) {
	case 0:
		b.WriteByte('-')
	case 1:
		b.WriteByte('+')
	}
	intDigits := r.IntN(22)
	fracDigits := r.IntN(22)
	if intDigits == 0 && fracDigits == 0 {
		intDigits = 1
	}
	for range intDigits {
		b.WriteByte(byte('0' + r.IntN(10)))
	}
	if fracDigits > 0 || r.IntN(8) == 0 {
		b.WriteByte('.')
		for range fracDigits {
			b.WriteByte(byte('0' + r.IntN(10)))
		}
	}
	if r.IntN(3) == 0 {
		b.WriteByte("eE"[r.IntN(2)])
		switch r.IntN(3) {
		case 0:
			b.WriteByte('-')
		case 1:
			b.WriteByte('+')
		}
		b.WriteString(strconv.Itoa(r.IntN(70)))
	}
	return b.String()
}

func Test_ParseFloat_MatchesStrconv_Generated(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xf10a7, 0x5eed)) //nolint:gosec // deterministic test data
	for range 30_000 {
		requireFloatMatchesStrconv(t, genFloatString(r))
	}
}

// Test_ParseFloat_ShortestRoundTrip parses the shortest decimal of random floats back to the same bits.
func Test_ParseFloat_ShortestRoundTrip(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xdec1, 0x1a1)) //nolint:gosec // deterministic test data
	for range 30_000 {
		f := math.Float64frombits(r.Uint64())
		if math.IsNaN(f) || math.IsInf(f, 0) {
			continue
		}
		s := strconv.FormatFloat(f, 'g', -1, 64)
		got, err := ParseFloat64(s)
		require.NoError(t, err, s)
		require.Equal(t, math.Float64bits(f), math.Float64bits(got), "ParseFloat64(%q)", s)
		gotB, err := ParseFloat64([]byte(s))
		require.NoError(t, err, s)
		require.Equal(t, math.Float64bits(f), math.Float64bits(gotB), "ParseFloat64(bytes %q)", s)
	}
	for range 30_000 {
		f := math.Float32frombits(r.Uint32())
		if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
			continue
		}
		s := strconv.FormatFloat(float64(f), 'g', -1, 32)
		got, err := ParseFloat32(s)
		require.NoError(t, err, s)
		require.Equal(t, math.Float32bits(f), math.Float32bits(got), "ParseFloat32(%q)", s)
	}
	// Decimal digits a person writes: few digits, a few places.
	for range 30_000 {
		f := float64(r.IntN(10_000_000)) / math.Pow10(r.IntN(8))
		s := strconv.FormatFloat(f, 'f', -1, 64)
		got, err := ParseFloat64(s)
		require.NoError(t, err, s)
		require.Equal(t, math.Float64bits(f), math.Float64bits(got), "ParseFloat64(%q)", s)
	}
}

// Test_ParseFloat32_Midpoints checks decimals around the halfway points between
// float32s, where rounding through a float64 can pick the wrong side. Each
// midpoint is cut to 12 to 17 digits, with the neighboring last digits too.
func Test_ParseFloat32_Midpoints(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x3d, 0x32)) //nolint:gosec // deterministic test data
	checked := 0
	for range 8_000 {
		f := math.Float32frombits(r.Uint32())
		next := math.Nextafter32(f, float32(math.Inf(1)))
		if math.IsNaN(float64(f)) || math.IsInf(float64(next), 0) || f < 0 {
			continue
		}
		// Exact in float64: both are 24-bit values.
		mid := (float64(f) + float64(next)) / 2
		for sig := 12; sig <= 17; sig++ {
			digits, exp10, ok := strings.Cut(strconv.FormatFloat(mid, 'e', sig-1, 64), "e")
			require.True(t, ok)
			digits = strings.Replace(digits, ".", "", 1)
			base, err := strconv.ParseInt(digits, 10, 64)
			require.NoError(t, err)
			e10, err := strconv.Atoi(exp10)
			require.NoError(t, err)
			// The digits moved the point sig-1 places, so the exponent drops by as much.
			exponent := strconv.Itoa(e10 - (sig - 1))
			for delta := int64(-1); delta <= 1; delta++ {
				in := strconv.FormatInt(base+delta, 10) + "e" + exponent
				want, wantErr := strconv.ParseFloat(in, 32)
				got, err := ParseFloat32(in)
				if wantErr != nil {
					require.Error(t, err, in)
					continue
				}
				require.NoError(t, err, in)
				require.Equal(t, math.Float32bits(float32(want)), math.Float32bits(got), "ParseFloat32(%q)", in)
				checked++
			}
		}
	}
	require.Greater(t, checked, 10_000)
}

func Test_ParseFloat_ErrorShape(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want error
	}{
		{"", strconv.ErrSyntax},
		{"-", strconv.ErrSyntax},
		{"1e", strconv.ErrSyntax},
		{"x", strconv.ErrSyntax},
		{"1e309", strconv.ErrRange},
		{"-1e309", strconv.ErrRange},
		{strings.Repeat("9", 400), strconv.ErrRange},
	} {
		v, err := ParseFloat64(tc.in)
		var numErr *strconv.NumError
		require.ErrorAs(t, err, &numErr, tc.in)
		require.Equal(t, "ParseFloat64", numErr.Func, tc.in)
		require.Equal(t, tc.in, numErr.Num, tc.in)
		require.ErrorIs(t, err, tc.want, tc.in)
		require.Zero(t, v, tc.in)

		v32, err := ParseFloat32(tc.in)
		require.ErrorAs(t, err, &numErr, tc.in)
		require.Equal(t, "ParseFloat32", numErr.Func, tc.in)
		require.Equal(t, tc.in, numErr.Num, tc.in)
		require.ErrorIs(t, err, tc.want, tc.in)
		require.Zero(t, v32, tc.in)
	}

	// An error from a byte slice must not alias it.
	in := []byte(strings.Repeat("9", 400))
	_, err := ParseFloat64(in)
	var numErr *strconv.NumError
	require.ErrorAs(t, err, &numErr)
	for i := range in {
		in[i] = 'x'
	}
	require.Equal(t, strings.Repeat("9", 400), numErr.Num)
}

func Test_ParseFloat_NegativeZero(t *testing.T) {
	t.Parallel()

	for _, s := range []string{"-0", "-0.0", "-0e5", "-.0", "-0e400", "-0." + strings.Repeat("0", 400)} {
		f, err := ParseFloat64(s)
		require.NoError(t, err, s)
		require.True(t, f == 0 && math.Signbit(f), "%q must be negative zero", s)
		f32, err := ParseFloat32(s)
		require.NoError(t, err, s)
		require.True(t, f32 == 0 && math.Signbit(float64(f32)), "%q must be negative zero", s)
	}
	for _, s := range []string{"0", "+0", "0.0", "0e400"} {
		f, err := ParseFloat64(s)
		require.NoError(t, err, s)
		require.False(t, math.Signbit(f), s)
	}
}

func FuzzParseFloat(f *testing.F) {
	for _, s := range []string{
		"0.3", "5e-324", "3.4028235e38", "1.7976931348623159e308", "9007199254740993", "-0",
		"0.30000000000000004", "123456789012345678901234567890", "1e-400", "1e400", ".5", "5.", "1e+", "0x1p-2",
		"inf", "nan", "1_0", "１", "1e99999999999999999999", "2.2250738585072011e-308",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		requireFloatMatchesStrconv(t, s)
	})
}

// Benchmark_ParseFloat_Shapes covers the shapes that decide which path a value takes.
func Benchmark_ParseFloat_Shapes(b *testing.B) {
	for _, tc := range []struct{ name, in string }{
		{"integer", "42"},
		{"half", "0.5"},
		{"decimal", "12345.6789"},
		{"coordinate", "-122.419416"},
		{"price", "19.99"},
		{"exponent", "1.5e10"},
		{"wide-15", "123456.789012345"},
		{"wide-17", "0.12345678901234567"},
		{"tiny", "1e-300"},
	} {
		b.Run(tc.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseFloat64(tc.in); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(tc.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := strconv.ParseFloat(tc.in, 64); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
