package utils

import (
	"slices"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

// refIndexControl is the scalar reference; exempt is noControlExemption or '\t'.
func refIndexControl(s string, exempt byte) int {
	for i := range len(s) {
		c := s[i]
		if (c < 0x20 || c == 0x7f) && c != exempt {
			return i
		}
	}
	return -1
}

// checkIndexControl compares both scanners on both input types with the reference.
func checkIndexControl(tb testing.TB, s string) {
	tb.Helper()
	if want := refIndexControl(s, noControlExemption); IndexControl(s) != want || IndexControl([]byte(s)) != want {
		tb.Fatalf("IndexControl(%q) = %d/%d, want %d", s, IndexControl(s), IndexControl([]byte(s)), want)
	}
	if want := refIndexControl(s, '\t'); IndexControlExceptTab(s) != want || IndexControlExceptTab([]byte(s)) != want {
		tb.Fatalf("IndexControlExceptTab(%q) = %d/%d, want %d", s, IndexControlExceptTab(s), IndexControlExceptTab([]byte(s)), want)
	}
}

func Test_IndexControl(t *testing.T) {
	t.Parallel()
	// Every byte value at every position of clean prefixes up to 40 bytes, so each scan shape sees each byte.
	clean := strings.Repeat("abcdefghijklmnopqrstuvwxyz0123456789", 2)
	for n := 0; n <= 40; n++ {
		checkIndexControl(t, clean[:n])
		for c := range 256 {
			for pos := 0; pos <= n; pos++ {
				checkIndexControl(t, string(slices.Insert([]byte(clean[:n]), pos, byte(c))))
			}
		}
	}

	// Bytes >= 0x80 never match, including the UTF-8 C1 controls unicode.IsControl flags.
	require.Equal(t, -1, IndexControl("caf\xc3\xa9 \xc2\x85 \xff\x80\x9f"))
	require.Equal(t, 6, strings.IndexFunc("caf\xc3\xa9 \xc2\x85", unicode.IsControl))
	require.Equal(t, -1, IndexControlExceptTab("a\tb\tc"))
	require.Equal(t, 1, IndexControl("a\tb\tc"))
	require.Equal(t, 3, IndexControlExceptTab("a\tb\x7f"))
	require.Equal(t, -1, IndexControl(""))
	require.Equal(t, 0, IndexControl("\x00"))
	// Only HTAB is exempt: the neighboring controls still match.
	require.Equal(t, 0, IndexControlExceptTab("\x08"))
	require.Equal(t, 0, IndexControlExceptTab("\n"))
}

func FuzzIndexControl(f *testing.F) {
	f.Add("")
	f.Add("plain header value")
	f.Add("tab\tseparated\tvalue")
	f.Add("crlf\r\ninjection")
	f.Add("del\x7fbyte")
	f.Add("utf8 caf\xc3\xa9 and \xc2\x85")
	f.Add(strings.Repeat("x", 31) + "\x01")
	f.Fuzz(func(t *testing.T, s string) {
		checkIndexControl(t, s)
	})
}

func Benchmark_IndexControl(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"clean-16B", "text/html;charse"},
		{"clean-64B", strings.Repeat("no-control-bytes", 4)},
		{"ctl-at-end-64B", strings.Repeat("no-control-bytes", 4)[:63] + "\r"},
	}
	var res int
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res = IndexControl(input.value)
			}
		})
		b.Run(input.name+"/fiber-except-tab", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res = IndexControlExceptTab(input.value)
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res = strings.IndexFunc(input.value, unicode.IsControl)
			}
		})
	}
	_ = res
}
