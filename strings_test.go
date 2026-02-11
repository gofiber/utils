// ⚡️ Fiber is an Express-inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io
//
// Package utils provides ASCII-optimized string case conversion functions for Fiber.
// The test suite verifies ToUpper and ToLower for correctness, performance, and allocation behavior.

package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCase defines a test case for case conversion.
type TestCase struct {
	name        string
	input       string
	upper       string // Expected output for ToUpper
	lower       string // Expected output for ToLower
	upperNoConv bool   // True if ToUpper should return input without allocation
	lowerNoConv bool   // True if ToLower should return input without allocation
}

var testCases = []TestCase{
	// Edge cases
	{name: "empty", input: "", upper: "", lower: "", upperNoConv: true, lowerNoConv: true},
	{name: "numbers", input: "1234567890", upper: "1234567890", lower: "1234567890", upperNoConv: true, lowerNoConv: true},
	{name: "symbols", input: "!@#$%^&*()", upper: "!@#$%^&*()", lower: "!@#$%^&*()", upperNoConv: true, lowerNoConv: true},
	{name: "single-lower", input: "a", upper: "A", lower: "a", lowerNoConv: true},
	{name: "single-upper", input: "A", upper: "A", lower: "a", upperNoConv: true},

	// ASCII letter conversion tests
	{name: "lowercase", input: "abcdefghijklmnopqrstuvwxyz", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz", lowerNoConv: true},
	{name: "uppercase", input: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz", upperNoConv: true},
	{name: "mixed-case", input: "aBcDeFgHiJkLmNoPqRsTuVwXyZ", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz"},

	// Non-ASCII characters should not be converted
	{name: "non-ascii", input: "µßäöü", upper: "µßäöü", lower: "µßäöü", upperNoConv: true, lowerNoConv: true},
	{name: "mixed-ascii-non-ascii", input: "Goµ", upper: "GOµ", lower: "goµ"},
}

var benchmarkCoreCases = []TestCase{
	{name: "empty", input: "", upper: "", lower: ""},
	{name: "http-get", input: "get", upper: "GET", lower: "get"},
	{name: "http-get-upper", input: "GET", upper: "GET", lower: "get"},
	{
		name: "header-content-type-mixed", input: "Content-Type: text/html; charset=utf-8",
		upper: "CONTENT-TYPE: TEXT/HTML; CHARSET=UTF-8",
		lower: "content-type: text/html; charset=utf-8",
	},
	{name: "large-lower", input: strings.Repeat("a", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-upper", input: strings.Repeat("A", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-mixed", input: strings.Repeat("aB", 32), upper: strings.Repeat("AB", 32), lower: strings.Repeat("ab", 32)},
}

func Test_ToUpper(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToUpper(tc.input)
			require.Equal(t, tc.upper, result, "ToUpper failed for %s", tc.name)
			if tc.upperNoConv {
				// Warm up the function to handle any first-run initialization costs
				for i := 0; i < 10; i++ {
					_ = ToUpper(tc.input)
				}

				allocs := testing.AllocsPerRun(100, func() {
					_ = ToUpper(tc.input)
				})

				if !raceEnabled {
					require.Zero(t, allocs, "ToUpper should not allocate for %s", tc.name)
				} else {
					// In race mode, just log the allocation count to avoid false failures
					t.Logf("ToUpper allocations for %s (with race detector): %f", tc.name, allocs)
				}
			}
		})
	}
}

func Test_ToLower(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToLower(tc.input)
			require.Equal(t, tc.lower, result, "ToLower failed for %s", tc.name)
			if tc.lowerNoConv {
				// Warm up the function to handle any first-run initialization costs
				for i := 0; i < 10; i++ {
					_ = ToLower(tc.input)
				}

				allocs := testing.AllocsPerRun(100, func() {
					_ = ToLower(tc.input)
				})

				if !raceEnabled {
					require.Zero(t, allocs, "ToLower should not allocate for %s", tc.name)
				} else {
					// In race mode, just log the allocation count to avoid false failures
					t.Logf("ToLower allocations for %s (with race detector): %f", tc.name, allocs)
				}
			}
		})
	}
}

func Test_ToLowerMut(t *testing.T) {
	t.Parallel()

	buf := []byte("Content-Type")
	s := UnsafeString(buf)
	out := ToLowerMut(s)

	require.Equal(t, "content-type", out)
	require.Equal(t, []byte("content-type"), buf, "backing bytes should be mutated in-place")
}

func Test_ToUpperMut(t *testing.T) {
	t.Parallel()

	buf := []byte("content-type")
	s := UnsafeString(buf)
	out := ToUpperMut(s)

	require.Equal(t, "CONTENT-TYPE", out)
	require.Equal(t, []byte("CONTENT-TYPE"), buf, "backing bytes should be mutated in-place")
}

func Test_ASCII_EdgeCases(t *testing.T) {
	// Test ASCII characters from 0 to 127 for ToUpper and ToLower
	// This ensures that all basic ASCII characters are handled correctly.
	t.Parallel()
	for i := 0; i < 128; i++ {
		idx := i
		c := byte(idx)
		t.Run(fmt.Sprintf("ASCII-char-%03d", idx), func(t *testing.T) {
			t.Parallel()
			s := string(c)
			upperExpected := strings.ToUpper(s)
			lowerExpected := strings.ToLower(s)
			require.Equal(t, upperExpected, ToUpper(s), "ToUpper failed for ASCII-char-%03d", idx)
			require.Equal(t, lowerExpected, ToLower(s), "ToLower failed for ASCII-char-%03d", idx)
		})
	}
}

func Benchmark_ToUpper(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var fiberRes, stdRes string
			b.Run("fiber", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					fiberRes = ToUpper(tc.input)
				}
				require.Equal(b, tc.upper, fiberRes)
			})
			b.Run("fiber/mut", func(b *testing.B) {
				template := []byte(tc.input)
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					fiberRes = ToUpperMut(UnsafeString(work))
				}
				require.Equal(b, tc.upper, fiberRes)
			})
			b.Run("default", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					stdRes = strings.ToUpper(tc.input)
				}
				require.Equal(b, tc.upper, stdRes)
			})
		})
	}
}

func Benchmark_ToLower(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var fiberRes, stdRes string
			b.Run("fiber", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					fiberRes = ToLower(tc.input)
				}
				require.Equal(b, tc.lower, fiberRes)
			})
			b.Run("fiber/mut", func(b *testing.B) {
				template := []byte(tc.input)
				work := make([]byte, len(template))
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					copy(work, template)
					fiberRes = ToLowerMut(UnsafeString(work))
				}
				require.Equal(b, tc.lower, fiberRes)
			})
			b.Run("default", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					stdRes = strings.ToLower(tc.input)
				}
				require.Equal(b, tc.lower, stdRes)
			})
		})
	}
}

func Test_AddTrailingSlashString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "/"},
		{name: "slash-only", in: "/", want: "/"},
		{name: "short-no-slash", in: "abc", want: "abc/"},
		{name: "short-with-slash", in: "abc/", want: "abc/"},
		{name: "path-no-slash", in: "/api/v1/users", want: "/api/v1/users/"},
		{name: "path-with-slash", in: "/api/v1/users/", want: "/api/v1/users/"},
		{name: "double-slash", in: "abc//", want: "abc//"},
		{name: "unicode", in: "/日本語", want: "/日本語/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := AddTrailingSlashString(tc.in)
			require.Equal(t, tc.want, result)
		})
	}
}

func Benchmark_AddTrailingSlashString(b *testing.B) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "slash-only", input: "/"},
		{name: "short-no-slash", input: "abc"},
		{name: "short-with-slash", input: "abc/"},
		{name: "path-no-slash", input: "/api/v1/users"},
		{name: "path-with-slash", input: "/api/v1/users/"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var res string
			for b.Loop() {
				res = AddTrailingSlashString(tc.input)
			}
			_ = res
		})
	}
}
