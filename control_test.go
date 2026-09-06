package utils

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

// refIndexControl is the scalar reference for IndexControl (exempt ==
// noControlExemption) and IndexControlExceptTab (exempt == '\t').
func refIndexControl(s string, exempt byte) int {
	for i := range len(s) {
		c := s[i]
		if (c < 0x20 || c == 0x7f) && c != exempt {
			return i
		}
	}
	return -1
}

func Test_IndexControl(t *testing.T) {
	t.Parallel()
	check := func(s string) {
		t.Helper()
		require.Equal(t, refIndexControl(s, noControlExemption), IndexControl(s), "IndexControl(%q)", s)
		require.Equal(t, refIndexControl(s, noControlExemption), IndexControl([]byte(s)), "IndexControl(%q) bytes", s)
		require.Equal(t, refIndexControl(s, '\t'), IndexControlExceptTab(s), "IndexControlExceptTab(%q)", s)
		require.Equal(t, refIndexControl(s, '\t'), IndexControlExceptTab([]byte(s)), "IndexControlExceptTab(%q) bytes", s)
	}

	// Every byte value at every position of a clean prefix of every length
	// up to two words past the unrolled loop, so each scan shape (byte-wise,
	// one word, two words, overlapping tail) sees each byte.
	clean := strings.Repeat("abcdefghijklmnopqrstuvwxyz0123456789", 2)
	for n := 0; n <= 40; n++ {
		check(clean[:n])
		for c := range 256 {
			buf := []byte(clean[:n])
			for pos := 0; pos <= n; pos++ {
				b := append(append(append([]byte{}, buf[:pos]...), byte(c)), buf[pos:]...)
				check(string(b))
			}
		}
	}

	// Bytes >= 0x80 never match, including UTF-8 encoded C1 controls that
	// unicode.IsControl would flag.
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
		require.Equal(t, refIndexControl(s, noControlExemption), IndexControl(s))
		require.Equal(t, refIndexControl(s, noControlExemption), IndexControl([]byte(s)))
		require.Equal(t, refIndexControl(s, '\t'), IndexControlExceptTab(s))
		require.Equal(t, refIndexControl(s, '\t'), IndexControlExceptTab([]byte(s)))
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
