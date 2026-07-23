package utils

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AppendHexEncode(t *testing.T) {
	t.Parallel()
	all := make([]byte, 256)
	for i := range all {
		all[i] = byte(i)
	}

	// Every length from empty through several SWAR words, sweeping the
	// 8-byte/4-byte/scalar split, plus the full byte alphabet.
	for length := 0; length <= 3*8+3; length++ {
		src := all[100 : 100+length]
		want := hex.EncodeToString(src)
		require.Equal(t, want, string(AppendHexEncode(nil, src)))
		require.Equal(t, want, string(AppendHexEncode(nil, string(src))))
	}
	require.Equal(t, hex.EncodeToString(all), string(AppendHexEncode(nil, all)))

	// Appending must preserve existing dst content.
	got := AppendHexEncode([]byte("id="), []byte{0xDE, 0xAD, 0xBE, 0xEF})
	require.Equal(t, "id=deadbeef", string(got))

	// Round trip through the stdlib decoder.
	decoded, err := hex.DecodeString(string(AppendHexEncode(nil, "hello, hex world!")))
	require.NoError(t, err)
	require.Equal(t, "hello, hex world!", string(decoded))
}

func Benchmark_AppendHexEncode(b *testing.B) {
	for _, size := range []int{8, 64, 512} {
		src := make([]byte, size)
		for i := range src {
			src[i] = byte(i * 7)
		}
		dst := make([]byte, 0, 2*size)
		b.Run(strconv.Itoa(size)+"B/fiber", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				AppendHexEncode(dst, src)
			}
		})
		b.Run(strconv.Itoa(size)+"B/default", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				hex.AppendEncode(dst, src)
			}
		})
	}
}
