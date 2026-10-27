// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"net"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_FunctionName(t *testing.T) {
	t.Parallel()
	require.Equal(t, "github.com/gofiber/utils/v2.Test_UUIDv4", FunctionName(Test_UUIDv4))
	require.Equal(t, "github.com/gofiber/utils/v2.Test_FunctionName.func1", FunctionName(func() {}))

	dummyint := 20
	require.Equal(t, "int", FunctionName(dummyint))

	require.Empty(t, FunctionName(nil))

	var nilFunc func()
	require.Empty(t, FunctionName(nilFunc))

	type myFunc func()
	var mf myFunc
	require.Empty(t, FunctionName(mf))

	var iface any = mf
	require.Empty(t, FunctionName(iface))

	var ptr *int
	require.Equal(t, "*int", FunctionName(ptr))

	type sampleStruct struct{}
	var s sampleStruct
	require.Equal(t, "utils.sampleStruct", FunctionName(s))
	require.Equal(t, "*utils.sampleStruct", FunctionName(&s))

	require.Contains(t, FunctionName(failingReader{}.Read), "failingReader.Read")
	require.Contains(t, FunctionName(failingReader.Read), "failingReader.Read")

	pc := reflect.ValueOf(Test_FunctionName).Pointer()
	require.Equal(t, runtime.FuncForPC(pc).Name(), FunctionName(Test_FunctionName))
}

func Test_UUIDv4(t *testing.T) {
	t.Parallel()
	res := UUIDv4()
	require.Len(t, res, 36)
	require.NotEqual(t, "00000000-0000-0000-0000-000000000000", res)
	require.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, res)
}

func Test_UUID(t *testing.T) {
	t.Parallel()
	seen := make(map[string]struct{})
	for range 64 {
		res := UUID()
		parsed, err := uuid.Parse(res)
		require.NoError(t, err)
		require.Equal(t, res, parsed.String())
		require.Equal(t, uuid.Version(4), parsed.Version())
		require.Equal(t, uuid.RFC4122, parsed.Variant())
		seen[res] = struct{}{}
	}
	require.Len(t, seen, 64)
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

// Test_UUIDv4_RandomFailure is not parallel: it swaps the random source of the uuid
// package, which parallel tests use too.
func Test_UUIDv4_RandomFailure(t *testing.T) {
	errDry := errors.New("entropy exhausted")
	uuid.SetRand(failingReader{err: errDry})
	t.Cleanup(func() { uuid.SetRand(nil) })

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		UUIDv4()
	}()
	err, ok := recovered.(error)
	require.True(t, ok, "UUIDv4 panicked with %v, want an error", recovered)
	require.ErrorIs(t, err, errDry)
	require.EqualError(t, err, "utils: failed to generate secure UUID: entropy exhausted")
	require.PanicsWithError(t, "utils: failed to generate secure UUID: entropy exhausted", func() { UUID() })

	uuid.SetRand(nil)
	require.Len(t, UUIDv4(), 36)
}

func Test_UUIDv4_Concurrency(t *testing.T) {
	t.Parallel()
	iterations := 1000
	var res string
	ch := make(chan string, iterations)
	results := make(map[string]string)
	for range iterations {
		go func() {
			ch <- UUIDv4()
		}()
	}
	for range iterations {
		res = <-ch
		results[res] = res
	}
	require.Len(t, results, iterations)
}

func Test_GenerateSecureToken(t *testing.T) {
	t.Parallel()
	token := GenerateSecureToken(32)
	require.Len(t, token, 43) // unpadded base64
	require.NotEmpty(t, token)

	token8 := GenerateSecureToken(8)
	require.Len(t, token8, 11)

	token16 := GenerateSecureToken(16)
	require.Len(t, token16, 22)

	token2 := GenerateSecureToken(32)
	require.NotEqual(t, token, token2)

	// A length below 1 means 32 bytes.
	tokenZero := GenerateSecureToken(0)
	require.Len(t, tokenZero, 43)

	tokenNegative := GenerateSecureToken(-1)
	require.Len(t, tokenNegative, 43)
}

func Test_GenerateSecureToken_Concurrency(t *testing.T) {
	t.Parallel()
	iterations := 1000
	ch := make(chan string, iterations)
	results := make(map[string]string)
	for range iterations {
		go func() {
			ch <- GenerateSecureToken(32)
		}()
	}
	for range iterations {
		res := <-ch
		results[res] = res
	}
	require.Len(t, results, iterations)
}

func Test_SecureToken(t *testing.T) {
	t.Parallel()
	token := SecureToken()
	require.Len(t, token, 43)
	require.NotEmpty(t, token)

	token2 := SecureToken()
	require.NotEqual(t, token, token2)
}

func saturate(v int64) int { return int(min(v, math.MaxInt)) }

func Test_ConvertToBytes(t *testing.T) {
	t.Parallel()
	require.Equal(t, 0, ConvertToBytes(""))
	require.Equal(t, 42, ConvertToBytes("42"))

	require.Equal(t, 0, ConvertToBytes(""))

	require.Equal(t, 42, ConvertToBytes("42"))
	require.Equal(t, 0, ConvertToBytes("0"))
	require.Equal(t, 1, ConvertToBytes("1"))
	require.Equal(t, 999, ConvertToBytes("999"))

	require.Equal(t, 42, ConvertToBytes("42b"))
	require.Equal(t, 42, ConvertToBytes("42B"))
	require.Equal(t, 42, ConvertToBytes("42 b"))
	require.Equal(t, 42, ConvertToBytes("42 B"))

	require.Equal(t, 42*1000, ConvertToBytes("42k"))
	require.Equal(t, 42*1000, ConvertToBytes("42K"))
	require.Equal(t, 42*1000, ConvertToBytes("42kb"))
	require.Equal(t, 42*1000, ConvertToBytes("42KB"))
	require.Equal(t, 42*1000, ConvertToBytes("42 kb"))
	require.Equal(t, 42*1000, ConvertToBytes("42 KB"))

	require.Equal(t, 42*1000000, ConvertToBytes("42M"))
	require.Equal(t, 42*1000000, ConvertToBytes("42m"))
	require.Equal(t, 42*1000000, ConvertToBytes("42MB"))
	require.Equal(t, 42*1000000, ConvertToBytes("42mb"))
	require.Equal(t, int(42.5*1000000), ConvertToBytes("42.5MB"))

	require.Equal(t, saturate(42*1000000000), ConvertToBytes("42G"))
	require.Equal(t, saturate(42*1000000000), ConvertToBytes("42g"))
	require.Equal(t, saturate(42*1000000000), ConvertToBytes("42GB"))
	require.Equal(t, saturate(42*1000000000), ConvertToBytes("42gb"))

	require.Equal(t, saturate(42*1000000000000), ConvertToBytes("42T"))
	require.Equal(t, saturate(42*1000000000000), ConvertToBytes("42t"))
	require.Equal(t, saturate(42*1000000000000), ConvertToBytes("42TB"))
	require.Equal(t, saturate(42*1000000000000), ConvertToBytes("42tb"))

	require.Equal(t, saturate(42*1000000000000000), ConvertToBytes("42P"))
	require.Equal(t, saturate(42*1000000000000000), ConvertToBytes("42p"))
	require.Equal(t, saturate(42*1000000000000000), ConvertToBytes("42PB"))
	require.Equal(t, saturate(42*1000000000000000), ConvertToBytes("42pb"))

	// E is a unit only before b or i, else an exponent marker.
	require.Equal(t, saturate(1000000000000000000), ConvertToBytes("1EB"))
	require.Equal(t, saturate(1000000000000000000), ConvertToBytes("1eb"))
	require.Equal(t, saturate(1000000000000000000), ConvertToBytes("1 EB"))
	require.Equal(t, saturate(1500000000000000000), ConvertToBytes("1.5EB"))
	require.Equal(t, saturate(9000000000000000000), ConvertToBytes("9EB"))
	require.Equal(t, math.MaxInt, ConvertToBytes("10EB")) // capped
	require.Equal(t, saturate(1<<60), ConvertToBytes("1EiB"))
	require.Equal(t, saturate(1<<60), ConvertToBytes("1Ei"))
	require.Equal(t, saturate(1500), ConvertToBytes("1.5e3B")) // an exponent, then bytes
	require.Equal(t, 1, ConvertToBytes("1E"))
	require.Equal(t, 1, ConvertToBytes("1e"))
	require.Equal(t, 5, ConvertToBytes("5 eggs"))

	require.Equal(t, 42*1024, ConvertToBytes("42Ki"))
	require.Equal(t, 42*1024, ConvertToBytes("42KiB"))
	require.Equal(t, 42*1024, ConvertToBytes("42kib"))
	require.Equal(t, 42*1024*1024, ConvertToBytes("42MiB"))
	require.Equal(t, saturate(42*1024*1024*1024), ConvertToBytes("42GiB"))
	require.Equal(t, saturate(42*1024*1024*1024*1024), ConvertToBytes("42TiB"))
	require.Equal(t, saturate(42*1024*1024*1024*1024*1024), ConvertToBytes("42PiB"))
	require.Equal(t, int(1.5*1024), ConvertToBytes("1.5KiB"))

	require.Equal(t, 0, ConvertToBytes("-5k"))
	require.Equal(t, 0, ConvertToBytes("-5.0k"))

	require.Equal(t, 0, ConvertToBytes("string"))
	require.Equal(t, 0, ConvertToBytes("MB"))
	require.Equal(t, 0, ConvertToBytes("invalidunit"))
	require.Equal(t, 42, ConvertToBytes("42X"))
	require.Equal(t, 123, ConvertToBytes("123a k"))
	require.Equal(t, 0, ConvertToBytes("42.5.5MB"))

	require.Equal(t, int(1.5*1000), ConvertToBytes("1.5k"))
	require.Equal(t, int(2.25*1000000), ConvertToBytes("2.25m"))
	require.Equal(t, int(0.5*1000000000), ConvertToBytes("0.5g"))

	require.Equal(t, 100*1000, ConvertToBytes("100 k"))
	require.Equal(t, 100*1000, ConvertToBytes("100  k"))
}

// "5 kilobytes" used to be read as 5 KiB: the 'i' that makes a unit binary comes right after the
// unit letter, and not at the start of a word.
func Test_ConvertToBytes_UnitWords(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]int64{
		"5 kilobytes":  5000,
		"5 Kilobytes":  5000,
		"5 gigabytes":  5000000000,
		"5 GIGABYTES":  5000000000,
		"5 megabytes":  5000000,
		"5 terabytes":  5000000000000,
		"5 petabytes":  5000000000000000,
		"5 kibibytes":  5120,
		"5 KIBIBYTES":  5120,
		"5 gibibytes":  5 << 30,
		"5 Gi":         5 << 30,
		"5Ki":          5120,
		"5KI":          5120,
		"5KiB":         5120,
		"5kib":         5120,
		"5 Ki B":       5120,
		"5Ki!":         5120,
		"5Ki.":         5120,
		"1.5 kilobyte": 1500,
	} {
		require.Equal(t, saturate(want), ConvertToBytes(in), in)
	}
}

// Fractional sizes that a float product truncates: "65.35KB" is 65350, not 65349.
func Test_ConvertToBytes_Fractions(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]int64{
		"65.35KB":    65350,
		"1.025M":     1025000,
		"130.51MB":   130510000,
		"132.593G":   132593000000,
		"1.001K":     1001,
		"1.15KB":     1150,
		"4.35K":      4350,
		"0.001K":     1,
		"0.0009K":    0, // 0.9 of a byte
		"0.5":        0,
		".5K":        500,
		"1.K":        1000,
		"2.675Ki":    2739, // 2739.2
		"0.1Mi":      104857,
		"1.5KiB":     1536,
		"1.5e3K":     1500000,
		"1.5E3K":     1500000,
		"007.50K":    7500,
		"1.250000K":  1250,
		"12.0000000": 12,
	} {
		require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in) // a 32-bit int saturates
	}
}

var (
	sizeUnits = map[string]int64{
		"": 1, "B": 1, "K": 1e3, "KB": 1e3, "M": 1e6, "MB": 1e6, "G": 1e9, "T": 1e12, "P": 1e15,
		"Ki": 1 << 10, "KiB": 1 << 10, "Mi": 1 << 20, "Gi": 1 << 30, "Ti": 1 << 40, "Pi": 1 << 50,
		"EB": 1e18, "Ei": 1 << 60, "EiB": 1 << 60,
	}
	sizeShape = regexp.MustCompile(`^\+?(\d+\.?\d*|\.\d+)(?:[eE]([+-]?\d+))?([A-Za-z]*)$`)
)

// modelSize is the oracle for ConvertToBytes: the exact size s denotes, from every digit,
// capped at the largest int64.
func modelSize(s string) (int64, bool) {
	m := sizeShape.FindStringSubmatch(s)
	if m == nil {
		return 0, false
	}
	mult, known := sizeUnits[m[3]]
	if !known {
		return 0, false
	}
	intPart, fracPart, _ := strings.Cut(m[1], ".")
	digits := strings.TrimLeft(intPart+fracPart, "0")
	if digits == "" {
		return 0, true
	}

	const maxExp = 1 << 40
	exp := new(big.Int)
	if m[2] != "" {
		exp.SetString(m[2], 10)
	}
	e := int64(maxExp)
	switch {
	case exp.CmpAbs(big.NewInt(maxExp)) <= 0:
		e = exp.Int64()
	case exp.Sign() < 0:
		e = -maxExp
	}
	e -= int64(len(fracPart))

	// The value is in [10^(n-1+e), 10^(n+e)); the unit scales it by less than 10^19.
	n := int64(len(digits))
	switch {
	case n-1+e >= 19:
		return math.MaxInt64, true
	case n+e+19 <= 0:
		return 0, true
	}
	value := new(big.Rat).SetInt(mustBigInt(digits))
	scale := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(max(e, -e)), nil))
	if e >= 0 {
		value.Mul(value, scale)
	} else {
		value.Quo(value, scale)
	}
	value.Mul(value, big.NewRat(mult, 1))
	q := new(big.Int).Quo(value.Num(), value.Denom())
	if !q.IsInt64() {
		return math.MaxInt64, true
	}
	return q.Int64(), true
}

func mustBigInt(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("not a number: " + s)
	}
	return n
}

// Test_ConvertToBytes_Exact compares random numbers of up to 60 digits with modelSize.
func Test_ConvertToBytes_Exact(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0xb17e5, 0x5eed)) //nolint:gosec // deterministic test data
	dropped := 0
	for range 4_000 {
		digits := make([]byte, 1+r.IntN(60))
		mix := r.IntN(3)
		for i := range digits {
			switch {
			case mix == 1 && r.IntN(4) != 0:
				digits[i] = '0'
			case mix == 2 && r.IntN(4) != 0:
				digits[i] = '9'
			default:
				digits[i] = byte('0' + r.IntN(10))
			}
		}
		num := string(digits)
		if r.IntN(4) != 0 {
			point := r.IntN(len(num) + 1)
			num = num[:point] + "." + num[point:]
		}
		if r.IntN(8) == 0 {
			num = "+" + num
		}
		if r.IntN(3) == 0 { // 10^-85 to 10^25 reaches every branch of the scaling
			exp := strconv.Itoa(r.IntN(111) - 85)
			if exp[0] != '-' && r.IntN(2) == 0 {
				exp = "+" + exp
			}
			num += string("eE"[r.IntN(2)]) + exp
		}
		for suffix := range sizeUnits {
			in := num + suffix
			want, ok := modelSize(in)
			require.True(t, ok, in)
			require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in)
			if want > 0 && want < math.MaxInt64 && len(strings.TrimLeft(string(digits), "0")) > 20 {
				dropped++
			}
		}
	}
	require.Greater(t, dropped, 4_000, "too few results that lose digits and still fit")
}

func Test_ConvertToBytes_Large(t *testing.T) {
	t.Parallel()
	// Past 2^53 a float64 cannot count bytes.
	for in, want := range map[string]int64{
		"9007199254740993K":     9007199254740993000,
		"9007199254740993B":     9007199254740993,
		"9007199254740993.000K": 9007199254740993000,
		"1000000.000000000001T": 1000000000000000001,
	} {
		require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in)
	}

	require.Equal(t, math.MaxInt, ConvertToBytes("9223372036854775807B"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9223372036854775808"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9999999999999999999K"))
	require.Equal(t, math.MaxInt, ConvertToBytes("18446744073709551615PiB"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9999999999999999.999P"))

	// More digits than a uint64 holds still count.
	require.Equal(t, math.MaxInt, ConvertToBytes("99999999999999999999"))
	require.Equal(t, math.MaxInt, ConvertToBytes("99999999999999999999K"))
	require.Equal(t, math.MaxInt, ConvertToBytes("99999999999999999999.5K"))
	require.Equal(t, 1000, ConvertToBytes("1.00000000000000000000001K"))
	require.Equal(t, 999, ConvertToBytes("0.99999999999999999999K"))
}

// '+' and the exponent scale exactly as the digits do; malformed input gives 0.
func Test_ConvertToBytes_SignAndExponent(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]int64{
		"+2.01k":                             2010,
		"+1.001K":                            1001,
		"+65.35KB":                           65350,
		".82e1GB":                            8200000000,
		"+.5K":                               500,
		"1e3K":                               1000000,
		"1E3":                                1000,
		"1e+3K":                              1000000,
		"5e-1K":                              500,
		"5.e1K":                              50000,
		"1.5e3K":                             1500000,
		"+5K":                                5000,
		"0e999K":                             0,
		"0.000e5K":                           0,
		"1e-999K":                            0,
		"1e-35K":                             0,
		"3e-17P":                             0,
		"3e-10P":                             300000,
		"12345678901234567890e-20P":          123456789012345,
		"12345678901234567890e-25P":          1234567890,
		"18446744073709551615e-34P":          1, // the widest product
		"18446744073709551615e-35P":          0, // past it
		"1e-38K":                             0, // the last place the two-step division covers
		"1e-39K":                             0, // the first it skips
		"1e-40K":                             0,
		"1e20":                               math.MaxInt64,
		"18446744073709551.7K":               math.MaxInt64, // the high word equals the divisor
		"18446744073709551.6K":               math.MaxInt64, // fits a uint64, not an int
		"99999999999999999999e-20K":          999,           // a dropped integer digit still counts
		"100000000000000000000e-20K":         1000,
		"1e18446744073709551616K":            math.MaxInt64, // an exponent that would wrap to 0
		"1e-18446744073709551616K":           0,
		"99999999999999999e5K":               math.MaxInt64,
		"184467440737095516.161e-5Mi":        1934281311383406679,
		"1844674407370955.1616199e-3Mi":      1934281311383406679,
		"184467440737095.51616199e-2Mi":      1934281311383406679,
		"0.000000000931322574615478515625Gi": 1, // exactly 1: 2^-30 Gi in 21 digits
		"3.5527136788005009293556213378906250e-15PiB": 4, // exactly 4: 2^-48 PiB in 34 digits
		"3813733478401949641952.6807962784e-15Ti":     4193244304741554190,
		"542202111511.42934007e-2Gi":                  5821850841909335364,
		"73634580.130256530332175T":                   math.MaxInt64,
		"1.0e400K":                                    math.MaxInt64,
		"1844674407370955.1615Ki":                     1888946593147858085,
		"9e18K":                                       math.MaxInt64,
		"1e999K":                                      math.MaxInt64,
		"1e19":                                        math.MaxInt64,
		"-5K":                                         0,
		"-5.0K":                                       0,
		"-0K":                                         0,
		"5+3K":                                        0,
		"++5K":                                        0,
		"+-5K":                                        0,
		"1e3.5K":                                      0,
		"1e3e3K":                                      0,
		"e3K":                                         0,
		"1.2.3K":                                      0,
		"1_000K":                                      0,
		"0x10K":                                       0,
		"1 000K":                                      0,
		"1e0000003":                                   1000,
	} {
		require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in)
	}

	// 2^63 and 2^63-0.5 both land on the largest int.
	for _, in := range []string{
		"+9223372036854775808.0",
		"9223372036854775808.0e0",
		"+9223372036854775807.5",
		"9223372036854775.808e3",
		"92233720368547758.08e2",
	} {
		require.Equal(t, math.MaxInt, ConvertToBytes(in), in)
	}
}

// Digits kept or dropped after 1844674407370955161, the first mantissa that takes no more digits.
func Test_ConvertToBytes_DroppedDigits(t *testing.T) {
	t.Parallel()

	for _, prefix := range []string{"1844674407370955159", "1844674407370955160", "1844674407370955161", "1844674407370955162"} {
		for _, d1 := range "069" {
			for _, d2 := range "05" {
				for _, tail := range []string{"", "99"} {
					digits := prefix + string(d1) + string(d2) + tail
					for point := 0; point <= len(digits); point++ {
						num := digits[:point] + "." + digits[point:]
						for _, exp := range []string{"", "e-3", "e-5", "e2"} {
							for _, unit := range []string{"", "Ki", "Mi"} {
								in := num + exp + unit
								want, ok := modelSize(in)
								require.True(t, ok, in)
								require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in)
							}
						}
					}
				}
			}
		}
	}
}

// Test_ConvertToBytes_Long pins digit counts and exponents beyond any fixed width.
func Test_ConvertToBytes_Long(t *testing.T) {
	t.Parallel()

	zeros := strings.Repeat("0", 1<<21) // more than 2^20, a limit the exponent once had
	for in, want := range map[string]int{
		"." + zeros[:1_048_557] + "1e10485760K": math.MaxInt, // the smallest input a limit of 2^20 gets wrong
		"0." + zeros + "1e2097152K":             100,
		"0." + zeros + "1e2097154K":             10000,
		"1" + zeros + "e-2097152K":              1000,
		"1" + zeros + "1e-2097152K":             10000,
		"." + zeros[:1<<20] + "1e1048576K":      100,
		"." + zeros + "1e99999999999K":          math.MaxInt,
		"1" + zeros + "e-99999999999K":          0,
		zeros + "1" + zeros + "K":               math.MaxInt,
		strings.Repeat("9", 1<<20) + "K":        math.MaxInt,
		"0." + strings.Repeat("9", 1<<20) + "K": 999,
	} {
		require.Equal(t, want, ConvertToBytes(in), "%.30s...", in)
	}

	// An 8-digit exponent, offset by 11 million digits.
	huge := strings.Repeat("0", 11_000_000)
	require.Equal(t, 100, ConvertToBytes("0."+huge+"1e11000000K"))
	require.Equal(t, 1000, ConvertToBytes("1"+huge+"e-11000000K"))
}

// Test_scaleDigits checks the digit-by-digit scaling on its own against big.Int.
func Test_scaleDigits(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewPCG(0x5ca1e, 0xd16)) //nolint:gosec // deterministic test data
	mults := []uint64{1, 1e3, 1e6, 1e9, 1e12, 1e15, 1 << 10, 1 << 20, 1 << 30, 1 << 40, 1 << 50}
	limit := big.NewInt(math.MaxInt64)
	saturated := 0
	for range 20_000 {
		digits := make([]byte, 20+r.IntN(26))
		for i := range digits {
			switch r.IntN(4) {
			case 0:
				digits[i] = '0'
			case 1:
				digits[i] = '9'
			default:
				digits[i] = byte('0' + r.IntN(10))
			}
		}
		digits[0] = byte('1' + r.IntN(9))
		text := string(digits)
		if r.IntN(2) == 0 {
			point := r.IntN(len(text) + 1)
			text = text[:point] + "." + text[point:]
		}
		if r.IntN(4) == 0 {
			text = "+" + text
		}
		mult := mults[r.IntN(len(mults))]
		down := int64(r.IntN(66)) - 5

		whole, _ := new(big.Int).SetString(strings.NewReplacer(".", "", "+", "").Replace(text), 10)
		want := whole.Mul(whole, new(big.Int).SetUint64(mult))
		if down > 0 {
			want.Quo(want, new(big.Int).Exp(big.NewInt(10), big.NewInt(down), nil))
		} else {
			want.Mul(want, new(big.Int).Exp(big.NewInt(10), big.NewInt(-down), nil))
		}
		if want.Cmp(limit) > 0 {
			want.Set(limit)
			saturated++
		}
		require.Equal(t, min(want.Int64(), math.MaxInt), int64(scaleDigits(text, len(text), down, mult)), "%s down %d mult %d", text, down, mult)
	}
	require.Greater(t, saturated, 1_000)

	// Edges of the cap: 19 places above the point, a uint64 past int64, the largest int64.
	for _, tc := range []struct {
		digits string
		down   int64
		mult   uint64
		want   int64
	}{
		{"20000000000000000000", 1, 10, math.MaxInt64},
		{"930000000000000000000", 2, 1, math.MaxInt64},
		{"92233720368547758070", 1, 1, math.MaxInt64},
		{"92233720368547758080", 1, 1, math.MaxInt64},
		{"92233720368547758060", 1, 1, 9223372036854775806},
		{"10000000000000000000", 0, 1, math.MaxInt64},
	} {
		require.Equal(t, min(tc.want, math.MaxInt), int64(scaleDigits(tc.digits, len(tc.digits), tc.down, tc.mult)), "%s down %d", tc.digits, tc.down)
	}
}

// The result is never negative and, where modelSize applies, equal to it.
func FuzzConvertToBytes(f *testing.F) {
	for _, s := range []string{
		"42", "42B", "42KB", "1.5MiB", "+2.01k", ".82e1GB", "1e3K", "9223372036854775808.0",
		"+9223372036854775807.5", "99999999999999999999K", "1.00000000000000000000001K", "-5K", "1e999K",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got := ConvertToBytes(s)
		require.GreaterOrEqual(t, got, 0, s)
		if want, ok := modelSize(s); ok {
			require.Equal(t, min(want, math.MaxInt), int64(got), s)
		}
	})
}

// parseByteSizeShape is ParseByteSize' grammar as a regexp: a number, blanks, and a unit if any. ASCII
// classes, not (?i), which would fold the Kelvin sign into k.
var parseByteSizeShape = regexp.MustCompile(`^\+?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?[ \t]*([bB]|[kKmMgGtTpP][iI]?[bB]?|[eE](?:[iI][bB]?|[bB]))?$`)

// requireNumError requires err to be the *strconv.NumError ParseByteSize gives for in.
func requireNumError(tb testing.TB, err, want error, in string) {
	tb.Helper()
	var ne *strconv.NumError
	require.ErrorAs(tb, err, &ne, in)
	require.Equal(tb, "ParseByteSize", ne.Func, in)
	require.Equal(tb, in, ne.Num, in)
	require.ErrorIs(tb, err, want, in)
}

func Test_ParseByteSize(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]int64{
		"0": 0, "42": 42, "+42": 42, "42b": 42, "42 B": 42, "42\tb": 42, "0.5": 0,
		"42k": 42000, "42KB": 42000, "42 kb": 42000, "42Ki": 42 << 10, "42KiB": 42 << 10, "42kib": 42 << 10,
		"1.5M": 1500000, ".5k": 500, "5.k": 5000, "1e3": 1000, "1.5e3K": 1500000, "1e3 B": 1000,
		"1MiB": 1 << 20, "1GB": 1e9, "1Gi": 1 << 30, "1TB": 1e12, "1PiB": 1 << 50,
		"1EB": 1e18, "1Eb": 1e18, "1Ei": 1 << 60, "1EiB": 1 << 60, "9EB": 9e18,
		strconv.Itoa(math.MaxInt - 1): int64(math.MaxInt - 1),
	} {
		got, err := ParseByteSize(in)
		if want >= math.MaxInt { // a 32-bit int cannot hold it
			require.Zero(t, got, in)
			requireNumError(t, err, strconv.ErrRange, in)
			continue
		}
		require.NoError(t, err, in)
		require.Equal(t, int(want), got, in)
	}

	// Where ConvertToBytes guesses, ParseByteSize fails: no number, a sign, text after the unit,
	// a unit that is not listed, a bare E, which is an exponent marker.
	for _, in := range []string{
		"", " ", "k", "KB", "-5", "-5k", "+-5", "5 X", "5kilobytes", "5 kb extra", "5k ", " 5k", "5kbb", "5kbi",
		"5ki b", "5 kB!", "5e", "5E", "5e+", "5e-", "5 e", "5eg", "5ex", "5bi", "5bb", "1.2.3", "1..2", ".", "+.",
		"5,000", "0x10", "5\n", "5k\n", "1_000", "\u0663kb", "5kb\x00", "5 kib b", strings.Repeat("9", 40) + "x",
		"5\u212aB", "5K\u0131B", "5\u017f", // letters that only fold to k, i or s
	} {
		got, err := ParseByteSize(in)
		require.Zero(t, got, in)
		requireNumError(t, err, strconv.ErrSyntax, in)
	}

	// MaxInt itself is out of range, as ConvertToBytes saturates there.
	for _, in := range []string{strconv.Itoa(math.MaxInt), "99999999999999999999999", "10EB", "9.3EB", "8EiB", "1e30"} {
		got, err := ParseByteSize(in)
		require.Zero(t, got, in)
		requireNumError(t, err, strconv.ErrRange, in)
	}
}

// Test_ParseByteSize_AgreesWithConvertToBytes: an accepted size scales as ConvertToBytes does.
func Test_ParseByteSize_AgreesWithConvertToBytes(t *testing.T) {
	t.Parallel()
	units := []string{
		"", "b", "k", "kb", "ki", "kib", "m", "mb", "mi", "mib", "g", "gb", "gi", "gib",
		"t", "tb", "ti", "tib", "p", "pb", "pi", "pib", "eb", "ei", "eib",
	}
	for _, num := range []string{"0", "1", "42", "+7", "1.5", ".5", "5.", "0.0009", "1e3", "1.5e3", "12345678901234567890", "99999999999999999999"} {
		for _, gap := range []string{"", " ", "\t", "  "} {
			for _, unit := range units {
				for _, in := range []string{num + gap + unit, num + gap + strings.ToUpper(unit)} {
					want := ConvertToBytes(in)
					got, err := ParseByteSize(in)
					if want == math.MaxInt {
						requireNumError(t, err, strconv.ErrRange, in)
						continue
					}
					require.NoError(t, err, in)
					require.Equal(t, want, got, in)
				}
			}
		}
	}
}

// Test_ParseByteSize_Grammar holds ParseByteSize to the regexp over random strings, edits and short units.
func Test_ParseByteSize_Grammar(t *testing.T) {
	t.Parallel()
	r := rand.New(rand.NewPCG(3, 9)) //nolint:gosec // deterministic test data
	const alphabet = "0123456789.+-eE \tbBkKiImMgGtTpPx"
	check := func(in string) {
		t.Helper()
		_, err := ParseByteSize(in)
		if syntax := errors.Is(err, strconv.ErrSyntax); syntax == parseByteSizeShape.MatchString(in) {
			t.Fatalf("ParseByteSize(%q) error = %v, but the grammar matches: %v", in, err, !syntax)
		}
	}
	for range 100_000 {
		b := make([]byte, r.IntN(11))
		for i := range b {
			b[i] = alphabet[r.IntN(len(alphabet))]
		}
		if len(b) > 0 && r.IntN(4) != 0 {
			b[0] = byte('0' + r.IntN(10))
		}
		check(string(b))
	}
	for _, valid := range []string{"42", "+4.2e1 KiB", ".5k", "5.EB", "1e3 b", "7 eib", "0.5Ti"} {
		for pos := range len(valid) {
			for _, c := range []byte(alphabet) {
				check(valid[:pos] + string(c) + valid[pos+1:])
			}
			check(valid[:pos])
			check(valid[:pos] + valid[pos+1:])
		}
	}
	for c := range 256 {
		check("5" + string([]byte{byte(c)}))
		check("5 " + string([]byte{byte(c)}))
		for d := range 256 {
			check("5" + string([]byte{byte(c), byte(d)}))
		}
	}
}

// FuzzParseByteSize: a size is rejected as text, out of range, or read as ConvertToBytes does.
func FuzzParseByteSize(f *testing.F) {
	for _, s := range []string{"42", "42KB", "1.5MiB", "+2.01k", "1e3K", "1EB", "9223372036854775807", "-5K", "5 X", "5e", ""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := ParseByteSize(s)
		switch {
		case err == nil:
			require.True(t, parseByteSizeShape.MatchString(s), s)
			require.Equal(t, ConvertToBytes(s), got, s)
			require.Less(t, got, math.MaxInt, s)
		case errors.Is(err, strconv.ErrRange):
			require.True(t, parseByteSizeShape.MatchString(s), s)
			require.Equal(t, math.MaxInt, ConvertToBytes(s), s)
			require.Zero(t, got, s)
		default:
			requireNumError(t, err, strconv.ErrSyntax, s)
			require.False(t, parseByteSizeShape.MatchString(s), s)
			require.Zero(t, got, s)
		}
	})
}

func Test_ConvertToBytes_DigitDetection(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected int
		desc     string
	}{
		{"0", 0, "digit 0"},
		{"1", 1, "digit 1"},
		{"9", 9, "digit 9"},
		{"123", 123, "multiple digits"},
		{"123k", 123000, "digits with unit"},
		{"a123", 0, "non-digit start"},
		{"12a3", 0, "non-digit in middle stops parsing"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.expected, ConvertToBytes(tc.input), "input: %s", tc.input)
		})
	}
}

func Test_GetArgument(t *testing.T) {
	originalArgs := os.Args

	defer func() { os.Args = originalArgs }()

	testArg := "test-arg"
	os.Args = []string{"cmd", testArg}

	require.True(t, GetArgument(testArg))
	require.False(t, GetArgument("missing-arg"))
}

func Test_GetArgument_Multiple(t *testing.T) {
	original := os.Args
	defer func() { os.Args = original }()

	os.Args = []string{"cmd", "-a", "-b", "--flag"}
	require.True(t, GetArgument("-a"))
	require.True(t, GetArgument("--flag"))
	require.False(t, GetArgument("-c"))
}

// os.Args[1:] panics when os.Args is empty, as in a process started without a program name.
func Test_GetArgument_NoArguments(t *testing.T) {
	original := os.Args
	defer func() { os.Args = original }()

	for _, args := range [][]string{nil, {}, {"cmd"}} {
		os.Args = args
		require.NotPanics(t, func() { require.False(t, GetArgument("-a")) }, "%q", args)
	}
	os.Args = []string{"-a"}
	require.False(t, GetArgument("-a"))
}

func Test_IncrementIPRange(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    net.IP
		expected net.IP
	}{
		{net.IP{192, 168, 1, 1}, net.IP{192, 168, 1, 2}},
		{net.IP{192, 168, 1, 254}, net.IP{192, 168, 1, 255}},
		{net.IP{192, 168, 1, 255}, net.IP{192, 168, 2, 0}},
		{net.IP{255, 255, 255, 255}, net.IP{0, 0, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.input.String(), func(t *testing.T) {
			t.Parallel()
			IncrementIPRange(c.input)
			require.Equal(t, c.expected, c.input)
		})
	}
}

func Test_IncrementIPRange_IPv6(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    net.IP
		expected net.IP
	}{
		{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{net.IP{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.input.String(), func(t *testing.T) {
			t.Parallel()
			IncrementIPRange(c.input)
			require.Equal(t, c.expected, c.input)
		})
	}
}

func Benchmark_ConvertToBytes(b *testing.B) {
	var res int
	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = ConvertToBytes("42B")
		}
		require.Equal(b, 42, res)
	})
}

func Benchmark_ConvertToBytes_Shapes(b *testing.B) {
	for _, in := range []string{
		"42", "42B", "42KB", "42 MB", "1.5MB", "130.51MB", "1.5GiB", "9007199254740993K", "1e3K",
		"1.2345678901234567890123K",
		"0.000000000931322574615478515625Gi",
	} {
		b.Run(in, func(b *testing.B) {
			b.ReportAllocs()
			var res int
			for b.Loop() {
				res = ConvertToBytes(in)
			}
			require.NotZero(b, res)
		})
	}
}

func Benchmark_GenerateSecureToken(b *testing.B) {
	var res string
	b.Run("16_bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = GenerateSecureToken(16)
		}
		require.Len(b, res, 22)
	})
	b.Run("32_bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = GenerateSecureToken(32)
		}
		require.Len(b, res, 43)
	})
}

func Benchmark_TokenGenerators(b *testing.B) {
	var res string
	b.Run("UUIDv4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = UUIDv4()
		}
		require.Len(b, res, 36)
	})
	b.Run("SecureToken", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = SecureToken()
		}
		require.Len(b, res, 43)
	})
}
