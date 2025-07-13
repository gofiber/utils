// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FunctionName(t *testing.T) {
	t.Parallel()
	require.Equal(t, "github.com/gofiber/utils/v2.Test_UUID", FunctionName(Test_UUID))
	require.Equal(t, "github.com/gofiber/utils/v2.Test_FunctionName.func1", FunctionName(func() {}))

	dummyint := 20
	require.Equal(t, "int", FunctionName(dummyint))

	// nil interface should return empty string
	require.Equal(t, "", FunctionName(nil))

	// typed nil function should also return empty string
	var nilFunc func()
	require.Equal(t, "", FunctionName(nilFunc))

	// typed nil custom func should return empty string
	type myFunc func()
	var mf myFunc
	require.Equal(t, "", FunctionName(mf))

	// typed nil function inside interface should return empty string
	var iface any = mf
	require.Equal(t, "", FunctionName(iface))

	// typed nil pointer should return its type name
	var ptr *int
	require.Equal(t, "*int", FunctionName(ptr))

	// struct and pointer types should include package name
	type sampleStruct struct{}
	var s sampleStruct
	require.Equal(t, "utils.sampleStruct", FunctionName(s))
	require.Equal(t, "*utils.sampleStruct", FunctionName(&s))
}

func Test_UUID(t *testing.T) {
	t.Parallel()
	res := UUID()
	require.Len(t, res, 36)
	require.NotEqual(t, "00000000-0000-0000-0000-000000000000", res)
}

func Test_UUID_Concurrency(t *testing.T) {
	t.Parallel()
	iterations := 1000
	var res string
	ch := make(chan string, iterations)
	results := make(map[string]string)
	for i := 0; i < iterations; i++ {
		go func() {
			ch <- UUID()
		}()
	}
	for i := 0; i < iterations; i++ {
		res = <-ch
		results[res] = res
	}
	require.Len(t, results, iterations)
}

func Test_UUIDv4(t *testing.T) {
	t.Parallel()
	res := UUIDv4()
	require.Len(t, res, 36)
	require.NotEqual(t, "00000000-0000-0000-0000-000000000000", res)
}

func Test_UUIDv4_Concurrency(t *testing.T) {
	t.Parallel()
	iterations := 1000
	var res string
	ch := make(chan string, iterations)
	results := make(map[string]string)
	for i := 0; i < iterations; i++ {
		go func() {
			ch <- UUIDv4()
		}()
	}
	for i := 0; i < iterations; i++ {
		res = <-ch
		results[res] = res
	}
	require.Len(t, results, iterations)
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
		c := c
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
		c := c
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
		for n := 0; n < b.N; n++ {
			res = ConvertToBytes("42B")
		}
		require.Equal(b, 42, res)
	})
}

// go test -v -run=^$ -bench=Benchmark_UUID -benchmem -count=2
func Benchmark_UUID(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = UUID()
		}
		require.Len(b, res, 36)
	})
	b.Run("default", func(b *testing.B) {
		rnd := make([]byte, 16)
		_, _ = rand.Read(rnd) //nolint: errcheck // No need to check error
		for n := 0; n < b.N; n++ {
			res = fmt.Sprintf("%x-%x-%x-%x-%x", rnd[0:4], rnd[4:6], rnd[6:8], rnd[8:10], rnd[10:])
		}
		require.Len(b, res, 36)
	})
}
