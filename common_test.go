// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"math"
	"math/big"
	"math/rand/v2"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FunctionName(t *testing.T) {
	t.Parallel()
	require.Equal(t, "github.com/gofiber/utils/v2.Test_UUIDv4", FunctionName(Test_UUIDv4))
	require.Equal(t, "github.com/gofiber/utils/v2.Test_FunctionName.func1", FunctionName(func() {}))

	dummyint := 20
	require.Equal(t, "int", FunctionName(dummyint))

	// nil interface should return empty string
	require.Empty(t, FunctionName(nil))

	// typed nil function should also return empty string
	var nilFunc func()
	require.Empty(t, FunctionName(nilFunc))

	// typed nil custom func should return empty string
	type myFunc func()
	var mf myFunc
	require.Empty(t, FunctionName(mf))

	// typed nil function inside interface should return empty string
	var iface any = mf
	require.Empty(t, FunctionName(iface))

	// typed nil pointer should return its type name
	var ptr *int
	require.Equal(t, "*int", FunctionName(ptr))

	// struct and pointer types should include package name
	type sampleStruct struct{}
	var s sampleStruct
	require.Equal(t, "utils.sampleStruct", FunctionName(s))
	require.Equal(t, "*utils.sampleStruct", FunctionName(&s))
}

func Test_UUIDv4(t *testing.T) {
	t.Parallel()
	res := UUIDv4()
	require.Len(t, res, 36)
	require.NotEqual(t, "00000000-0000-0000-0000-000000000000", res)
	require.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, res)
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
	// Test with 32 bytes
	token := GenerateSecureToken(32)
	require.Len(t, token, 43) // base64 encoding of 32 bytes
	require.NotEmpty(t, token)

	// Test custom length
	token8 := GenerateSecureToken(8)
	require.Len(t, token8, 11) // base64 of 8 bytes ~11 chars

	token16 := GenerateSecureToken(16)
	require.Len(t, token16, 22) // base64 of 16 bytes ~22 chars

	// Test uniqueness
	token2 := GenerateSecureToken(32)
	require.NotEqual(t, token, token2)

	// Test invalid length defaults to 32
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

	// Test uniqueness
	token2 := SecureToken()
	require.NotEqual(t, token, token2)
}

func Test_ConvertToBytes(t *testing.T) {
	t.Parallel()
	// initial assertions
	require.Equal(t, 0, ConvertToBytes(""))
	require.Equal(t, 42, ConvertToBytes("42"))

	// Test empty string
	require.Equal(t, 0, ConvertToBytes(""))

	// Test basic numbers (digit detection optimization)
	require.Equal(t, 42, ConvertToBytes("42"))
	require.Equal(t, 0, ConvertToBytes("0"))
	require.Equal(t, 1, ConvertToBytes("1"))
	require.Equal(t, 999, ConvertToBytes("999"))

	// Test with 'b' and 'B' suffixes
	require.Equal(t, 42, ConvertToBytes("42b"))
	require.Equal(t, 42, ConvertToBytes("42B"))
	require.Equal(t, 42, ConvertToBytes("42 b"))
	require.Equal(t, 42, ConvertToBytes("42 B"))

	// Test sizeMultipliers array usage (k/K - 1e3)
	require.Equal(t, 42*1000, ConvertToBytes("42k"))
	require.Equal(t, 42*1000, ConvertToBytes("42K"))
	require.Equal(t, 42*1000, ConvertToBytes("42kb"))
	require.Equal(t, 42*1000, ConvertToBytes("42KB"))
	require.Equal(t, 42*1000, ConvertToBytes("42 kb"))
	require.Equal(t, 42*1000, ConvertToBytes("42 KB"))

	// Test sizeMultipliers array usage (m/M - 1e6)
	require.Equal(t, 42*1000000, ConvertToBytes("42M"))
	require.Equal(t, 42*1000000, ConvertToBytes("42m"))
	require.Equal(t, 42*1000000, ConvertToBytes("42MB"))
	require.Equal(t, 42*1000000, ConvertToBytes("42mb"))
	require.Equal(t, int(42.5*1000000), ConvertToBytes("42.5MB"))

	// Test sizeMultipliers array usage (g/G - 1e9)
	require.Equal(t, 42*1000000000, ConvertToBytes("42G"))
	require.Equal(t, 42*1000000000, ConvertToBytes("42g"))
	require.Equal(t, 42*1000000000, ConvertToBytes("42GB"))
	require.Equal(t, 42*1000000000, ConvertToBytes("42gb"))

	// Test sizeMultipliers array usage (t/T - 1e12)
	require.Equal(t, 42*1000000000000, ConvertToBytes("42T"))
	require.Equal(t, 42*1000000000000, ConvertToBytes("42t"))
	require.Equal(t, 42*1000000000000, ConvertToBytes("42TB"))
	require.Equal(t, 42*1000000000000, ConvertToBytes("42tb"))

	// Test sizeMultipliers array usage (p/P - 1e15)
	require.Equal(t, 42*1000000000000000, ConvertToBytes("42P"))
	require.Equal(t, 42*1000000000000000, ConvertToBytes("42p"))
	require.Equal(t, 42*1000000000000000, ConvertToBytes("42PB"))
	require.Equal(t, 42*1000000000000000, ConvertToBytes("42pb"))

	// Test binary units (KiB = 1024)
	require.Equal(t, 42*1024, ConvertToBytes("42Ki"))
	require.Equal(t, 42*1024, ConvertToBytes("42KiB"))
	require.Equal(t, 42*1024, ConvertToBytes("42kib"))
	require.Equal(t, 42*1024*1024, ConvertToBytes("42MiB"))
	require.Equal(t, 42*1024*1024*1024, ConvertToBytes("42GiB"))
	require.Equal(t, 42*1024*1024*1024*1024, ConvertToBytes("42TiB"))
	require.Equal(t, 42*1024*1024*1024*1024*1024, ConvertToBytes("42PiB"))
	require.Equal(t, int(1.5*1024), ConvertToBytes("1.5KiB"))

	// Negative sizes are unrecognized
	require.Equal(t, 0, ConvertToBytes("-5k"))
	require.Equal(t, 0, ConvertToBytes("-5.0k"))

	// Test edge cases and error conditions
	require.Equal(t, 0, ConvertToBytes("string"))
	require.Equal(t, 0, ConvertToBytes("MB"))
	require.Equal(t, 0, ConvertToBytes("invalidunit"))
	require.Equal(t, 42, ConvertToBytes("42X"))     // invalid unit
	require.Equal(t, 123, ConvertToBytes("123a k")) // invalid unit
	require.Equal(t, 0, ConvertToBytes("42.5.5MB")) // invalid format

	// Test decimal numbers with various units
	require.Equal(t, int(1.5*1000), ConvertToBytes("1.5k"))
	require.Equal(t, int(2.25*1000000), ConvertToBytes("2.25m"))
	require.Equal(t, int(0.5*1000000000), ConvertToBytes("0.5g"))

	// Test space handling
	require.Equal(t, 100*1000, ConvertToBytes("100 k"))
	require.Equal(t, 100*1000, ConvertToBytes("100  k")) // multiple spaces
}

// Test_ConvertToBytes_Fractions pins fractional sizes that a float product
// truncates ("65.35KB" is 65350, not 65349).
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
		"0.5":        0, // half a byte
		".5K":        500,
		"1.K":        1000, // trailing point
		"2.675Ki":    2739, // 2739.2
		"0.1Mi":      104857,
		"1.5KiB":     1536,
		"1.5e3K":     1500000, // exponent: float path
		"1.5E3K":     1500000,
		"007.50K":    7500,
		"1.250000K":  1250,
		"12.0000000": 12,
	} {
		require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in) // a 32-bit int saturates
	}
}

// Test_ConvertToBytes_Exact checks up to nine digits with a point anywhere,
// before every unit, against big.Rat.
func Test_ConvertToBytes_Exact(t *testing.T) {
	t.Parallel()

	units := map[string]int64{
		"": 1, "B": 1, "K": 1e3, "KB": 1e3, "M": 1e6, "MB": 1e6, "G": 1e9, "T": 1e12, "P": 1e15,
		"Ki": 1 << 10, "KiB": 1 << 10, "Mi": 1 << 20, "Gi": 1 << 30, "Ti": 1 << 40, "Pi": 1 << 50,
	}
	r := rand.New(rand.NewPCG(0xb17e5, 0x5eed)) //nolint:gosec // deterministic test data
	maxInt := new(big.Int).SetInt64(math.MaxInt)
	for range 10_000 {
		digits := make([]byte, 1+r.IntN(9))
		for i := range digits {
			digits[i] = byte('0' + r.IntN(10))
		}
		num := string(digits)
		if r.IntN(4) != 0 {
			point := r.IntN(len(num) + 1)
			num = num[:point] + "." + num[point:]
		}
		for suffix, mult := range units {
			in := num + suffix
			exact, ok := new(big.Rat).SetString(num)
			require.True(t, ok, num)
			exact.Mul(exact, big.NewRat(mult, 1))
			want := new(big.Int).Quo(exact.Num(), exact.Denom())
			if want.Cmp(maxInt) > 0 {
				want = maxInt
			}
			require.Equal(t, want.Int64(), int64(ConvertToBytes(in)), in)
		}
	}
}

func Test_ConvertToBytes_Large(t *testing.T) {
	t.Parallel()
	// Past 2^53 a float64 cannot count bytes; min() is for 32-bit ints.
	for in, want := range map[string]int64{
		"9007199254740993K":     9007199254740993000,
		"9007199254740993B":     9007199254740993,
		"9007199254740993.000K": 9007199254740993000,
		"1000000.000000000001T": 1000000000000000001, // 10^18 + 1
	} {
		require.Equal(t, min(want, math.MaxInt), int64(ConvertToBytes(in)), in)
	}

	// Saturation at the largest int, however it is reached.
	require.Equal(t, math.MaxInt, ConvertToBytes("9223372036854775807B"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9223372036854775808"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9999999999999999999K"))
	require.Equal(t, math.MaxInt, ConvertToBytes("18446744073709551615PiB"))
	require.Equal(t, math.MaxInt, ConvertToBytes("9999999999999999.999P"))

	// More digits than a uint64 holds: plain digits clamp, a unit makes them
	// invalid, and with a point the float path decides.
	require.Equal(t, math.MaxInt, ConvertToBytes("99999999999999999999"))
	require.Equal(t, 0, ConvertToBytes("99999999999999999999K"))
	require.Equal(t, math.MaxInt, ConvertToBytes("99999999999999999999.5K"))
	require.Equal(t, 1000, ConvertToBytes("1.00000000000000000000001K")) // 23 decimals
}

// Test_ConvertToBytes_NeedsPoint pins that a sign or an exponent is a size only
// next to a point; without one it is not recognized.
func Test_ConvertToBytes_NeedsPoint(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]int{
		"1e3K":   0,
		"+5K":    0,
		"1.5e3K": 1500000,
		"+5.0K":  5000,
	} {
		require.Equal(t, want, ConvertToBytes(in), in)
	}
}

func Test_ConvertToBytes_DigitDetection(t *testing.T) {
	t.Parallel()
	// Test the new direct byte comparison digit detection
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

	// Reset os.Args
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

// Additional coverage for IncrementIPRange using IPv6 addresses
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

// go test -v -run=^$ -bench=Benchmark_ConvertToBytes -benchmem -count=2
func Benchmark_ConvertToBytes(b *testing.B) {
	var res int
	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = ConvertToBytes("42B")
		}
		require.Equal(b, 42, res)
	})
}

// go test -v -run=^$ -bench=Benchmark_ConvertToBytes_Shapes -benchmem -count=2
func Benchmark_ConvertToBytes_Shapes(b *testing.B) {
	for _, in := range []string{"42", "42B", "42KB", "42 MB", "1.5MB", "130.51MB", "1.5GiB", "9007199254740993K"} {
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

// go test -v -run=^$ -bench=Benchmark_GenerateSecureToken -benchmem -count=2
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

// go test -v -run=^$ -bench=Benchmark_TokenGenerators -benchmem -count=2
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
