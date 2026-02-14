// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			for i := 0; i < b.N; i++ {
				res = AddTrailingSlashBytes(tc.input)
			}
			_ = res
		})
	}
}
