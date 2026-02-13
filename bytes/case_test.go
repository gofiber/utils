package bytes

import (
	stdbytes "bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type benchCase struct {
	name  string
	input string
	lower string
	upper string
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
	{name: "large-lower", input: strings.Repeat("a", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-upper", input: strings.Repeat("A", 64), upper: strings.Repeat("A", 64), lower: strings.Repeat("a", 64)},
	{name: "large-mixed", input: strings.Repeat("aB", 32), upper: strings.Repeat("AB", 32), lower: strings.Repeat("ab", 32)},
}

func Test_ToLowerBytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, []byte("/my/name/is/:param/*"), UnsafeToLower([]byte("/MY/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my1/name/is/:param/*"), UnsafeToLower([]byte("/MY1/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my2/name/is/:param/*"), UnsafeToLower([]byte("/MY2/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my3/name/is/:param/*"), UnsafeToLower([]byte("/MY3/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my4/name/is/:param/*"), UnsafeToLower([]byte("/MY4/NAME/IS/:PARAM/*")))
}

func Test_ToUpperBytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, []byte("/MY/NAME/IS/:PARAM/*"), UnsafeToUpper([]byte("/my/name/is/:param/*")))
	require.Equal(t, []byte("/MY1/NAME/IS/:PARAM/*"), UnsafeToUpper([]byte("/my1/name/is/:param/*")))
	require.Equal(t, []byte("/MY2/NAME/IS/:PARAM/*"), UnsafeToUpper([]byte("/my2/name/is/:param/*")))
	require.Equal(t, []byte("/MY3/NAME/IS/:PARAM/*"), UnsafeToUpper([]byte("/my3/name/is/:param/*")))
	require.Equal(t, []byte("/MY4/NAME/IS/:PARAM/*"), UnsafeToUpper([]byte("/my4/name/is/:param/*")))
}

func Test_ToLower_NoMutation(t *testing.T) {
	in := []byte("Content-Type")
	snapshot := append([]byte(nil), in...)
	out := ToLower(in)

	require.Equal(t, "content-type", string(out))
	require.Equal(t, snapshot, in)
}

func Test_ToUpper_NoMutation(t *testing.T) {
	in := []byte("content-type")
	snapshot := append([]byte(nil), in...)
	out := ToUpper(in)

	require.Equal(t, "CONTENT-TYPE", string(out))
	require.Equal(t, snapshot, in)
}

func TestUnsafeToLower_Mutates(t *testing.T) {
	in := []byte("Content-Type")
	out := UnsafeToLower(in)
	require.Same(t, &out[0], &in[0])
	require.Equal(t, "content-type", string(in))
}

func TestUnsafeToUpper_Mutates(t *testing.T) {
	in := []byte("content-type")
	out := UnsafeToUpper(in)
	require.Same(t, &out[0], &in[0])
	require.Equal(t, "CONTENT-TYPE", string(in))
}

func Test_ToLowerBytes_Edge(t *testing.T) {
	t.Parallel()

	cases := [][]byte{
		{},
		[]byte("123"),
		[]byte("!@#"),
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, stdbytes.ToLower(c), UnsafeToLower(c))
		})
	}
}

func Test_ToUpperBytes_Edge(t *testing.T) {
	t.Parallel()

	cases := [][]byte{
		{},
		[]byte("123"),
		[]byte("!@#"),
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, stdbytes.ToUpper(c), UnsafeToUpper(c))
		})
	}
}

func Benchmark_ToLowerBytes(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			template := []byte(tc.input)
			want := []byte(tc.lower)
			var res []byte

			b.Run("fiber", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					res = ToLower(template)
				}
				require.Equal(b, want, res)
			})
			b.Run("fiber/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(template))
				for i := 0; i < b.N; i++ {
					copy(work, template)
					res = UnsafeToLower(work)
				}
				require.Equal(b, want, res)
			})
			b.Run("default", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					res = stdbytes.ToLower(template)
				}
				require.Equal(b, want, res)
			})
		})
	}
}

func Benchmark_ToUpperBytes(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			template := []byte(tc.input)
			want := []byte(tc.upper)
			var res []byte

			b.Run("fiber", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					res = ToUpper(template)
				}
				require.Equal(b, want, res)
			})
			b.Run("fiber/unsafe", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(template))
				for i := 0; i < b.N; i++ {
					copy(work, template)
					res = UnsafeToUpper(work)
				}
				require.Equal(b, want, res)
			})
			b.Run("default", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					res = stdbytes.ToUpper(template)
				}
				require.Equal(b, want, res)
			})
		})
	}
}
