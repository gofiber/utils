package utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
