// ⚡️ Fiber is an Express-inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io
//
// Package utils provides ASCII-optimized string case conversion functions for Fiber.
// The test suite verifies ToUpper and ToLower for correctness, performance, and allocation behavior.

package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// isRaceEnabled detects if the race detector is enabled
// This helps determine if allocation tests should be strict
func isRaceEnabled() bool {
	return raceEnabled
}

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

var benchmarkCases = []TestCase{
	// Benchmark cases for performance testing
	{name: "empty", input: "", upper: "", lower: "", upperNoConv: true, lowerNoConv: true},
	{name: "single-lower", input: "a", upper: "A", lower: "a", lowerNoConv: true},
	{name: "single-upper", input: "A", upper: "A", lower: "a", upperNoConv: true},
	{name: "numbers", input: "1234567890", upper: "1234567890", lower: "1234567890", upperNoConv: true, lowerNoConv: true},
	{name: "http-get", input: "get", upper: "GET", lower: "get", lowerNoConv: true},
	{name: "http-get-upper", input: "GET", upper: "GET", lower: "get", upperNoConv: true},
	{name: "http-get-pascal", input: "Get", upper: "GET", lower: "get"},
	{
		name: "header-content-type-mixed", input: "Content-Type: text/html; charset=utf-8",
		upper: "CONTENT-TYPE: TEXT/HTML; CHARSET=UTF-8",
		lower: "content-type: text/html; charset=utf-8",
	},
	{
		name: "url-camel", input: "/api/v1/usersName?nameJohn&sortDesc",
		upper: "/API/V1/USERSNAME?NAMEJOHN&SORTDESC",
		lower: "/api/v1/usersname?namejohn&sortdesc",
	},
	{name: "15-char-mixed", input: "aBcDeFgHiJkLmNo", upper: "ABCDEFGHIJKLMNO", lower: "abcdefghijklmno"},
	{name: "16-char-mixed", input: "aBcDeFgHiJkLmNoP", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},
	{name: "24-char-mixed", input: "aBcDeFgHiJkLmNoPqRsTuvWx", upper: "ABCDEFGHIJKLMNOPQRSTUVWX", lower: "abcdefghijklmnopqrstuvwx"},
	{name: "25-char-mixed", input: "aBcDeFgHiJkLmNoPqRsTuVwXy", upper: "ABCDEFGHIJKLMNOPQRSTUVWXY", lower: "abcdefghijklmnopqrstuvwxy"},
	{name: "large-lower", input: strings.Repeat("a", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64), lowerNoConv: true},
	{name: "large-mixed", input: strings.Repeat("aB", 32), upper: strings.Repeat("AB", 32), lower: strings.Repeat("ab", 32)},
	{name: "large-upper", input: strings.Repeat("A", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64), upperNoConv: true},
	{name: "large-first-upper", input: "A" + strings.Repeat("a", 63), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-last-upper", input: strings.Repeat("a", 63) + "A", upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "very-large-lower", input: strings.Repeat("a", 256), upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256), lowerNoConv: true},
	{name: "very-large-mixed", input: strings.Repeat("aB", 128), upper: strings.Repeat("AB", 128), lower: strings.Repeat("ab", 128)},
	{name: "very-large-upper", input: strings.Repeat("A", 256), upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256), upperNoConv: true},
	{name: "very-large-first-upper", input: "A" + strings.Repeat("a", 255), upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256)},
	{name: "very-large-last-upper", input: strings.Repeat("a", 255) + "A", upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256)},
}

func Test_ToUpper(t *testing.T) {
	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ToUpper(tc.input)
			require.Equal(t, tc.upper, result, "ToUpper failed for %s", tc.name)
			if tc.upperNoConv {
				allocs := testing.AllocsPerRun(100, func() {
					_ = ToUpper(tc.input)
				})
				// Skip strict allocation check when race detector is enabled
				// as it introduces additional allocations that are not part of the actual function
				if !isRaceEnabled() {
					require.Zero(t, allocs, "ToUpper should not allocate for %s", tc.name)
				} else {
					// In race mode, just log the allocation count for debugging
					t.Logf("ToUpper allocations for %s (with race detector): %f", tc.name, allocs)
				}
			}
		})
	}
}

func Test_ToLower(t *testing.T) {
	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ToLower(tc.input)
			require.Equal(t, tc.lower, result, "ToLower failed for %s", tc.name)
			if tc.lowerNoConv {
				allocs := testing.AllocsPerRun(100, func() {
					_ = ToLower(tc.input)
				})
				// Skip strict allocation check when race detector is enabled
				// as it introduces additional allocations that are not part of the actual function
				if !isRaceEnabled() {
					require.Zero(t, allocs, "ToLower should not allocate for %s", tc.name)
				} else {
					// In race mode, just log the allocation count for debugging
					t.Logf("ToLower allocations for %s (with race detector): %f", tc.name, allocs)
				}
			}
		})
	}
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
	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = ToUpper(tc.input)
			}
			require.Equal(b, tc.upper, res)
		})
	}
}

func Benchmark_ToLower(b *testing.B) {
	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = ToLower(tc.input)
			}
			require.Equal(b, tc.lower, res)
		})
	}
}

func Benchmark_StdToUpper(b *testing.B) {
	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = strings.ToUpper(tc.input)
			}
			require.Equal(b, tc.upper, res)
		})
	}
}

func Benchmark_StdToLower(b *testing.B) {
	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			var res string
			for n := 0; n < b.N; n++ {
				res = strings.ToLower(tc.input)
			}
			require.Equal(b, tc.lower, res)
		})
	}
}

// Test_ToUpper_NonASCII_Allocations provides diagnostic information about allocation behavior
// This test helps debug allocation issues in different environments (with/without race detector)
func Test_ToUpper_NonASCII_Allocations(t *testing.T) {
	t.Parallel()
	input := "µßäöü"
	allocs := testing.AllocsPerRun(100, func() {
		result := ToUpper(input)
		if result != input {
			t.Errorf("ToUpper modified input: got %v, want %v", result, input)
		}
	})
	t.Logf("Input bytes: %v", []byte(input))
	t.Logf("Allocations for ToUpper(µßäöü): %f", allocs)
	t.Logf("Race detector enabled: %v", isRaceEnabled())
	t.Logf("CI environment: %v", os.Getenv("CI"))
	t.Logf("GITHUB_ACTIONS: %v", os.Getenv("GITHUB_ACTIONS"))

	// Only enforce zero allocations when not running with race detector
	if !isRaceEnabled() {
		require.Zero(t, allocs, "ToUpper should not allocate for non-ascii without race detector")
	}
}
