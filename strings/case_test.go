package strings

import (
	"testing"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	stdstrings "strings"
)

type benchCase struct {
	name  string
	input string
	lower string
	upper string
}

var benchCases = []benchCase{
	{name: "empty", input: "", lower: "", upper: ""},
	{name: "http-get", input: "get", lower: "get", upper: "GET"},
	{name: "http-get-upper", input: "GET", lower: "get", upper: "GET"},
	{
		name:  "header-content-type-mixed",
		input: "Content-Type: text/html; charset=utf-8",
		lower: "content-type: text/html; charset=utf-8",
		upper: "CONTENT-TYPE: TEXT/HTML; CHARSET=UTF-8",
	},
}

func TestToLower_NoMutation(t *testing.T) {
	in := "Content-Type"
	out := ToLower(in)
	if out != "content-type" {
		t.Fatalf("unexpected output: %q", out)
	}
	if in != "Content-Type" {
		t.Fatalf("input changed: %q", in)
	}
}

func TestToUpper_NoMutation(t *testing.T) {
	in := "content-type"
	out := ToUpper(in)
	if out != "CONTENT-TYPE" {
		t.Fatalf("unexpected output: %q", out)
	}
	if in != "content-type" {
		t.Fatalf("input changed: %q", in)
	}
}

func TestUnsafeToLower_Mutates(t *testing.T) {
	buf := []byte("Content-Type")
	s := unsafeconv.UnsafeString(buf)
	out := UnsafeToLower(s)
	if out != "content-type" {
		t.Fatalf("unexpected output: %q", out)
	}
	if string(buf) != "content-type" {
		t.Fatalf("backing bytes not mutated: %q", buf)
	}
}

func TestUnsafeToUpper_Mutates(t *testing.T) {
	buf := []byte("content-type")
	s := unsafeconv.UnsafeString(buf)
	out := UnsafeToUpper(s)
	if out != "CONTENT-TYPE" {
		t.Fatalf("unexpected output: %q", out)
	}
	if string(buf) != "CONTENT-TYPE" {
		t.Fatalf("backing bytes not mutated: %q", buf)
	}
}

func Benchmark_SubpkgToLower(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			var res string

			b.Run("subpkg", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = ToLower(tc.input)
				}
				if res != tc.lower {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("subpkg/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				template := []byte(tc.input)
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					res = UnsafeToLower(unsafeconv.UnsafeString(work))
				}
				if res != tc.lower {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("stdlib", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = stdstrings.ToLower(tc.input)
				}
				if res != tc.lower {
					b.Fatalf("unexpected output: %q", res)
				}
			})
		})
	}
}

func Benchmark_SubpkgToUpper(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			var res string

			b.Run("subpkg", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = ToUpper(tc.input)
				}
				if res != tc.upper {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("subpkg/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				template := []byte(tc.input)
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					res = UnsafeToUpper(unsafeconv.UnsafeString(work))
				}
				if res != tc.upper {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("stdlib", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = stdstrings.ToUpper(tc.input)
				}
				if res != tc.upper {
					b.Fatalf("unexpected output: %q", res)
				}
			})
		})
	}
}
