package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// refSplitTrim is the reference for SplitTrimSeq: split, trim, drop empties.
func refSplitTrim(s string, sep byte) []string {
	var out []string
	for _, part := range strings.Split(s, string(sep)) {
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

	// Breaking out of the loop stops the iteration after the first element.
	count := 0
	for elem := range SplitTrimSeq("a, b, c", ',') {
		require.Equal(t, "a", elem)
		count++
		break
	}
	require.Equal(t, 1, count)

	// Elements alias the input rather than copying it.
	in := []byte(" gzip , br")
	for elem := range SplitTrimSeq(in, ',') {
		require.Same(t, &in[1], &elem[0])
		break
	}
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
				for range SplitTrimSeq(input.value, ',') {
					count++
				}
			}
		})
		b.Run(input.name+"/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for part := range strings.SplitSeq(input.value, ",") {
					if TrimSpace(part) == "" {
						continue
					}
					count++
				}
			}
		})
	}
	_ = count
}
