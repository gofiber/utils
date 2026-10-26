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

// Test_AddTrailingSlashBytes_ResultIsOwned holds the doc comment to its word that
// whenever a '/' is appended the caller gets a new slice. An empty input used to
// get a view of the constant "/": two results shared one array, and writing to one
// was a fatal fault, which is what the deprecated ToLowerBytes does to every slice.
func Test_AddTrailingSlashBytes_ResultIsOwned(t *testing.T) {
	t.Parallel()

	spare := make([]byte, 0, 8)
	inputs := map[string][]byte{
		"nil":                       nil,
		"empty":                     {},
		"empty with spare capacity": spare,
		"no slash":                  []byte("ab"),
	}
	for name, in := range inputs {
		want := append(append([]byte{}, in...), '/')
		first, second := AddTrailingSlashBytes(in), AddTrailingSlashBytes(in)
		require.Equal(t, want, first, name)
		require.NotSame(t, &first[0], &second[0], "%s: two results share storage", name)

		first[len(first)-1] = 'x' // a fault if the result is a view of read-only memory
		require.Equal(t, want, second, "%s: writing to one result changed another", name)
		require.Zero(t, spare[:1][0], "%s: the input's spare capacity was written", name)
	}
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
