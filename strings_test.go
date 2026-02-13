// ⚡️ Fiber is an Express-inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			for n := 0; n < b.N; n++ {
				res = AddTrailingSlashString(tc.input)
			}
			_ = res
		})
	}
}
