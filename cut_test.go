package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func cutSamples() []string {
	return []string{
		"",
		";",
		";;",
		"a",
		"a;",
		";a",
		"a;b",
		"a;b;c",
		"text/html; charset=utf-8",
		"host:8080",
		"[::1]:8080",
		"no separator here",
		"trailing;;",
		"\x00;\xff",
	}
}

func Test_CutByte(t *testing.T) {
	t.Parallel()
	for _, s := range cutSamples() {
		for _, sep := range []byte{';', ':', 'a', '\x00'} {
			wantBefore, wantAfter, wantFound := strings.Cut(s, string(sep))
			before, after, found := CutByte(s, sep)
			require.Equal(t, wantFound, found, "CutByte(%q, %q)", s, sep)
			require.Equal(t, wantBefore, before, "CutByte(%q, %q)", s, sep)
			require.Equal(t, wantAfter, after, "CutByte(%q, %q)", s, sep)

			bb, ba, bf := CutByte([]byte(s), sep)
			wb, wa, wf := bytes.Cut([]byte(s), []byte{sep})
			require.Equal(t, wf, bf, "CutByte(%q, %q) bytes", s, sep)
			require.Equal(t, wb, bb, "CutByte(%q, %q) bytes", s, sep)
			require.Equal(t, wa, ba, "CutByte(%q, %q) bytes", s, sep)
			if !found {
				// Like bytes.Cut, the after part of a miss is nil.
				require.Nil(t, ba)
			}

			// LastCutByte splits at the last separator.
			lb, la, lf := LastCutByte(s, sep)
			if i := strings.LastIndexByte(s, sep); i >= 0 {
				require.True(t, lf)
				require.Equal(t, s[:i], lb)
				require.Equal(t, s[i+1:], la)
			} else {
				require.False(t, lf)
				require.Equal(t, s, lb)
				require.Empty(t, la)
			}
			lbb, lba, lbf := LastCutByte([]byte(s), sep)
			require.Equal(t, lf, lbf)
			require.Equal(t, lb, string(lbb))
			require.Equal(t, la, string(lba))
			if !lf {
				require.Nil(t, lba)
			}
		}
	}

	// The results alias the input: no copies are made.
	in := []byte("key=value")
	k, v, ok := CutByte(in, '=')
	require.True(t, ok)
	require.Same(t, &in[0], &k[0])
	require.Same(t, &in[4], &v[0])
}

func Benchmark_CutByte(b *testing.B) {
	input := "text/html; charset=utf-8"
	inputBytes := []byte(input)
	var before, after string
	var found bool
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			before, after, found = CutByte(input, ';')
		}
		require.True(b, found)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			before, after, found = strings.Cut(input, ";")
		}
		require.True(b, found)
	})
	_, _ = before, after
	var bb, ba []byte
	b.Run("fiber-bytes", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			bb, ba, found = CutByte(inputBytes, ';')
		}
		require.True(b, found)
	})
	b.Run("default-bytes", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			bb, ba, found = bytes.Cut(inputBytes, []byte{';'})
		}
		require.True(b, found)
	})
	_, _ = bb, ba
}

func Benchmark_LastCutByte(b *testing.B) {
	input := "[2001:db8::1]:8080"
	var before, after string
	var found bool
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			before, after, found = LastCutByte(input, ':')
		}
		require.True(b, found)
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if i := strings.LastIndexByte(input, ':'); i >= 0 {
				before, after, found = input[:i], input[i+1:], true
			}
		}
		require.True(b, found)
	})
	_, _ = before, after
}
