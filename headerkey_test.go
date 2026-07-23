package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func headerKeySamples() []string {
	return []string{
		"",
		"Content-Type",
		"content-type",
		"CONTENT-TYPE",
		"cOnTeNt-TyPe",
		"Host",
		"host",
		"x-request-id",
		"X-Request-ID",
		"x--double-dash",
		"-leading-dash",
		"trailing-dash-",
		"1-2-3",
		"etag",
		"WWW-Authenticate",
		"www-authenticate",
		"Header With Space",
		"Bad:Colon",
		"Tab\tKey",
		"non-ascii-é",
		"under_score",
		"tilde~ok",
	}
}

func Test_CanonicalHeaderKey(t *testing.T) {
	t.Parallel()
	for _, in := range headerKeySamples() {
		want := http.CanonicalHeaderKey(in)
		require.Equal(t, want, CanonicalHeaderKey(in), "input %q", in)
		require.Equal(t, want, string(CanonicalHeaderKey([]byte(in))), "input %q", in)
	}

	// Already-canonical and invalid keys come back without copying: for
	// []byte inputs the same backing array is returned.
	for _, in := range []string{"Content-Type", "Bad Key", "X-Id"} {
		b := []byte(in)
		out := CanonicalHeaderKey(b)
		require.Same(t, &b[0], &out[0], "input %q must not copy", in)
	}
}

func Benchmark_CanonicalHeaderKey(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"canonical", "Content-Type"},
		{"common-lower", "content-type"},
		{"custom-lower", "x-custom-request-id"},
	}
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				CanonicalHeaderKey(input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				http.CanonicalHeaderKey(input.value)
			}
		})
	}
}
