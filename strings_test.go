// ⚡️ Fiber is an Express-inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io
//
// Package utils contains ASCII-optimized string conversion functions for Fiber.
// This test suite verifies ToUpper and ToLower for correctness and performance,
// covering HTTP methods, headers, URLs, and string patterns (lowercase, uppercase,
// camelCase, kebab-case, snake_case, mixed). It ensures no allocations for
// no-conversion cases and tests SWAR optimization thresholds (15–24 bytes).

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

	// HTTP methods (3–7 bytes)
	{name: "http-get", input: "get", upper: "GET", lower: "get", lowerNoConv: true},
	{name: "http-get-upper", input: "GET", upper: "GET", lower: "get", upperNoConv: true},
	{name: "http-get-pascal", input: "Get", upper: "GET", lower: "get"},
	{name: "http-head", input: "head", upper: "HEAD", lower: "head", lowerNoConv: true},
	{name: "http-head-upper", input: "HEAD", upper: "HEAD", lower: "head", upperNoConv: true},
	{name: "http-head-pascal", input: "Head", upper: "HEAD", lower: "head"},
	{name: "http-post", input: "post", upper: "POST", lower: "post", lowerNoConv: true},
	{name: "http-post-upper", input: "POST", upper: "POST", lower: "post", upperNoConv: true},
	{name: "http-post-pascal", input: "Post", upper: "POST", lower: "post"},
	{name: "http-put", input: "put", upper: "PUT", lower: "put", lowerNoConv: true},
	{name: "http-put-upper", input: "PUT", upper: "PUT", lower: "put", upperNoConv: true},
	{name: "http-put-pascal", input: "Put", upper: "PUT", lower: "put"},
	{name: "http-patch", input: "patch", upper: "PATCH", lower: "patch", lowerNoConv: true},
	{name: "http-patch-upper", input: "PATCH", upper: "PATCH", lower: "patch", upperNoConv: true},
	{name: "http-patch-pascal", input: "Patch", upper: "PATCH", lower: "patch"},
	{name: "http-delete", input: "delete", upper: "DELETE", lower: "delete", lowerNoConv: true},
	{name: "http-delete-upper", input: "DELETE", upper: "DELETE", lower: "delete", upperNoConv: true},
	{name: "http-delete-pascal", input: "Delete", upper: "DELETE", lower: "delete"},
	{name: "http-options", input: "options", upper: "OPTIONS", lower: "options", lowerNoConv: true},
	{name: "http-options-upper", input: "OPTIONS", upper: "OPTIONS", lower: "options", upperNoConv: true},
	{name: "http-options-pascal", input: "Options", upper: "OPTIONS", lower: "options"},

	// HTTP headers (10–50 bytes)
	{
		name: "header-content-type", input: "content-type: application/json; charset=utf-8",
		upper: "CONTENT-TYPE: APPLICATION/JSON; CHARSET=UTF-8",
		lower: "content-type: application/json; charset=utf-8", lowerNoConv: true,
	},
	{
		name: "header-content-type-upper", input: "CONTENT-TYPE: APPLICATION/JSON; CHARSET=UTF-8",
		upper: "CONTENT-TYPE: APPLICATION/JSON; CHARSET=UTF-8",
		lower: "content-type: application/json; charset=utf-8", upperNoConv: true,
	},
	{
		name: "header-content-type-mixed", input: "CoNtEnT-TyPe: application/json",
		upper: "CONTENT-TYPE: APPLICATION/JSON",
		lower: "content-type: application/json",
	},
	{
		name: "header-forwarded", input: "x-forwarded-for: 192.168.0.1",
		upper: "X-FORWARDED-FOR: 192.168.0.1",
		lower: "x-forwarded-for: 192.168.0.1", lowerNoConv: true,
	},
	{
		name: "header-forwarded-upper", input: "X-FORWARDED-FOR: 192.168.0.1",
		upper: "X-FORWARDED-FOR: 192.168.0.1",
		lower: "x-forwarded-for: 192.168.0.1", upperNoConv: true,
	},
	{
		name: "header-forwarded-mixed", input: "X-FoRwArDeD-FoR: 192.168.0.1",
		upper: "X-FORWARDED-FOR: 192.168.0.1",
		lower: "x-forwarded-for: 192.168.0.1",
	},
	{
		name: "header-upgrade", input: "Upgrade: WebSocket",
		upper: "UPGRADE: WEBSOCKET",
		lower: "upgrade: websocket",
	},

	// URLs (20–50 bytes)
	{
		name: "url", input: "/api/v1/users?name=John&sort=desc&limit=10",
		upper: "/API/V1/USERS?NAME=JOHN&SORT=DESC&LIMIT=10",
		lower: "/api/v1/users?name=john&sort=desc&limit=10",
	},
	{
		name: "url-upper", input: "/API/V1/USERS?NAME=JOHN&SORT=DESC&LIMIT=10",
		upper: "/API/V1/USERS?NAME=JOHN&SORT=DESC&LIMIT=10",
		lower: "/api/v1/users?name=john&sort=desc&limit=10", upperNoConv: true,
	},
	{
		name: "url-camel", input: "/api/v1/usersName?nameJohn&sortDesc",
		upper: "/API/V1/USERSNAME?NAMEJOHN&SORTDESC",
		lower: "/api/v1/usersname?namejohn&sortdesc",
	},

	// Length thresholds (15–25 bytes)
	{name: "15-char-lower", input: "abcdefghijklmno", upper: "ABCDEFGHIJKLMNO", lower: "abcdefghijklmno", lowerNoConv: true},
	{name: "15-char-mixed", input: "aBcDeFgHiJkLmNo", upper: "ABCDEFGHIJKLMNO", lower: "abcdefghijklmno"},
	{name: "15-char-upper", input: "ABCDEFGHIJKLMNO", upper: "ABCDEFGHIJKLMNO", lower: "abcdefghijklmno", upperNoConv: true},
	{name: "16-char-lower", input: "abcdefghijklmnop", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop", lowerNoConv: true},
	{name: "16-char-mixed", input: "aBcDeFgHiJkLmNoP", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},
	{name: "16-char-upper", input: "ABCDEFGHIJKLMNOP", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop", upperNoConv: true},
	{name: "17-char-lower", input: "abcdefghijklmnopq", upper: "ABCDEFGHIJKLMNOPQ", lower: "abcdefghijklmnopq", lowerNoConv: true},
	{name: "17-char-mixed", input: "aBcDeFgHiJkLmNoPq", upper: "ABCDEFGHIJKLMNOPQ", lower: "abcdefghijklmnopq"},
	{name: "17-char-upper", input: "ABCDEFGHIJKLMNOPQ", upper: "ABCDEFGHIJKLMNOPQ", lower: "abcdefghijklmnopq", upperNoConv: true},
	{name: "23-char-lower", input: "abcdefghijklmnopqrstuvw", upper: "ABCDEFGHIJKLMNOPQRSTUVW", lower: "abcdefghijklmnopqrstuvw", lowerNoConv: true},
	{name: "23-char-mixed", input: "aBcDeFgHiJkLmNoPqRsTuVw", upper: "ABCDEFGHIJKLMNOPQRSTUVW", lower: "abcdefghijklmnopqrstuvw"},
	{name: "23-char-upper", input: "ABCDEFGHIJKLMNOPQRSTUVW", upper: "ABCDEFGHIJKLMNOPQRSTUVW", lower: "abcdefghijklmnopqrstuvw", upperNoConv: true},
	{name: "24-char-lower", input: "abcdefghijklmnopqrstuvwx", upper: "ABCDEFGHIJKLMNOPQRSTUVWX", lower: "abcdefghijklmnopqrstuvwx", lowerNoConv: true},
	{name: "24-char-mixed", input: "aBcDeFgHiJkLmNoPqRsTuvWx", upper: "ABCDEFGHIJKLMNOPQRSTUVWX", lower: "abcdefghijklmnopqrstuvwx"},
	{name: "24-char-upper", input: "ABCDEFGHIJKLMNOPQRSTUVWX", upper: "ABCDEFGHIJKLMNOPQRSTUVWX", lower: "abcdefghijklmnopqrstuvwx", upperNoConv: true},
	{name: "25-char-lower", input: "abcdefghijklmnopqrstuvwxy", upper: "ABCDEFGHIJKLMNOPQRSTUVWXY", lower: "abcdefghijklmnopqrstuvwxy", lowerNoConv: true},
	{name: "25-char-mixed", input: "aBcDeFgHiJkLmNoPqRsTuVwXy", upper: "ABCDEFGHIJKLMNOPQRSTUVWXY", lower: "abcdefghijklmnopqrstuvwxy"},
	{name: "25-char-upper", input: "ABCDEFGHIJKLMNOPQRSTUVWXY", upper: "ABCDEFGHIJKLMNOPQRSTUVWXY", lower: "abcdefghijklmnopqrstuvwxy", upperNoConv: true},

	// Position-specific (16 bytes)
	{name: "16-char-first-lower", input: "aBCDEFGHIJKLMNOP", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},
	{name: "16-char-last-lower", input: "ABCDEFGHIJKLMNOp", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},
	{name: "16-char-first-upper", input: "Abcdefghijklmnop", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},
	{name: "16-char-last-upper", input: "abcdefghijklmnoP", upper: "ABCDEFGHIJKLMNOP", lower: "abcdefghijklmnop"},

	// Large strings (64 bytes)
	{name: "large-lower", input: strings.Repeat("a", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64), lowerNoConv: true},
	{name: "large-mixed", input: strings.Repeat("aB", 32), upper: strings.Repeat("AB", 32), lower: strings.Repeat("ab", 32)},
	{name: "large-upper", input: strings.Repeat("A", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64), upperNoConv: true},
	{name: "large-first-upper", input: "A" + strings.Repeat("a", 63), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-last-upper", input: strings.Repeat("a", 63) + "A", upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},

	// Very large strings (256 bytes)
	{name: "very-large-lower", input: strings.Repeat("a", 256), upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256), lowerNoConv: true},
	{name: "very-large-upper", input: strings.Repeat("A", 256), upper: strings.Repeat("A", 256), lower: strings.Repeat("a", 256), upperNoConv: true},
	{name: "very-large-mixed", input: strings.Repeat("aB", 128), upper: strings.Repeat("AB", 128), lower: strings.Repeat("ab", 128)},
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
				require.Zero(t, allocs, "ToUpper should not allocate for %s", tc.name)
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
				require.Zero(t, allocs, "ToLower should not allocate for %s", tc.name)
			}
		})
	}
}

func Test_ASCII_EdgeCases(t *testing.T) {
	t.Parallel()
	for i := 0; i < 128; i++ {
		c := byte(i)
		s := string(c)
		upperExpected := strings.ToUpper(s)
		lowerExpected := strings.ToLower(s)
		t.Run(fmt.Sprintf("ASCII-%d", i), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, upperExpected, ToUpper(s), "ToUpper failed for ASCII %d", i)
			require.Equal(t, lowerExpected, ToLower(s), "ToLower failed for ASCII %d", i)
		})
	}
}

func Test_NonASCII_Unchanged(t *testing.T) {
	t.Parallel()
	nonASCII := "µßäöü"
	t.Run("non-ascii", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, nonASCII, ToUpper(nonASCII), "ToUpper altered non-ASCII")
		require.Equal(t, nonASCII, ToLower(nonASCII), "ToLower altered non-ASCII")
	})
}

func Benchmark_ToUpper(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var allocsPerRun float64
			b.ReportAllocs()
			var res string
			allocsPerRun = testing.AllocsPerRun(100, func() {
				res = ToUpper(tc.input)
			})
			require.Equal(b, tc.upper, res)
			if tc.upperNoConv {
				require.Zerof(b, allocsPerRun, "ToUpper should not allocate for %s", tc.name)
			}
		})
	}
}

func Benchmark_ToLower(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var allocsPerRun float64
			b.ReportAllocs()
			var res string
			allocsPerRun = testing.AllocsPerRun(100, func() {
				res = ToLower(tc.input)
			})
			require.Equal(b, tc.lower, res)
			if tc.lowerNoConv {
				require.Zerof(b, allocsPerRun, "ToLower should not allocate for %s", tc.name)
			}
		})
	}
}

func Benchmark_StdToUpper(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var res string
			for n := 0; n < b.N; n++ {
				res = strings.ToUpper(tc.input)
			}
			require.Equal(b, tc.upper, res)
		})
	}
}

func Benchmark_StdToLower(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var res string
			for n := 0; n < b.N; n++ {
				res = strings.ToLower(tc.input)
			}
			require.Equal(b, tc.lower, res)
		})
	}
}
