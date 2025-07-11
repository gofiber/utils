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

func Benchmark_ToLowerBytes(b *testing.B) {
	input := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:',.<>/?`~")
	expected := bytes.ToLower(input)
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLowerBytes(input)
		}
		require.True(b, bytes.Equal(expected, res))
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.ToLower(input)
		}
		require.True(b, bytes.Equal(expected, res))
	})
}

func Benchmark_ToUpperBytes(b *testing.B) {
	input := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:',.<>/?`~")
	expected := bytes.ToUpper(input)
	var res []byte

	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToUpperBytes(input)
		}
		require.Equal(b, expected, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = bytes.ToUpper(input)
		}
		require.Equal(b, expected, res)
	})
}
