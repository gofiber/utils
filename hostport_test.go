package utils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func hostPortSamples() []string {
	return []string{
		"",
		":",
		"::",
		":80",
		"example.com",
		"example.com:",
		"example.com:80",
		"example.com:8080:90",
		"host%eth0:1",
		"[::1]",
		"[::1]:",
		"[::1]:8080",
		"[fe80::1%eth0]:443",
		"[::1]x:80",
		"[::1]:80:90",
		"[::1:80",
		"::1]:80",
		"[]:80",
		"[[::1]]:80",
		"[::1]]:80",
		"a[b:80",
		"a]b:80",
		"[a]:[b]",
		"[a]:b]",
		"1.2.3.4:65535",
		" example.com:80",
		"example.com :80",
		"\x00:\x00",
	}
}

// checkSplitHostPort compares both input types with net.SplitHostPort.
func checkSplitHostPort(tb testing.TB, in string) {
	tb.Helper()
	wantHost, wantPort, err := net.SplitHostPort(in)
	host, port, ok := SplitHostPort(in)
	require.Equal(tb, err == nil, ok, "input %q: %v", in, err)
	if !ok {
		require.Empty(tb, host, "input %q", in)
		require.Empty(tb, port, "input %q", in)
	} else {
		require.Equal(tb, wantHost, host, "input %q", in)
		require.Equal(tb, wantPort, port, "input %q", in)
	}

	bh, bp, bok := SplitHostPort([]byte(in))
	require.Equal(tb, ok, bok, "input %q bytes", in)
	require.Equal(tb, host, string(bh), "input %q bytes", in)
	require.Equal(tb, port, string(bp), "input %q bytes", in)
	if !bok {
		require.Nil(tb, bh, "input %q bytes", in)
		require.Nil(tb, bp, "input %q bytes", in)
	}
}

func Test_SplitHostPort(t *testing.T) {
	t.Parallel()
	for _, in := range hostPortSamples() {
		checkSplitHostPort(t, in)
	}

	// The parts alias the input.
	in := []byte("[::1]:8080")
	host, port, ok := SplitHostPort(in)
	require.True(t, ok)
	require.Same(t, &in[1], &host[0])
	require.Same(t, &in[6], &port[0])
}

func FuzzSplitHostPort(f *testing.F) {
	for _, s := range hostPortSamples() {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		checkSplitHostPort(t, s)
	})
}

func Benchmark_SplitHostPort(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"host-port", "example.com:8080"},
		{"ipv6-port", "[2001:db8::1]:8080"},
		{"missing-port", "example.com"},
	}
	var host, port string
	var ok bool
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				host, port, ok = SplitHostPort(input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var err error
				host, port, err = net.SplitHostPort(input.value)
				ok = err == nil
			}
		})
	}
	_ = host + port
	_ = ok
}
