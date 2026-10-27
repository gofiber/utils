package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func refSplitTrim(s string, sep byte) []string {
	var out []string
	for part := range strings.SplitSeq(s, string(sep)) {
		if part = TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func splitTrimSamples() []string {
	return []string{
		"",
		",",
		" , ,",
		"\t,\r\n",
		"a",
		" a ",
		"a,b",
		"a, b ,c",
		",,a,,",
		"a,",
		",a",
		"gzip, deflate, br",
		"text/html;q=0.9, */*;q=0.8",
		"keep-alive,\tUpgrade",
		" leading, trailing ",
		"no separator at all",
	}
}

func Test_SplitTrimSeq(t *testing.T) {
	t.Parallel()
	for _, in := range splitTrimSamples() {
		for _, sep := range []byte{',', ';', ' '} {
			want := refSplitTrim(in, sep)
			var got []string
			for elem := range SplitTrimSeq(in, sep) {
				got = append(got, elem)
			}
			require.Equal(t, want, got, "SplitTrimSeq(%q, %q)", in, sep)

			var gotBytes []string
			for elem := range SplitTrimSeq([]byte(in), sep) {
				gotBytes = append(gotBytes, string(elem))
			}
			require.Equal(t, want, gotBytes, "SplitTrimSeq(%q, %q) bytes", in, sep)
		}
	}

	count := 0
	seq := SplitTrimSeq("a, b, c", ',')
	for elem := range seq {
		require.Equal(t, "a", elem)
		count++
		break
	}
	require.Equal(t, 1, count)

	// A retained sequence restarts after a break.
	for range 2 {
		var got []string
		for elem := range seq {
			got = append(got, elem)
		}
		require.Equal(t, []string{"a", "b", "c"}, got)
	}

	in := []byte(" gzip , br")
	for elem := range SplitTrimSeq(in, ',') {
		require.Same(t, &in[1], &elem[0])
		break
	}
}

// countSplitTrimSeq and countSplitSeq range outside b.Loop, whose body disables inlining.
func countSplitTrimSeq(s string, sep byte) int {
	n := 0
	for range SplitTrimSeq(s, sep) {
		n++
	}
	return n
}

func countSplitSeq(s, sep string) int {
	n := 0
	for part := range strings.SplitSeq(s, sep) {
		if TrimSpace(part) != "" {
			n++
		}
	}
	return n
}

func Benchmark_SplitTrimSeq(b *testing.B) {
	inputs := []struct {
		name  string
		value string
	}{
		{"accept-encoding", "gzip, deflate, br"},
		{"accept", "text/html, application/xhtml+xml, application/xml;q=0.9, image/avif, image/webp, */*;q=0.8"},
	}
	var count int
	for _, input := range inputs {
		b.Run(input.name+"/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				count += countSplitTrimSeq(input.value, ',')
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				count += countSplitSeq(input.value, ",")
			}
		})
	}
	_ = count
}
