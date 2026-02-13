package strings

import (
	"fmt"
	stdstrings "strings"
	"testing"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/stretchr/testify/require"
)

type benchCase struct {
	name  string
	input string
	lower string
	upper string
}

var testCases = []benchCase{
	{name: "empty", input: "", upper: "", lower: ""},
	{name: "numbers", input: "1234567890", upper: "1234567890", lower: "1234567890"},
	{name: "symbols", input: "!@#$%^&*()", upper: "!@#$%^&*()", lower: "!@#$%^&*()"},
	{name: "single-lower", input: "a", upper: "A", lower: "a"},
	{name: "single-upper", input: "A", upper: "A", lower: "a"},
	{name: "lowercase", input: "abcdefghijklmnopqrstuvwxyz", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz"},
	{name: "uppercase", input: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz"},
	{name: "mixed-case", input: "aBcDeFgHiJkLmNoPqRsTuVwXyZ", upper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", lower: "abcdefghijklmnopqrstuvwxyz"},
	{name: "non-ascii", input: "µßäöü", upper: "µßäöü", lower: "µßäöü"},
	{name: "mixed-ascii-non-ascii", input: "Goµ", upper: "GOµ", lower: "goµ"},
}

var benchmarkCoreCases = []benchCase{
	{name: "empty", input: "", upper: "", lower: ""},
	{name: "http-get", input: "get", upper: "GET", lower: "get"},
	{name: "http-get-upper", input: "GET", upper: "GET", lower: "get"},
	{
		name:  "header-content-type-mixed",
		input: "Content-Type: text/html; charset=utf-8",
		upper: "CONTENT-TYPE: TEXT/HTML; CHARSET=UTF-8",
		lower: "content-type: text/html; charset=utf-8",
	},
	{name: "large-lower", input: stdstrings.Repeat("a", 64), upper: stdstrings.Repeat("A", 64), lower: stdstrings.Repeat("a", 64)},
	{name: "large-upper", input: stdstrings.Repeat("A", 64), upper: stdstrings.Repeat("A", 64), lower: stdstrings.Repeat("a", 64)},
	{name: "large-mixed", input: stdstrings.Repeat("aB", 32), upper: stdstrings.Repeat("AB", 32), lower: stdstrings.Repeat("ab", 32)},
}

func Test_ToUpper(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.upper, ToUpper(tc.input))
		})
	}
}

func Test_ToLower(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.lower, ToLower(tc.input))
		})
	}
}

func Test_ASCII_EdgeCases(t *testing.T) {
	t.Parallel()
	for i := 0; i < 128; i++ {
		idx := i
		c := byte(idx)
		t.Run(fmt.Sprintf("ASCII-char-%03d", idx), func(t *testing.T) {
			t.Parallel()
			s := string(c)
			require.Equal(t, stdstrings.ToUpper(s), ToUpper(s))
			require.Equal(t, stdstrings.ToLower(s), ToLower(s))
		})
	}
}

func Test_ToLower_NoMutation(t *testing.T) {
	in := "Content-Type"
	out := ToLower(in)
	require.Equal(t, "content-type", out)
	require.Equal(t, "Content-Type", in)
}

func Test_ToUpper_NoMutation(t *testing.T) {
	in := "content-type"
	out := ToUpper(in)
	require.Equal(t, "CONTENT-TYPE", out)
	require.Equal(t, "content-type", in)
}

func TestUnsafeToLower_Mutates(t *testing.T) {
	buf := []byte("Content-Type")
	s := unsafeconv.UnsafeString(buf)
	out := UnsafeToLower(s)
	require.Equal(t, "content-type", out)
	require.Equal(t, "content-type", string(buf))
}

func TestUnsafeToUpper_Mutates(t *testing.T) {
	buf := []byte("content-type")
	s := unsafeconv.UnsafeString(buf)
	out := UnsafeToUpper(s)
	require.Equal(t, "CONTENT-TYPE", out)
	require.Equal(t, "CONTENT-TYPE", string(buf))
}

func Benchmark_ToUpper(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tc.input)))
			var fiberRes, stdRes string
			b.Run("fiber", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					fiberRes = ToUpper(tc.input)
				}
				require.Equal(b, tc.upper, fiberRes)
			})
			b.Run("fiber/unsafe", func(b *testing.B) {
				template := []byte(tc.input)
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					fiberRes = UnsafeToUpper(unsafeconv.UnsafeString(work))
				}
				require.Equal(b, tc.upper, fiberRes)
			})
			b.Run("default", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					stdRes = stdstrings.ToUpper(tc.input)
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
			var fiberRes, stdRes string
			b.Run("fiber", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					fiberRes = ToLower(tc.input)
				}
				require.Equal(b, tc.lower, fiberRes)
			})
			b.Run("fiber/unsafe", func(b *testing.B) {
				template := []byte(tc.input)
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					fiberRes = UnsafeToLower(unsafeconv.UnsafeString(work))
				}
				require.Equal(b, tc.lower, fiberRes)
			})
			b.Run("default", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					stdRes = stdstrings.ToLower(tc.input)
				}
				require.Equal(b, tc.lower, stdRes)
			})
		})
	}
}
