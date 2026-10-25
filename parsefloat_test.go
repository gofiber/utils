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

// floatGrammar is the syntax the parsers accept; strconv also takes inf, nan and hex floats.
var floatGrammar = regexp.MustCompile(`^[+-]?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?$`)

// wantFloat is the float a correct parser returns for the number s, in strconv's
// convention (±Inf and ErrRange on overflow). strconv is right up to 700 bytes; it
// keeps at most 800 digits, misplaces the point past them and reads five digits of
// an exponent, so longer numbers go to big.Rat, which rounds exactly.
func wantFloat(tb testing.TB, s string, bitSize int) (float64, error) {
	tb.Helper()

	if len(s) <= 700 {
		return strconv.ParseFloat(s, bitSize)
	}
	mantissa, exp, _ := strings.Cut(strings.ToLower(strings.TrimLeft(s, "+-")), "e")
	intPart, fracPart, _ := strings.Cut(mantissa, ".")
	digits := strings.TrimLeft(intPart+fracPart, "0")
	trimmed := strings.TrimRight(digits, "0")

	// The number is trimmed * 10^e. An exponent beyond 2^40 is as good as 2^40:
	// no input has that many digits to offset it.
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

	// It lies in [10^(n-1+e), 10^(n+e)): out of the float range one way or the other
	// whatever the digits, or big.Rat rounds it.
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

// requireFloatMatches compares both parsers and input types with the correct
// answer, bit for bit.
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
		// Inside the grammar only a range error is possible; the reference returns Inf, these 0.
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
		"1.797693134862315708145274237317043567981e308", "3.4028234663852885981170418348451692544e38", // the largest floats, digits cut off
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

// genFloatString draws a string of the grammar with digit counts around the
// 15, 17 and 19 digit limits.
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
// float32s, cut to 12 to 17 digits with the neighboring last digits.
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
		requireFloatMatches(t, s)
	})
}

// Test_ParseFloat_ExactPathBounds checks both sides of the edges of the exact path
// (2^24 and 2^53, exponents 10 and 22) against strconv.
func Test_ParseFloat_ExactPathBounds(t *testing.T) {
	t.Parallel()

	for _, m := range []uint64{1, 7, 123456789, 1<<24 - 1, 1 << 24, 1<<24 + 1, 1<<53 - 1, 1 << 53, 1<<53 + 1, 9999999999999999} {
		for exp := int64(-25); exp <= 25; exp++ {
			in := strconv.FormatUint(m, 10) + "e" + strconv.FormatInt(exp, 10)
			requireFloatMatches(t, in)

			want64, err := strconv.ParseFloat(in, 64)
			require.NoError(t, err, in)
			if got, ok := scaleExact(m, exp, 64); ok {
				require.Equal(t, math.Float64bits(want64), math.Float64bits(got), "scaleExact(%s, 64)", in)
			}
			want32, err := strconv.ParseFloat(in, 32)
			if err != nil {
				continue
			}
			if got, ok := scaleExact(m, exp, 32); ok {
				require.Equal(t, math.Float64bits(want32), math.Float64bits(got), "scaleExact(%s, 32)", in)
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

	// A long exponent keeps its sign and is clamped.
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

// Test_Pow10Tab derives each table entry again in 4096-bit floating point.
func Test_Pow10Tab(t *testing.T) {
	t.Parallel()

	require.Len(t, pow10Tab, pow10Max-pow10Min+1)

	const prec = 4096 // 10^347 has 1153 bits
	one := new(big.Float).SetPrec(prec).SetInt64(1)
	mask := new(big.Int).SetUint64(math.MaxUint64)
	for e := pow10Min; e <= pow10Max; e++ {
		p := new(big.Float).SetPrec(prec).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(max(e, -e))), nil))
		if e < 0 {
			p.Quo(one, p)
		}
		var frac big.Float
		p.MantExp(&frac) // p = frac * 2^exp with frac in [0.5, 1)
		m, _ := frac.SetMantExp(&frac, 128).Int(nil)

		got := pow10Tab[e-pow10Min]
		require.Equal(t, new(big.Int).Rsh(m, 64).Uint64(), got.hi, "1e%d high word", e)
		require.Equal(t, new(big.Int).And(m, mask).Uint64(), got.lo, "1e%d low word", e)
	}
}

// requireEiselLemire compares eiselLemire with strconv for both sizes (+Inf is
// strconv's answer to an overflow) and returns how many results it gave.
func requireEiselLemire(tb testing.TB, man uint64, exp10 int) int {
	tb.Helper()

	decided := 0
	in := strconv.FormatUint(man, 10) + "e" + strconv.Itoa(exp10)
	for _, bitSize := range []int{64, 32} {
		got, ok := eiselLemire(man, int64(exp10), bitSize)
		if !ok {
			continue
		}
		want, err := strconv.ParseFloat(in, bitSize)
		if err != nil {
			require.ErrorIs(tb, err, strconv.ErrRange, "eiselLemire(%q, %d)", in, bitSize)
		}
		require.Equal(tb, math.Float64bits(want), math.Float64bits(got), "eiselLemire(%q, %d)", in, bitSize)
		decided++
	}
	return decided
}

// Test_EiselLemire_MatchesStrconv compares random inputs, across the whole
// exponent range and into overflow and underflow, with strconv. It requires an
// answer for nearly all of them, so a function that always gives up fails.
func Test_EiselLemire_MatchesStrconv(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xe15e1, 0x1e3a)) //nolint:gosec // deterministic test data
	total, decided := 0, 0
	for range 40_000 {
		var man uint64
		switch r.IntN(3) {
		case 0:
			man = r.Uint64()
		case 1:
			man = r.Uint64() >> r.IntN(64)
		default:
			man = r.Uint64N(1e17) + 1
		}
		if man == 0 {
			continue
		}
		exp10 := r.IntN(2*pow10Max+100) - pow10Max - 50
		if r.IntN(3) == 0 { // the subnormal range of float64 and float32
			exp10 = -345 + r.IntN(30)
			if r.IntN(2) == 0 {
				exp10 = -75 + r.IntN(45)
			}
		}
		decided += requireEiselLemire(t, man, exp10)
		total += 2
	}
	require.GreaterOrEqual(t, decided, total*999/1000, "eiselLemire gave up on %d of %d", total-decided, total)
}

// Test_EiselLemire_Range covers the edges of the table: where a 64-bit mantissa
// becomes zero or infinite whatever its digits, and the subnormals just above.
func Test_EiselLemire_Range(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		man  uint64
		exp  int
		want float64
	}{
		{math.MaxUint64, pow10Min - 1, 0},
		{math.MaxUint64, pow10Min, 0x1p-1072},
		{1, pow10Max, 1e308},
		{1, pow10Max + 1, math.Inf(1)},
		{math.MaxUint64, pow10Max, math.Inf(1)},
		{17976931348623157, 292, math.MaxFloat64},
		{17976931348623159, 292, math.Inf(1)},
		{5, -324, 0x1p-1074},
		{2, -324, 0},
		{3, -324, 0x1p-1074},
		{22250738585072011, -324, 0x0.fffffffffffffp-1022},
		{22250738585072012, -324, 0x1p-1022},
	} {
		got, ok := eiselLemire(tc.man, int64(tc.exp), 64)
		in := strconv.FormatUint(tc.man, 10) + "e" + strconv.Itoa(tc.exp)
		require.True(t, ok, in)
		want, err := strconv.ParseFloat(in, 64)
		if err != nil {
			require.ErrorIs(t, err, strconv.ErrRange, in)
		}
		require.Equal(t, math.Float64bits(want), math.Float64bits(got), in)
		require.Equal(t, math.Float64bits(tc.want), math.Float64bits(got), in)
		requireEiselLemire(t, tc.man, tc.exp)
	}
	for _, in := range []string{"1e38", "1e39", "3.4028235e38", "3.4028236e38", "1e-45", "7e-46", "7.1e-46", "1.4e-45", "2.8e-45", "1.1754942e-38", "1.1754943e-38"} {
		requireFloatMatches(t, in)
	}
}

// Test_ParseFloat_Halfway parses the exact halfway points between floats, up to
// 40-digit integers, and the decimals just above and below each.
func Test_ParseFloat_Halfway(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x4a1f, 0x77)) //nolint:gosec // deterministic test data
	checked := 0
	for _, f := range []struct {
		mantBits int
		minExp   int
		maxExp   int
	}{
		{52, 44, 135}, // float64: from fractions of 20 digits to integers of 40
		{23, 4, 100},  // float32
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
				whole, frac, _ := strings.Cut(text, ".")
				if man, err := strconv.ParseUint(whole+frac, 10, 64); err == nil {
					for delta := range uint64(3) {
						requireEiselLemire(t, man+delta-1, -len(frac))
					}
				}
			}
		}
	}
	require.Greater(t, checked, 1_000)
}

// Test_EiselLemire_NotHalfway covers products that look like a halfway point in
// their top bits but are not one.
func Test_EiselLemire_NotHalfway(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		man uint64
		exp int64
	}{{37, 46}, {41, -24}, {74, 46}, {82, -24}, {103, -60}, {129, 47}} {
		got, ok := eiselLemire(tc.man, tc.exp, 64)
		in := strconv.FormatUint(tc.man, 10) + "e" + strconv.FormatInt(tc.exp, 10)
		require.True(t, ok, in)
		want, err := strconv.ParseFloat(in, 64)
		require.NoError(t, err, in)
		require.Equal(t, math.Float64bits(want), math.Float64bits(got), in)
	}
}

// Test_EiselLemire_Carry covers roundings into the next binade and the edges of the normal range.
func Test_EiselLemire_Carry(t *testing.T) {
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

// Test_ParseFloat_UsesEiselLemire pins scaleExact to answer past the exact range,
// halfway points and exact products included, and to give up only where it must.
func Test_ParseFloat_UsesEiselLemire(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in       string
		mantissa uint64
		exp      int64
		bitSize  int
		decided  bool
	}{
		{"12345678901234567e-17", 12345678901234567, -17, 64, true},
		{"1e-300", 1, -300, 64, true},
		{"17976931348623157e292", 17976931348623157, 292, 64, true},
		{"314159265e-8", 314159265, -8, 32, true},
		{"16777219", 16777219, 0, 32, true},                     // halfway, rounds up to the even float32
		{"9007199254740993", 9007199254740993, 0, 64, true},     // halfway, rounds down to the even float64
		{"4503599627370497.5", 45035996273704975, -1, 64, true}, // halfway with a fraction
		{"1.0000000000000000", 10000000000000000, -16, 64, true},
		{"1e23", 1, 23, 64, false}, // halfway, but 5^23 does not fit a uint64
	} {
		got, ok := scaleExact(tc.mantissa, tc.exp, tc.bitSize)
		require.Equal(t, tc.decided, ok, tc.in)
		if ok {
			want, err := strconv.ParseFloat(tc.in, tc.bitSize)
			require.NoError(t, err, tc.in)
			require.Equal(t, math.Float64bits(want), math.Float64bits(got), tc.in)
		}
	}
}

// Test_exactDyadic checks the integers and exact fractions it rounds against
// strconv, and the exponents where the table of powers of ten ends.
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

// Test_roundBits compares the rounding to 24 bits with big.Float, ties included.
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

// Test_ParseFloat_Long parses random numbers of 20 to 1,500 digits, their
// decimal point and exponent chosen to land in and just outside the float64 range.
func Test_ParseFloat_Long(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x10c9, 0xd161)) //nolint:gosec // deterministic test data
	for range 2_000 {
		digits := make([]byte, 20+r.IntN(1_480))
		for i := range digits {
			digits[i] = byte('0' + r.IntN(10))
		}
		switch r.IntN(4) {
		case 0: // a short significant part, then zeros
			for i := 2 + r.IntN(30); i < len(digits); i++ {
				digits[i] = '0'
			}
		case 1: // a digit at each end and zeros between
			for i := 1; i < len(digits)-1; i++ {
				digits[i] = '0'
			}
		}
		point := r.IntN(len(digits) + 1)
		text := string(digits[:point])
		if point < len(digits) || r.IntN(4) == 0 {
			text += "." + string(digits[point:])
		}
		// The decimal exponent of the leading digit lands between -335 and 315, or
		// the exponent is far past anything the digits can offset.
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

// halfwayText returns the exact decimal of (2j+1) * 2^(unit-1), which lies halfway
// between j * 2^unit and its successor.
func halfwayText(j uint64, unit int) string {
	tie := new(big.Float).SetPrec(64).SetMantExp(new(big.Float).SetUint64(2*j+1), unit-1)
	text := tie.Text('f', 1100)
	if strings.Contains(text, ".") {
		text = strings.TrimRight(strings.TrimRight(text, "0"), ".")
	}
	return text
}

// belowText returns the decimal text with its last digit one lower.
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

// aboveText returns the decimal text with a nonzero digit after its last.
func aboveText(text string) string {
	if strings.Contains(text, ".") {
		return text + "1"
	}
	return text + ".1"
}

// Test_ParseFloat_HalfwayEverywhere parses the exact halfway points between
// floats in every binade, from the subnormals to the largest, and the decimals
// just above and below each. They have up to 767 digits.
func Test_ParseFloat_HalfwayEverywhere(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xa11, 0x7e5)) //nolint:gosec // deterministic test data
	checked := 0
	for _, f := range []struct {
		mantBits, minUnit, maxExp int
	}{
		{52, -1074, 1023}, // float64
		{23, -149, 127},   // float32
	} {
		for k := f.minUnit; k <= f.maxExp; k++ {
			if k%3 != 0 && k > f.minUnit+30 && k < f.maxExp-30 {
				continue
			}
			unit := max(k-f.mantBits, f.minUnit) // floats of this binade are j * 2^unit
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

// Test_ParseFloat_Subnormals parses the shortest and the 17 digit decimals of
// random subnormals and the smallest normals back to the same bits.
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

// Test_ParseFloat_OverflowEdge parses the halfway point above the largest float of
// each size, which rounds to infinity, and the decimals next to it.
func Test_ParseFloat_OverflowEdge(t *testing.T) {
	t.Parallel()

	for _, f := range []struct {
		j    uint64
		unit int
	}{
		{1<<53 - 1, 971}, // float64: the largest is (2^53-1) * 2^971
		{1<<24 - 1, 104}, // float32: (2^24-1) * 2^104
	} {
		text := halfwayText(f.j, f.unit)
		for _, in := range []string{text, belowText(t, text), aboveText(text)} {
			requireFloatMatches(t, in)
		}
	}
}

// Test_ParseFloat_Cutoff sweeps the numbers that start around 1844674407370955161,
// the first mantissa that takes no more digits, followed by digits that would fit
// the mantissa before it, with the point and an exponent anywhere.
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

// Test_ParseFloat_HugeInputs parses numbers whose 11 million digits offset an
// exponent of 8 digits, which a reader that stops counting at 2^20 gets wrong.
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
		{"1", "e-99999999", 0},     // 1e-88999999
		{"1", "e-110000000", 0},    // 1e-99000000
		{"0.", "1e99999999", inf},  // 1e88999998
		{"0.", "1e110000000", inf}, // 1e98999999
	} {
		in := tc.head + zeros + tc.tail
		// The error text carries the whole input, so a failure prints a part of it.
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

// Test_ParseFloat_SlowPathAllocs pins what the slow paths do not allocate: a
// number that lost a digit, a halfway point and an overflow cost nothing beyond
// the error itself.
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
