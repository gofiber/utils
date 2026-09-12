// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// assertIPGrammarParity pins the three invariants that tie this module's two IP
// grammars together: IsIPv4 and IsIPv6 accept exactly what net.ParseIP accepts
// for their family, and each agrees with the netip-based parser beside it once
// the zone suffix net.ParseIP rejects is accounted for. Test_IsIP_Parity and
// FuzzIsIP share it, so the two grammars cannot drift apart unnoticed again.
func assertIPGrammarParity(t *testing.T, s string) {
	t.Helper()

	hasColon := strings.ContainsRune(s, ':')
	netOK := net.ParseIP(s) != nil

	is4 := IsIPv4(s)
	is6 := IsIPv6(s)
	if want := netOK && !hasColon; is4 != want {
		t.Fatalf("IsIPv4(%q) = %v, want %v (net.ParseIP)", s, is4, want)
	}
	if want := netOK && hasColon; is6 != want {
		t.Fatalf("IsIPv6(%q) = %v, want %v (net.ParseIP)", s, is6, want)
	}

	// net.ParseIP is netip.ParseAddr minus zones, so the parsers agree with the
	// validators on every input that carries no zone.
	if _, ok := ParseIPv4(s); ok != is4 {
		t.Fatalf("ParseIPv4(%q) ok = %v but IsIPv4 = %v", s, ok, is4)
	}
	addr, ok6 := ParseIPv6(s)
	if want := ok6 && addr.Zone() == ""; want != is6 {
		t.Fatalf("ParseIPv6(%q) ok = %v zone = %q but IsIPv6 = %v", s, ok6, addr.Zone(), is6)
	}
}

// Test_IsIP_Parity runs the shared IP corpus through both grammars.
func Test_IsIP_Parity(t *testing.T) {
	t.Parallel()

	for _, in := range parseIPSamples() {
		assertIPGrammarParity(t, in)
	}
}

// Test_IsIPv6_HexFieldDigitCap covers fields longer than four hex digits whose
// value still fits in 0xFFFF. A value-only bound accepted these; net.ParseIP,
// the reference IsIPv6 follows, rejects them.
func Test_IsIPv6_HexFieldDigitCap(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		"00001::2",
		"cec::7dcE:0D568",
		"::00000",
		"00000::",
		"1:2:3:4:5:6:7:00008",
		"00000000000000001::",
	} {
		require.False(t, IsIPv6(in), "over-long hex field: %q", in)
		require.Nil(t, net.ParseIP(in), "net.ParseIP disagrees: %q", in)
	}

	// Exactly four digits stays valid, leading zeros included.
	for _, in := range []string{
		"0001:2:3:4:5:6:7:8",
		"0000::0000",
		"::0001",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
	} {
		require.True(t, IsIPv6(in), "four-digit field rejected: %q", in)
		require.NotNil(t, net.ParseIP(in), "net.ParseIP disagrees: %q", in)
	}
}

func Test_IsIPv4(t *testing.T) {
	t.Parallel()

	require.True(t, IsIPv4("255.255.255.255"))
	require.True(t, IsIPv4("174.23.33.100"))
	require.True(t, IsIPv4("127.0.0.1"))
	require.True(t, IsIPv4("0.0.0.0"))
	require.True(t, IsIPv4("1.1.1.1"))

	require.False(t, IsIPv4(".0.0.0"))
	require.False(t, IsIPv4("0.0.0."))
	require.False(t, IsIPv4("0.0.0"))
	require.False(t, IsIPv4(".0.0.0."))
	require.False(t, IsIPv4("0.0.0.0.0"))
	require.False(t, IsIPv4("0"))
	require.False(t, IsIPv4(""))
	require.False(t, IsIPv4("2345:0425:2CA1::0567:5673:23b5"))
	require.False(t, IsIPv4("invalid"))
	require.False(t, IsIPv4("189.12.34.260"))
	require.False(t, IsIPv4("189.12.260.260"))
	require.False(t, IsIPv4("189.260.260.260"))
	require.False(t, IsIPv4("255.255.255.256"))
	require.False(t, IsIPv4("999.999.999.999"))
	require.False(t, IsIPv4("9999.9999.9999.9999"))
	require.False(t, IsIPv4("192168.1.1"))
	require.False(t, IsIPv4("192.1681.1"))
	require.False(t, IsIPv4("192a168.1.1"))
}

func Test_IsIPv6(t *testing.T) {
	t.Parallel()

	// Valid Cases
	require.True(t, IsIPv6("9396:9549:b4f7:8ed0:4791:1330:8c06:e62d"))
	require.True(t, IsIPv6("2345:0425:2CA1::0567:5673:23b5"))
	require.True(t, IsIPv6("2001:1:2:3:4:5:6:7"))
	require.True(t, IsIPv6("::1"), "IPv6 loopback address with leading ellipsis")
	require.True(t, IsIPv6("1::"), "Shortest possible IPv6 address with trailing ellipsis")
	require.True(t, IsIPv6("2001:db8::8a2e:370:7334"), "Valid IPv6 address with ellipsis in the middle")
	require.True(t, IsIPv6("::ffff:192.0.2.128"), "IPv4-mapped IPv6 address")
	require.True(t, IsIPv6("::192.168.1.1"), "IPv4-compatible IPv6 address")
	require.True(t, IsIPv6("ffff::"), "Address with a single non-zero group and leading ellipsis")

	// Boundary values
	require.True(t, IsIPv6("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), "Maximum value for all groups")

	// Invalid cases
	require.False(t, IsIPv6("1::2::3"), "IPv6 address with multiple ellipses")
	require.False(t, IsIPv6("2001::25de::cade"), "IPv6 address with multiple ellipses in the middle")
	require.False(t, IsIPv6("g:a:b:c:d:e:f:1"), "Invalid character 'g' in address")
	require.False(t, IsIPv6("::192.168.1.256"), "Invalid IPv4 segment in IPv6 address")
	require.False(t, IsIPv6("1:2:3:4:5:6:7:8:9"), "IPv6 address too long")
	require.False(t, IsIPv6("1:2:3"), "IPv6 address too short")
	require.False(t, IsIPv6("1:::3:4"), "Consecutive colons not part of a valid ellipsis")
	require.False(t, IsIPv6("::ffff:256.0.0.1"), "Embedded IPv4 address with invalid IPv4 segment")
	require.False(t, IsIPv6("1111:2222:3333:4444:5555:6666:7777:8888:9999"), "IPv6 address with too many groups")
	require.False(t, IsIPv6("2001:db8::8a2e:370g:7334"), "IPv6 address with invalid hex character 'g'")
	require.False(t, IsIPv6("1.1.1.1"))
	require.False(t, IsIPv6("2001:1:2:3:4:5:6:"))
	require.False(t, IsIPv6(":1:2:3:4:5:6:"))
	require.False(t, IsIPv6("1:2:3:4:5:6:"))
	require.False(t, IsIPv6(""))
	require.False(t, IsIPv6("invalid"))
}

func Test_IsIPv6_EdgeCases(t *testing.T) {
	t.Parallel()

	// hex group larger than 4 digits
	require.False(t, IsIPv6("1:2:3:4:5:6:7:12345"))

	// IPv4 part at wrong position
	require.False(t, IsIPv6("1:2:3:4:1.2.3.4"))

	// IPv4 part would exceed length
	require.False(t, IsIPv6("1:2:3:4:5:6:7:1.2.3.4"))

	// ellipsis present but no zero group
	require.False(t, IsIPv6("1::2:3:4:5:6:7:8"))
}

func Test_IPWhitespace(t *testing.T) {
	t.Parallel()

	require.False(t, IsIPv4(" 1.1.1.1"))
	require.False(t, IsIPv4("1.1.1.1 "))
	require.False(t, IsIPv6(" ::1"))
	require.False(t, IsIPv6("::1 "))
}

// go test -v -run=^$ -bench=Benchmark_IsIPv4 -benchmem -count=6
func Benchmark_IsIPv4(b *testing.B) {
	ip := "174.23.33.100"
	var res bool

	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = IsIPv4(ip)
		}
		require.True(b, res)
	})

	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = net.ParseIP(ip) != nil
		}
		require.True(b, res)
	})
}

// go test -v -run=^$ -bench=Benchmark_IsIPv6 -benchmem -count=6
func Benchmark_IsIPv6(b *testing.B) {
	ip := "9396:9549:b4f7:8ed0:4791:1330:8c06:e62d"
	var res bool

	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = IsIPv6(ip)
		}
		require.True(b, res)
	})

	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = net.ParseIP(ip) != nil
		}
		require.True(b, res)
	})
}
