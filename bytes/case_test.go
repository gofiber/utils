package bytes

import (
	stdbytes "bytes"
	"testing"
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
	in := []byte("Content-Type")
	snapshot := append([]byte(nil), in...)
	out := ToLower(in)

	if string(out) != "content-type" {
		t.Fatalf("unexpected output: %q", out)
	}
	if !stdbytes.Equal(in, snapshot) {
		t.Fatalf("input was mutated: %q", in)
	}
}

func TestToUpper_NoMutation(t *testing.T) {
	in := []byte("content-type")
	snapshot := append([]byte(nil), in...)
	out := ToUpper(in)

	if string(out) != "CONTENT-TYPE" {
		t.Fatalf("unexpected output: %q", out)
	}
	if !stdbytes.Equal(in, snapshot) {
		t.Fatalf("input was mutated: %q", in)
	}
}

func TestUnsafeToLower_Mutates(t *testing.T) {
	in := []byte("Content-Type")
	out := UnsafeToLower(in)
	if &out[0] != &in[0] {
		t.Fatal("expected same backing array")
	}
	if string(in) != "content-type" {
		t.Fatalf("unexpected output: %q", in)
	}
}

func TestUnsafeToUpper_Mutates(t *testing.T) {
	in := []byte("content-type")
	out := UnsafeToUpper(in)
	if &out[0] != &in[0] {
		t.Fatal("expected same backing array")
	}
	if string(in) != "CONTENT-TYPE" {
		t.Fatalf("unexpected output: %q", in)
	}
}

func Benchmark_SubpkgToLower(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			input := []byte(tc.input)
			want := []byte(tc.lower)
			var res []byte

			b.Run("subpkg", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = ToLower(input)
				}
				if !stdbytes.Equal(want, res) {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("subpkg/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(input))
				for n := 0; n < b.N; n++ {
					copy(work, input)
					res = UnsafeToLower(work)
				}
				if !stdbytes.Equal(want, res) {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("stdlib", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = stdbytes.ToLower(input)
				}
				if !stdbytes.Equal(want, res) {
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
			input := []byte(tc.input)
			want := []byte(tc.upper)
			var res []byte

			b.Run("subpkg", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = ToUpper(input)
				}
				if !stdbytes.Equal(want, res) {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("subpkg/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(input))
				for n := 0; n < b.N; n++ {
					copy(work, input)
					res = UnsafeToUpper(work)
				}
				if !stdbytes.Equal(want, res) {
					b.Fatalf("unexpected output: %q", res)
				}
			})

			b.Run("stdlib", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = stdbytes.ToUpper(input)
				}
				if !stdbytes.Equal(want, res) {
					b.Fatalf("unexpected output: %q", res)
				}
			})
		})
	}
}
