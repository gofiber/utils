// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	largeStr = "/RePos/GoFiBer/FibEr/iSsues/187643/CoMmEnts/RePos/GoFiBer/FibEr/iSsues/CoMmEnts"
	upperStr = "/REPOS/GOFIBER/FIBER/ISSUES/187643/COMMENTS/REPOS/GOFIBER/FIBER/ISSUES/COMMENTS"
	lowerStr = "/repos/gofiber/fiber/issues/187643/comments/repos/gofiber/fiber/issues/comments"
)

func Test_ToLowerBytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, []byte("/my/name/is/:param/*"), ToLowerBytes([]byte("/MY/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my1/name/is/:param/*"), ToLowerBytes([]byte("/MY1/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my2/name/is/:param/*"), ToLowerBytes([]byte("/MY2/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my3/name/is/:param/*"), ToLowerBytes([]byte("/MY3/NAME/IS/:PARAM/*")))
	require.Equal(t, []byte("/my4/name/is/:param/*"), ToLowerBytes([]byte("/MY4/NAME/IS/:PARAM/*")))
}

func Test_ToUpperBytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, []byte("/MY/NAME/IS/:PARAM/*"), ToUpperBytes([]byte("/my/name/is/:param/*")))
	require.Equal(t, []byte("/MY1/NAME/IS/:PARAM/*"), ToUpperBytes([]byte("/my1/name/is/:param/*")))
	require.Equal(t, []byte("/MY2/NAME/IS/:PARAM/*"), ToUpperBytes([]byte("/my2/name/is/:param/*")))
	require.Equal(t, []byte("/MY3/NAME/IS/:PARAM/*"), ToUpperBytes([]byte("/my3/name/is/:param/*")))
	require.Equal(t, []byte("/MY4/NAME/IS/:PARAM/*"), ToUpperBytes([]byte("/my4/name/is/:param/*")))
}

func Test_ToLowerBytes_NoMutation(t *testing.T) {
	t.Parallel()

	original := []byte("/MY/NAME/IS/:PARAM/*")
	snapshot := append([]byte(nil), original...)
	out := ToLowerBytes(original)

	require.Equal(t, snapshot, original, "input slice should not be mutated")
	require.Equal(t, []byte("/my/name/is/:param/*"), out)
}

func Test_ToUpperBytes_NoMutation(t *testing.T) {
	t.Parallel()

	original := []byte("/my/name/is/:param/*")
	snapshot := append([]byte(nil), original...)
	out := ToUpperBytes(original)

	require.Equal(t, snapshot, original, "input slice should not be mutated")
	require.Equal(t, []byte("/MY/NAME/IS/:PARAM/*"), out)
}

func Test_ToLowerBytesMut(t *testing.T) {
	t.Parallel()

	buf := []byte("/MY/NAME/IS/:PARAM/*")
	out := ToLowerBytesMut(buf)
	require.Same(t, &buf[0], &out[0], "should return same backing array")
	require.Equal(t, []byte("/my/name/is/:param/*"), buf)
}

func Test_ToUpperBytesMut(t *testing.T) {
	t.Parallel()

	buf := []byte("/my/name/is/:param/*")
	out := ToUpperBytesMut(buf)
	require.Same(t, &buf[0], &out[0], "should return same backing array")
	require.Equal(t, []byte("/MY/NAME/IS/:PARAM/*"), buf)
}

func Benchmark_ToLowerBytes(b *testing.B) {
	for _, tc := range benchmarkCoreCases {
		b.Run(tc.name, func(b *testing.B) {
			template := []byte(tc.input)
			want := []byte(tc.lower)
			var res []byte

			b.Run("fiber", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = ToLowerBytes(template)
				}
				require.True(b, bytes.Equal(want, res))
			})
			b.Run("fiber/mut", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					res = ToLowerBytesMut(work)
				}
				require.True(b, bytes.Equal(want, res))
			})
			b.Run("default", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = bytes.ToLower(template)
				}
				require.True(b, bytes.Equal(want, res))
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
				for n := 0; n < b.N; n++ {
					res = ToUpperBytes(template)
				}
				require.Equal(b, want, res)
			})
			b.Run("fiber/mut", func(b *testing.B) {
				b.ReportAllocs()
				work := make([]byte, len(template))
				for n := 0; n < b.N; n++ {
					copy(work, template)
					res = ToUpperBytesMut(work)
				}
				require.Equal(b, want, res)
			})
			b.Run("default", func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					res = bytes.ToUpper(template)
				}
				require.Equal(b, want, res)
			})
		})
	}
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
			require.Equal(t, bytes.ToLower(c), ToLowerBytes(c))
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
			require.Equal(t, bytes.ToUpper(c), ToUpperBytes(c))
		})
	}
}

func Test_AddTrailingSlashBytes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   []byte
		want []byte
	}{
		{name: "empty", in: []byte(""), want: []byte("/")},
		{name: "slash-only", in: []byte("/"), want: []byte("/")},
		{name: "short-no-slash", in: []byte("abc"), want: []byte("abc/")},
		{name: "short-with-slash", in: []byte("abc/"), want: []byte("abc/")},
		{name: "path-no-slash", in: []byte("/api/v1/users"), want: []byte("/api/v1/users/")},
		{name: "path-with-slash", in: []byte("/api/v1/users/"), want: []byte("/api/v1/users/")},
		{name: "double-slash", in: []byte("abc//"), want: []byte("abc//")},
		{name: "unicode", in: []byte("/日本語"), want: []byte("/日本語/")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := AddTrailingSlashBytes(tc.in)
			require.Equal(t, tc.want, result)
		})
	}
}

func Test_AddTrailingSlashBytes_NoMutation(t *testing.T) {
	t.Parallel()

	original := []byte("test")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	_ = AddTrailingSlashBytes(original)

	require.Equal(t, originalCopy, original, "original slice should not be mutated")
}

func Test_AddTrailingSlashBytes_ReturnsSame(t *testing.T) {
	t.Parallel()

	input := []byte("test/")
	result := AddTrailingSlashBytes(input)
	require.Equal(t, input, result)
	require.Same(t, &input[0], &result[0], "should return same slice instance")
}

func Benchmark_AddTrailingSlashBytes(b *testing.B) {
	cases := []struct {
		name  string
		input []byte
	}{
		{name: "empty", input: []byte("")},
		{name: "slash-only", input: []byte("/")},
		{name: "short-no-slash", input: []byte("abc")},
		{name: "short-with-slash", input: []byte("abc/")},
		{name: "path-no-slash", input: []byte("/api/v1/users")},
		{name: "path-with-slash", input: []byte("/api/v1/users/")},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var res []byte
			for n := 0; n < b.N; n++ {
				res = AddTrailingSlashBytes(tc.input)
			}
			_ = res
		})
	}
}
