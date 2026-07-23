package utils

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func jsonStringInputs() []string {
	return []string{
		"",
		"plain ascii value",
		"with \"quotes\" and \\backslashes\\",
		"tabs\tnewlines\nreturns\r",
		"controls \x00\x01\x1f\x7f",
		"html <b>bold</b> & more",
		"café naïve 日本語",
		"line\u2028and\u2029separators",
		"invalid \xff\xfe utf8 \xc3(",
		"trailing continuation \xc3",
		strings.Repeat("0123456789abcdef", 8),
		strings.Repeat("x", 7) + "\"" + strings.Repeat("y", 9),
	}
}

func Test_AppendJSONString(t *testing.T) {
	t.Parallel()
	for _, in := range jsonStringInputs() {
		want, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, string(want), string(AppendJSONString(nil, in)), "input %q", in)
		require.Equal(t, string(want), string(AppendJSONString(nil, []byte(in))), "input %q", in)
	}

	// Every single byte value agrees with encoding/json.
	for i := range 256 {
		in := string([]byte{byte(i)})
		want, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, string(want), string(AppendJSONString(nil, in)), "byte %#x", i)
	}

	// Appending must preserve existing dst content.
	got := AppendJSONString([]byte(`{"msg":`), "hi \"there\"")
	require.Equal(t, `{"msg":"hi \"there\""`, string(got)) //nolint:testifylint // byte-exact comparison is the contract; JSONEq would hide encoding differences
}

func Benchmark_AppendJSONString(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"clean-16B", "GET /index.html"},
		{"clean-64B", strings.Repeat("clean-log", 7) + "x"},
		{"escaped-64B", strings.Repeat(`says "hi" `, 6) + "done"},
		{"unicode-64B", strings.Repeat("café 日本 ", 6) + "end"},
	}
	for _, input := range inputs {
		dst := make([]byte, 0, 3*len(input.value))
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendJSONString(dst, input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := json.Marshal(input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
