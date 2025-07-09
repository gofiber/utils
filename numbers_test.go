package utils

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseUint(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want uint64
		err  error
	}{
		{"0", 0, nil},
		{"42", 42, nil},
		{"18446744073709551615", math.MaxUint64, nil},
		{"", 0, strconv.ErrSyntax},
		{"123abc", 0, strconv.ErrSyntax},
		{"18446744073709551616", 0, strconv.ErrRange},
		{"-1", 0, strconv.ErrSyntax},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			v, err := ParseUint(tc.in)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			}
		})
		t.Run(tc.in+"_bytes", func(t *testing.T) {
			v, err := ParseUint([]byte(tc.in))
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			}
		})
	}
}

func Benchmark_ParseUint(b *testing.B) {
	str := "1234567890"
	var res uint64
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res, _ = ParseUint(str)
	}
	require.Equal(b, uint64(1234567890), res)
}

func Benchmark_StdParseUint(b *testing.B) {
	str := "1234567890"
	var res uint64
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res, _ = strconv.ParseUint(str, 10, 64)
	}
	require.Equal(b, uint64(1234567890), res)
}
