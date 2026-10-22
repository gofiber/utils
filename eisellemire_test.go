package utils

import (
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
