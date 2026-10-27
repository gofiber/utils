package utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// floatGrammar is the accepted syntax; strconv also takes inf, nan and hex floats.
var floatGrammar = regexp.MustCompile(`^[+-]?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?$`)

// wantFloat is the oracle: strconv while it is right (up to 700 bytes), then big.Rat,
// which rounds exactly.
func wantFloat(tb testing.TB, s string, bitSize int) (float64, error) {
	tb.Helper()

	if len(s) <= 700 {
		return strconv.ParseFloat(s, bitSize)
	}
	mantissa, exp, _ := strings.Cut(strings.ToLower(strings.TrimLeft(s, "+-")), "e")
	intPart, fracPart, _ := strings.Cut(mantissa, ".")
	digits := strings.TrimLeft(intPart+fracPart, "0")
	trimmed := strings.TrimRight(digits, "0")

	const maxExp = 1 << 40
	bigExp := new(big.Int)
	if exp != "" {
		bigExp.SetString(exp, 10)
	}
	e := int64(maxExp)
	switch {
	case bigExp.CmpAbs(big.NewInt(maxExp)) <= 0:
		e = bigExp.Int64()
	case bigExp.Sign() < 0:
		e = -maxExp
	}
	e += int64(len(digits)-len(trimmed)) - int64(len(fracPart))

	n := int64(len(trimmed))
	var f float64
	switch {
	case trimmed == "" || n+e < -400:
	case n-1+e > 400:
		f = math.Inf(1)
	default:
		num, ok := new(big.Int).SetString(trimmed, 10)
		require.True(tb, ok, "big.Int cannot read %.40q", trimmed)
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(max(e, -e)), nil)
		r := new(big.Rat)
		if e >= 0 {
			r.SetInt(num.Mul(num, scale))
		} else {
			r.SetFrac(num, scale)
		}
		if bitSize == 32 {
			f32, _ := r.Float32()
			f = float64(f32)
		} else {
			f, _ = r.Float64()
		}
	}
	if strings.HasPrefix(s, "-") {
		f = math.Copysign(f, -1) // a zero Rat has no sign
	}
	if math.IsInf(f, 0) {
		return f, &strconv.NumError{Func: "ParseFloat", Num: s, Err: strconv.ErrRange}
	}
	return f, nil
}

func requireFloatMatches(tb testing.TB, s string) {
	tb.Helper()

	if !floatGrammar.MatchString(s) {
		_, err := ParseFloat64(s)
		require.ErrorIs(tb, err, strconv.ErrSyntax, "ParseFloat64(%q)", s)
		_, err = ParseFloat32(s)
		require.ErrorIs(tb, err, strconv.ErrSyntax, "ParseFloat32(%q)", s)
		return
	}

	want64, wantErr64 := wantFloat(tb, s, 64)
	got64, err64 := ParseFloat64(s)
	gotB64, errB64 := ParseFloat64([]byte(s))
	requireFloatResult(tb, "ParseFloat64", s, want64, wantErr64, got64, err64)
	requireFloatResult(tb, "ParseFloat64(bytes)", s, want64, wantErr64, gotB64, errB64)

	want32, wantErr32 := wantFloat(tb, s, 32)
	got32, err32 := ParseFloat32(s)
	gotB32, errB32 := ParseFloat32([]byte(s))
	requireFloatResult(tb, "ParseFloat32", s, want32, wantErr32, float64(got32), err32)
	requireFloatResult(tb, "ParseFloat32(bytes)", s, want32, wantErr32, float64(gotB32), errB32)
}

func requireFloatResult(tb testing.TB, fn, s string, want float64, wantErr error, got float64, err error) {
	tb.Helper()

	if wantErr != nil {
		// Only range errors occur inside the grammar; the reference returns Inf, the parsers 0.
		require.ErrorIs(tb, wantErr, strconv.ErrRange, "%s(%.40q): reference", fn, s)
		require.ErrorIs(tb, err, strconv.ErrRange, "%s(%.40q)", fn, s)
		require.Zero(tb, got, "%s(%.40q) returned a value with an error", fn, s)
		return
	}
	require.NoError(tb, err, "%s(%.40q)", fn, s)
	require.Equal(tb, math.Float64bits(want), math.Float64bits(got),
		"%s(%.40q) = %v (%#x), want %v (%#x)", fn, s, got, math.Float64bits(got), want, math.Float64bits(want))
}

func Test_ParseFloat_Corpus(t *testing.T) {
	t.Parallel()

	corpus := []string{
		"0.1", "0.2", "0.3", "0.7", "0.9631", "1.1", "2.675", "9.95", "123.456", "1328.3",
		"88.9438", "1024.62342", "4740156.50398", "8544.4627727", "331.075823", "0.30000000000000004",
		"0", "-0", "+0", "0.0", "-0.0", "0e0", "-0e0", "0e400", "-0e400", "0.000e-400", "00000", ".0", "-.0e5",
		"1.", ".5", "-.5", "+.5e1", "5.e1", "1e2", "1E2", "1e+2", "1e-2", "-1.2e3", "+123.45", "0001.50", "1.0000000000",
		"16777215", "16777216", "16777217", "16777218", "33554433",
		"9007199254740991", "9007199254740992", "9007199254740993", "9007199254740994", "9007199254740995",
		"18446744073709551615", "18446744073709551616", "18446744073709551617", "99999999999999999999",
		"1e10", "1e11", "1e22", "1e23", "123e-22", "123e-23", "1e-10", "1e-11", "1e-22", "1e-23",
		"9007199254740992e22", "9007199254740992e-22", "9007199254740993e-22",
		"8.41e21", "5e-324", "4e-324", "2e-324", "2.4703282292062327e-324", "2.4703282292062328e-324",
		"4.9406564584124654e-324", "2.2250738585072014e-308", "2.2250738585072011e-308", "2.2250738585072012e-308",
		"1.7976931348623157e308", "1.7976931348623158e308", "1.7976931348623159e308", "1.797693134862315807e308",
		"3.4028234e38", "3.4028235e38", "3.4028236e38", "3.40282356779733661637539395458142568447e38",
		"3.4028235677973366e38", "3.4028235677973367e38", "1.1754943508222875e-38", "1.1754942807573643e-38",
		"1.797693134862315708145274237317043567981e308", "3.4028234663852885981170418348451692544e38",
		"1e-45", "7e-46", "7.006492321624085e-46", "7.006492321624086e-46", "1.401298464324817e-45",
		"1.00000017881393432617187499", "1.00000017881393432617187500", "1.00000017881393432617187501",
		"0.123456789012345678901234567890", "123456789012345678901234567890", "1234567890123456789.0123456789",
		"0.30000000000000004440892098500626", "0.1000000000000000055511151231257827",
		"1." + strings.Repeat("0", 40) + "1", "1." + strings.Repeat("0", 40),
		"9" + strings.Repeat("0", 30) + "e-30", strings.Repeat("9", 400), strings.Repeat("9", 400) + "e-390",
		"0." + strings.Repeat("0", 400) + "1", "1" + strings.Repeat("0", 400) + "e-400",
		"1e309", "1e308", "1e-324", "1e-400", "1e400", "1e-999999999999", "1e999999999999", "0e999999999999",
		"1e99999999999999999999", "1e-99999999999999999999", "1e3000", "1e-3000", "-1e3000", "-1e-3000",
		"1e1125899906842623", "1e1125899906842624", "1e1125899906842625", "1e-1125899906842625",
		"1e18446744073709551615", "1e18446744073709551616", "1e18446744073709551621", "1e-18446744073709551621",
		// Past what strconv reads right: 800 digits, five digits of an exponent.
		"1" + strings.Repeat("0", 800) + "e-800", "1" + strings.Repeat("0", 800) + "e-791",
		"1" + strings.Repeat("0", 799) + "1e-491", "1" + strings.Repeat("0", 1000) + "e-500",
		"1" + strings.Repeat("0", 1000) + "e-940", "1" + strings.Repeat("0", 1000) + "e-1000",
		strings.Repeat("9", 1500) + "e-1500", "0." + strings.Repeat("0", 1000) + strings.Repeat("9", 1000),
		"0." + strings.Repeat("0", 9691) + "1e100000", "0." + strings.Repeat("0", 11000) + "1e100000",
		"1" + strings.Repeat("0", 10000) + "e-100000", "0." + strings.Repeat("0", 20000) + "1" + strings.Repeat("0", 21) + "1e20001",
		// Not in the grammar: every one a syntax error.
		"", "+", "-", ".", "+.", "-.", "e5", ".e5", "1e", "1e+", "1e-", "1.2.3", "1e1.0", "1e5x", "abc", "1x", " 1", "1 ",
		"inf", "-inf", "+Inf", "infinity", "nan", "NaN", "0x1p-2", "0x10", "1_000", "1,5", "--1", "+-1", "1e--1", "１２",
	}
	for _, s := range corpus {
		requireFloatMatches(t, s)
	}
}

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

func Test_ParseFloat_Generated(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xf10a7, 0x5eed)) //nolint:gosec // deterministic test data
	for range 30_000 {
		requireFloatMatches(t, genFloatString(r))
	}
}

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
	for range 30_000 {
		f := float64(r.IntN(10_000_000)) / math.Pow10(r.IntN(8))
		s := strconv.FormatFloat(f, 'f', -1, 64)
		got, err := ParseFloat64(s)
		require.NoError(t, err, s)
		require.Equal(t, math.Float64bits(f), math.Float64bits(got), "ParseFloat64(%q)", s)
	}
}

// Test_ParseFloat32_Midpoints checks float32 halfway points cut to 12 to 17 digits, last digit ±1.
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
		requireFloatMatches(t, s)
	})
}

// Both sides of the exact path's edges: 2^24 and 2^53, exponents 10 and 22.
func Test_ParseFloat_ExactPathBounds(t *testing.T) {
	t.Parallel()

	for _, m := range []uint64{1, 7, 123456789, 1<<24 - 1, 1 << 24, 1<<24 + 1, 1<<53 - 1, 1 << 53, 1<<53 + 1, 9999999999999999} {
		for exp := int64(-25); exp <= 25; exp++ {
			in := strconv.FormatUint(m, 10) + "e" + strconv.FormatInt(exp, 10)
			requireFloatMatches(t, in)

			want64, err := strconv.ParseFloat(in, 64)
			require.NoError(t, err, in)
			if exactFloat64(m, exp) {
				require.Equal(t, math.Float64bits(want64), math.Float64bits(scale64(m, exp)), "scale64(%s)", in)
			}
			want32, err := strconv.ParseFloat(in, 32)
			if err != nil {
				continue
			}
			if got, ok := scale32(m, exp); ok {
				require.Equal(t, math.Float64bits(want32), math.Float64bits(got), "scale32(%s)", in)
			}
		}
	}
}

// Test_parseExp pins the exponent reader that ParseFloat and ConvertToBytes share.
func Test_parseExp(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		from int
		want int64
		ok   bool
	}{
		{"5", 0, 5, true},
		{"+5", 0, 5, true},
		{"-5", 0, -5, true},
		{"-0", 0, 0, true},
		{"0012", 0, 12, true},
		{"1e30", 2, 30, true},
		{"1e-30", 2, -30, true},
		{"", 0, 0, false},
		{"1e", 2, 0, false},
		{"+", 0, 0, false},
		{"-", 0, 0, false},
		{"1e+", 2, 0, false},
		{"5x", 0, 0, false},
		{"x5", 0, 0, false},
		{"--5", 0, 0, false},
		{"+-5", 0, 0, false},
		{"5-", 0, 0, false},
		{"5 ", 0, 0, false},
		{"5.0", 0, 0, false},
		{"1e3.5", 2, 0, false},
	} {
		got, ok := parseExp(tc.in, tc.from)
		require.Equal(t, tc.ok, ok, "%q from %d", tc.in, tc.from)
		require.Equal(t, tc.want, got, "%q from %d", tc.in, tc.from)
		got, ok = parseExp([]byte(tc.in), tc.from)
		require.Equal(t, tc.ok, ok, "%q from %d as bytes", tc.in, tc.from)
		require.Equal(t, tc.want, got, "%q from %d as bytes", tc.in, tc.from)
	}

	got, ok := parseExp("99999999999999999999999999", 0)
	require.True(t, ok)
	require.Equal(t, int64(maxParsedExp), got)
	got, ok = parseExp("-99999999999999999999999999", 0)
	require.True(t, ok)
	require.Equal(t, -int64(maxParsedExp), got)
	got, ok = parseExp(strings.Repeat("0", 100)+"7", 0)
	require.True(t, ok)
	require.Equal(t, int64(7), got)
}

func Test_ParseFloat_MantissaExponent(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xe15e1, 0x1e3a)) //nolint:gosec // deterministic test data
	for range 20_000 {
		var man uint64
		switch r.IntN(3) {
		case 0:
			man = r.Uint64()
		case 1:
			man = r.Uint64() >> r.IntN(64)
		default:
			man = r.Uint64N(1e17) + 1
		}
		exp10 := r.IntN(800) - 450
		if r.IntN(3) == 0 { // the subnormal range of float64 and float32
			exp10 = -345 + r.IntN(30)
			if r.IntN(2) == 0 {
				exp10 = -75 + r.IntN(45)
			}
		}
		requireFloatMatches(t, strconv.FormatUint(man, 10)+"e"+strconv.Itoa(exp10))
	}
}

// Where a 64-bit mantissa turns zero or infinite, and products that only look like a halfway point.
func Test_ParseFloat_RangeEdges(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		"18446744073709551615e-343", "18446744073709551615e-342", "1e308", "1e309", "18446744073709551615e308",
		"17976931348623157e292", "17976931348623159e292",
		"5e-324", "2e-324", "3e-324", "22250738585072011e-324", "22250738585072012e-324",
		"1e38", "1e39", "3.4028235e38", "3.4028236e38", "1e-45", "7e-46", "7.1e-46", "1.4e-45", "2.8e-45", "1.1754942e-38", "1.1754943e-38",
		"37e46", "41e-24", "74e46", "82e-24", "103e-60", "129e47",
	} {
		requireFloatMatches(t, in)
	}
}

// Exact halfway points up to 40-digit integers, and the decimals just above and below.
func Test_ParseFloat_Halfway(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x4a1f, 0x77)) //nolint:gosec // deterministic test data
	checked := 0
	for _, f := range []struct {
		mantBits int
		minExp   int
		maxExp   int
	}{
		{52, 44, 135},
		{23, 4, 100},
	} {
		for k := f.minExp; k <= f.maxExp; k++ {
			for range 24 {
				// (2m+1) * 2^(k-mantBits-1) is halfway between m * 2^(k-mantBits) and the next float.
				m := uint64(1)<<f.mantBits | r.Uint64N(1<<f.mantBits)
				tie := new(big.Float).SetPrec(128).SetMantExp(new(big.Float).SetUint64(2*m+1), k-f.mantBits-1)
				text := tie.Text('f', 80)
				text = strings.TrimRight(strings.TrimRight(text, "0"), ".")

				for _, in := range []string{text, belowText(t, text), aboveText(text)} {
					requireFloatMatches(t, in)
					checked++
				}
			}
		}
	}
	require.Greater(t, checked, 1_000)
}

// Roundings into the next binade and the edges of the normal range.
func Test_ParseFloat_Carry(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		"9007199254740991.5", "9007199254740991.75", "9007199254740991.25",
		"16777215.5", "16777215.75", "16777215.25", "16777216.5",
		"1.7976931348623157e308", "1.7976931348623158e308", "1.7976931348623159e308",
		"2.2250738585072014e-308", "2.2250738585072013e-308", "2.2250738585072011e-308",
		"3.4028235e38", "3.4028236e38", "1.17549435e-38", "1.1754943e-38",
	} {
		requireFloatMatches(t, in)
	}
}

// Which float32 inputs one exact multiply or divide rounds; it gives up at a halfway point.
func Test_scale32(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in       string
		mantissa uint64
		exp      int64
		decided  bool
	}{
		{"314159e-4", 314159, -4, true},
		{"314159265e-8", 314159265, -8, true},
		{"16777219", 16777219, 0, false}, // halfway: rounding twice may differ
		{"12345678901234567e-17", 12345678901234567, -17, false},
		{"1e-300", 1, -300, false},
		{"1e23", 1, 23, false},
	} {
		got, ok := scale32(tc.mantissa, tc.exp)
		require.Equal(t, tc.decided, ok, tc.in)
		if ok {
			want, err := strconv.ParseFloat(tc.in, 32)
			require.NoError(t, err, tc.in)
			require.Equal(t, math.Float64bits(want), math.Float64bits(got), tc.in)
		}
	}
}

// Integers and exact fractions against strconv, and where the table of powers of ten ends.
func Test_exactDyadic(t *testing.T) {
	t.Parallel()

	p5 := func(n int) uint64 { return pow10[n] >> n }
	for _, tc := range []struct {
		man     uint64
		exp     int64
		decided bool
	}{
		{1, 19, true},
		{1, 20, false},
		{3, 0, true},
		{math.MaxUint64, 0, true},
		{math.MaxUint64, 1, false},
		{p5(19), -19, true},
		{7 * p5(19), -19, true},
		{1, -19, false},
		{p5(19), -20, false},
		{p5(18) * 3, -18, true},
		{p5(18) * 3, -19, false},
		{math.MaxUint64 - 1, -1, false},
		{15, -1, true},
		{45035996273704975, -1, true},
		{9007199254740993, 0, true},
		{16777217, 0, true},
		{1 << 63, 1, false},
		{1 << 62, 1, false},
		{1 << 61, 1, true},
	} {
		in := strconv.FormatUint(tc.man, 10) + "e" + strconv.FormatInt(tc.exp, 10)
		for _, bitSize := range []int{64, 32} {
			got, ok := exactDyadic(tc.man, tc.exp, bitSize)
			if bitSize == 64 {
				require.Equal(t, tc.decided, ok, in)
			}
			if !ok {
				continue
			}
			want, err := strconv.ParseFloat(in, bitSize)
			if err != nil {
				require.ErrorIs(t, err, strconv.ErrRange, in)
			}
			require.Equal(t, math.Float64bits(want), math.Float64bits(got), "%s as float%d", in, bitSize)
		}
	}
}

// Compares the rounding to 24 bits with big.Float, ties included.
func Test_roundBits(t *testing.T) {
	t.Parallel()

	check := func(u uint64) {
		want, _ := new(big.Float).SetPrec(24).SetMode(big.ToNearestEven).SetUint64(u).Float64()
		require.Equal(t, math.Float64bits(want), math.Float64bits(roundBits(u, 24)), "roundBits(%d, 24)", u)
	}
	for _, u := range []uint64{0, 1, 1<<24 - 1, 1 << 24, 1<<24 + 1, 1<<25 + 1, 1<<25 + 2, 1<<25 + 3, 1<<25 + 6, 1<<63 - 1, 1 << 63, math.MaxUint64, math.MaxUint64 - 1<<39} {
		check(u)
	}
	r := rand.New(rand.NewPCG(7, 11)) //nolint:gosec // deterministic test data
	for range 100_000 {
		u := r.Uint64() >> r.IntN(64)
		check(u)
		check(u&^(1<<r.IntN(40)-1) | 1<<r.IntN(40)) // near a tie
	}
}

// Numbers around the length where strconv gets the text itself (maxRawLen).
func Test_ParseFloat_RawBoundary(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xb0d, 0x700)) //nolint:gosec // deterministic test data
	for n := maxRawLen - 6; n <= maxRawLen+6; n++ {
		for range 40 {
			digits := make([]byte, n)
			for i := range digits {
				digits[i] = byte('0' + r.IntN(10))
			}
			switch r.IntN(3) {
			case 0:
				for i := 2 + r.IntN(30); i < n; i++ {
					digits[i] = '0'
				}
			case 1:
				for i := range r.IntN(n - 20) {
					digits[i] = '0'
				}
			}
			point := r.IntN(n)
			text := string(digits[:point]) + "." + string(digits[point:])
			requireFloatMatches(t, text)
			requireFloatMatches(t, text+"e"+strconv.Itoa(r.IntN(700)-350))
			requireFloatMatches(t, "-"+text[:n-1]) // the sign counts toward the length too
		}
	}
}

// Random long numbers with exponents in and just outside the float64 range.
func Test_ParseFloat_Long(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x10c9, 0xd161)) //nolint:gosec // deterministic test data
	for range 2_000 {
		digits := make([]byte, 20+r.IntN(1_480))
		for i := range digits {
			digits[i] = byte('0' + r.IntN(10))
		}
		switch r.IntN(4) {
		case 0:
			for i := 2 + r.IntN(30); i < len(digits); i++ {
				digits[i] = '0'
			}
		case 1:
			for i := 1; i < len(digits)-1; i++ {
				digits[i] = '0'
			}
		}
		point := r.IntN(len(digits) + 1)
		text := string(digits[:point])
		if point < len(digits) || r.IntN(4) == 0 {
			text += "." + string(digits[point:])
		}
		exp := -335 + r.IntN(651) - point
		if r.IntN(40) == 0 {
			exp = (1_000_001 + r.IntN(1_000_000_000)) * (1 - 2*r.IntN(2))
		}
		if r.IntN(2) == 0 {
			text += "e" + strconv.Itoa(exp)
		}
		if r.IntN(3) == 0 {
			text = "-" + text
		}
		requireFloatMatches(t, text)
	}
}

func halfwayText(j uint64, unit int) string {
	tie := new(big.Float).SetPrec(64).SetMantExp(new(big.Float).SetUint64(2*j+1), unit-1)
	text := tie.Text('f', 1100)
	if strings.Contains(text, ".") {
		text = strings.TrimRight(strings.TrimRight(text, "0"), ".")
	}
	return text
}

func belowText(t *testing.T, text string) string {
	t.Helper()

	whole, frac, _ := strings.Cut(text, ".")
	digits, ok := new(big.Int).SetString(whole+frac, 10)
	require.True(t, ok, text)
	below := new(big.Int).Sub(digits, big.NewInt(1)).String()
	for len(below) < len(frac)+1 {
		below = "0" + below
	}
	if frac != "" {
		below = below[:len(below)-len(frac)] + "." + below[len(below)-len(frac):]
	}
	return below
}

func aboveText(text string) string {
	if strings.Contains(text, ".") {
		return text + "1"
	}
	return text + ".1"
}

// Halfway points of up to 767 digits in every binade, and the decimals either side.
func Test_ParseFloat_HalfwayEverywhere(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xa11, 0x7e5)) //nolint:gosec // deterministic test data
	checked := 0
	for _, f := range []struct {
		mantBits, minUnit, maxExp int
	}{
		{52, -1074, 1023},
		{23, -149, 127},
	} {
		for k := f.minUnit; k <= f.maxExp; k++ {
			if k%3 != 0 && k > f.minUnit+30 && k < f.maxExp-30 {
				continue
			}
			unit := max(k-f.mantBits, f.minUnit)
			nbits := k - unit + 1
			for range 2 {
				j := uint64(1)<<(nbits-1) | r.Uint64N(uint64(1)<<(nbits-1))
				text := halfwayText(j, unit)
				for _, in := range []string{text, belowText(t, text), aboveText(text)} {
					requireFloatMatches(t, in)
					checked++
				}
			}
		}
		// Halfway between zero and the smallest subnormal.
		text := halfwayText(0, f.minUnit)
		for _, in := range []string{text, belowText(t, text), aboveText(text)} {
			requireFloatMatches(t, in)
		}
	}
	require.Greater(t, checked, 1_000)
}

func Test_ParseFloat_Subnormals(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x5b, 0x40)) //nolint:gosec // deterministic test data
	for range 20_000 {
		f := math.Float64frombits(r.Uint64() & (1<<53 - 1))
		for _, s := range []string{strconv.FormatFloat(f, 'g', -1, 64), strconv.FormatFloat(f, 'e', 16, 64)} {
			got, err := ParseFloat64(s)
			require.NoError(t, err, s)
			require.Equal(t, math.Float64bits(f), math.Float64bits(got), s)
		}
		g := math.Float32frombits(r.Uint32() & (1<<24 - 1))
		s := strconv.FormatFloat(float64(g), 'g', -1, 32)
		got, err := ParseFloat32(s)
		require.NoError(t, err, s)
		require.Equal(t, math.Float32bits(g), math.Float32bits(got), s)
	}
}

// Halfway above the largest float of each size, which rounds to infinity, and its neighbors.
func Test_ParseFloat_OverflowEdge(t *testing.T) {
	t.Parallel()

	for _, f := range []struct {
		j    uint64
		unit int
	}{
		{1<<53 - 1, 971},
		{1<<24 - 1, 104},
	} {
		text := halfwayText(f.j, f.unit)
		for _, in := range []string{text, belowText(t, text), aboveText(text)} {
			requireFloatMatches(t, in)
		}
	}
}

// Digits after 1844674407370955161, the first mantissa that takes no more digits.
func Test_ParseFloat_Cutoff(t *testing.T) {
	t.Parallel()

	for _, prefix := range []string{"1844674407370955159", "1844674407370955160", "1844674407370955161", "1844674407370955162"} {
		for _, d := range "0569" {
			for _, tail := range []string{"", "0", "1", "99"} {
				digits := prefix + string(d) + tail
				for point := 0; point <= len(digits); point++ {
					for _, exp := range []string{"", "e-3", "e5", "e300", "e-310"} {
						requireFloatMatches(t, digits[:point]+"."+digits[point:]+exp)
					}
				}
			}
		}
	}
}

// 11 million digits offset an 8-digit exponent, which a reader that stops at 2^20 gets wrong.
func Test_ParseFloat_HugeInputs(t *testing.T) {
	t.Parallel()

	zeros := strings.Repeat("0", 11_000_000)
	inf := math.Inf(1)
	for _, tc := range []struct {
		head, tail string
		want       float64
	}{
		{"0.", "1e11000000", 0.1},
		{"0.", "1e11000002", 10},
		{"1", "e-11000000", 1},
		{"15", "e-11000001", 1.5},
		{"1", "e-99999999", 0},
		{"1", "e-110000000", 0},
		{"0.", "1e99999999", inf},
		{"0.", "1e110000000", inf},
	} {
		in := tc.head + zeros + tc.tail
		// The error text carries the whole input, hence %.60v.
		for name, parse := range map[string]func() (float64, error){
			"string": func() (float64, error) { return ParseFloat64(in) },
			"bytes":  func() (float64, error) { return ParseFloat64([]byte(in)) },
		} {
			f, err := parse()
			if math.IsInf(tc.want, 1) {
				if !errors.Is(err, strconv.ErrRange) {
					t.Fatalf("%s %s: want a range error, got %.60v", name, tc.tail, err)
				}
				require.Zero(t, f, name)
				continue
			}
			if err != nil {
				t.Fatalf("%s %s: unexpected error %.60v", name, tc.tail, err)
			}
			require.Equal(t, math.Float64bits(tc.want), math.Float64bits(f), "%s %s", name, tc.tail)
		}
	}
}

// The slow paths allocate nothing beyond the error itself.
func Test_ParseFloat_SlowPathAllocs(t *testing.T) {
	long := strings.Repeat("9", 1_000) + ".5"
	for _, tc := range []struct {
		in     string
		allocs int
	}{
		{"0.123456789012345678901234567890", 0},
		{"9007199254740993", 0},
		{"1e23", 0},
		{"5e-324", 0},
		{"1" + strings.Repeat("0", 800) + "e-800", 0},
		{"1e309", 1},
		{long, 1},
	} {
		// A byte slice also copies itself into the error.
		bytesAllocs := tc.allocs
		if tc.allocs > 0 {
			bytesAllocs++
		}
		b := []byte(tc.in)
		var err error
		require.Equal(t, tc.allocs, int(testing.AllocsPerRun(20, func() { _, err = ParseFloat64(tc.in) })), "ParseFloat64(%.40q)", tc.in)
		require.Equal(t, bytesAllocs, int(testing.AllocsPerRun(20, func() { _, err = ParseFloat64(b) })), "ParseFloat64(bytes %.40q)", tc.in)
		if tc.allocs == 0 {
			require.NoError(t, err)
		}
	}
}

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
		{"tie", "9007199254740993"},
		{"padded", "1.0000000000000000"},
		{"subnormal", "5e-324"},
		{"digits-21", "0.123456789012345678901"},
		{"digits-40", "3.141592653589793238462643383279502884197"},
	} {
		b.Run(tc.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseFloat64(tc.in); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(tc.name+"/fiber_bytes", func(b *testing.B) {
			in := []byte(tc.in)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseFloat64(in); err != nil {
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
