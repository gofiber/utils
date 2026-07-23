package utils

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseIPSamples() []string {
	return []string{
		// IPv4
		"0.0.0.0",
		"127.0.0.1",
		"203.0.113.195",
		"255.255.255.255",
		"1.2.3.4",
		"256.1.1.1",
		"1.2.3.4.5",
		"1.2.3",
		"01.2.3.4",
		"1.02.3.4",
		"1.2.3.04",
		"1.2.3.4 ",
		" 1.2.3.4",
		"1..3.4",
		"1.2.3.",
		".1.2.3",
		"1.2.3.4567",
		"a.b.c.d",
		"1.2.3.-4",
		"",
		// IPv6
		"::",
		"::1",
		"1::",
		"2001:db8::8a2e:370:7334",
		"2001:0db8:0000:0000:0000:ff00:0042:8329",
		"1:2:3:4:5:6:7:8",
		"1:2:3:4:5:6:7:8:9",
		"1:2:3:4:5:6:7::",
		"1::2:3:4:5:6:7:8",
		"1:::2",
		"1::2::3",
		":1:2:3",
		"1:2:3:",
		":",
		":::",
		"fe80::1%eth0",
		"fe80::1%",
		"%eth0",
		"::%25eth0",
		"::ffff:192.168.0.1",
		"::ffff:1.2.3.400",
		"1:2:3:4:5:6:1.2.3.4",
		"1:2:3:4:5:6:7:1.2.3.4",
		"::1.2.3.4",
		"1.2.3.4::",
		"12345::",
		"g::1",
		"FE80::A:b",
		"0:0:0:0:0:0:0:0",
		"1:2::10.0.0.1",
		"::.1.2.3",
	}
}

func Test_ParseIPv4(t *testing.T) {
	t.Parallel()
	for _, in := range parseIPSamples() {
		want, wantErr := netip.ParseAddr(in)
		wantOK := wantErr == nil && !strings.Contains(in, ":")
		got, ok := ParseIPv4(in)
		gotBytes, okBytes := ParseIPv4([]byte(in))
		require.Equal(t, wantOK, ok, "input %q", in)
		require.Equal(t, wantOK, okBytes, "input %q", in)
		if wantOK {
			require.Equal(t, want, got, "input %q", in)
			require.Equal(t, want, gotBytes, "input %q", in)
		} else {
			require.Equal(t, netip.Addr{}, got, "input %q", in)
		}
	}
}

func Test_ParseIPv6(t *testing.T) {
	t.Parallel()
	for _, in := range parseIPSamples() {
		want, wantErr := netip.ParseAddr(in)
		wantOK := wantErr == nil && strings.Contains(in, ":")
		got, ok := ParseIPv6(in)
		gotBytes, okBytes := ParseIPv6([]byte(in))
		require.Equal(t, wantOK, ok, "input %q", in)
		require.Equal(t, wantOK, okBytes, "input %q", in)
		if wantOK {
			require.Equal(t, want, got, "input %q", in)
			require.Equal(t, want, gotBytes, "input %q", in)
		} else {
			require.Equal(t, netip.Addr{}, got, "input %q", in)
		}
	}
}

func Benchmark_ParseIPv4(b *testing.B) {
	input := "203.0.113.195"
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := ParseIPv4(input); !ok {
				b.Fatal("failed to parse")
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := netip.ParseAddr(input); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Benchmark_ParseIPv6(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"compressed", "2001:db8::8a2e:370:7334"},
		{"full", "2001:0db8:0000:0000:0000:ff00:0042:8329"},
	}
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, ok := ParseIPv6(input.value); !ok {
					b.Fatal("failed to parse")
				}
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := netip.ParseAddr(input.value); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
