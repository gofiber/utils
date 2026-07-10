package swar

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

// wordFromLanes builds a word with the given 8 lane values, lane 0 first.
func wordFromLanes(lanes [8]byte) uint64 {
	return binary.LittleEndian.Uint64(lanes[:])
}

func Test_Load8_LaneOrder(t *testing.T) {
	t.Parallel()
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := uint64(0x0807060504030201)
	require.Equal(t, want, Load8(b, 0))
	require.Equal(t, want, Load8(string(b), 0))
	require.Equal(t, uint64(0x0908070605040302), Load8(b, 1))
	require.Equal(t, uint64(0x0A09080706050403), Load8(string(b), 2))
}

func Test_ToLowerWord_ToUpperWord_Exhaustive(t *testing.T) {
	t.Parallel()
	// Every byte value in every lane position, once with all other lanes
	// zero and once with all other lanes 0xFF, so lane isolation is proven.
	for c := range 256 {
		wantLower := byte(c)
		if wantLower >= 'A' && wantLower <= 'Z' {
			wantLower += 'a' - 'A'
		}
		wantUpper := byte(c)
		if wantUpper >= 'a' && wantUpper <= 'z' {
			wantUpper -= 'a' - 'A'
		}
		for lane := range 8 {
			for _, fill := range []byte{0x00, 0xFF, 'm', 'M'} {
				var lanes [8]byte
				for i := range lanes {
					lanes[i] = fill
				}
				lanes[lane] = byte(c)
				w := wordFromLanes(lanes)

				gotLower := ToLowerWord(w)
				gotUpper := ToUpperWord(w)
				require.Equal(t, wantLower, byte(gotLower>>(8*lane)),
					"ToLowerWord lane %d byte 0x%02x fill 0x%02x", lane, c, fill)
				require.Equal(t, wantUpper, byte(gotUpper>>(8*lane)),
					"ToUpperWord lane %d byte 0x%02x fill 0x%02x", lane, c, fill)
			}
		}
	}
}

func Test_MatchByteMask_Exhaustive(t *testing.T) {
	t.Parallel()
	// For every needle value and every lane value, the mask must be exact:
	// 0x80 iff equal, 0 otherwise — including needle 0x00, needle >= 0x80,
	// and adjacent-lane interference (the borrow-corruption trap).
	for needle := range 256 {
		for c := range 256 {
			for lane := range 8 {
				// Fill other lanes with the needle itself (worst case for
				// cross-lane borrows) and with needle^1 (no other matches).
				for _, fill := range []byte{byte(needle), byte(needle) ^ 1} {
					var lanes [8]byte
					for i := range lanes {
						lanes[i] = fill
					}
					lanes[lane] = byte(c)
					m := MatchByteMask(wordFromLanes(lanes), byte(needle))
					for i := range lanes {
						want := byte(0)
						if lanes[i] == byte(needle) {
							want = 0x80
						}
						require.Equal(t, want, byte(m>>(8*i)),
							"needle 0x%02x lanes %v lane %d", needle, lanes, i)
					}
				}
			}
		}
	}
}

func Test_MatchRangeMask_Exhaustive(t *testing.T) {
	t.Parallel()
	ranges := [][2]byte{
		{'0', '9'},
		{'a', 'z'},
		{'A', 'Z'},
		{0x00, 0x1F},
		{0x09, 0x0D},
		{0x00, 0x7F},
		{0x20, 0x20},
		{0x7F, 0x7F},
		{0x00, 0x00},
	}
	for _, r := range ranges {
		lo, hi := r[0], r[1]
		for c := range 256 {
			for lane := range 8 {
				var lanes [8]byte
				for i := range lanes {
					lanes[i] = 0xFF // never in range: high bit set
				}
				lanes[lane] = byte(c)
				m := MatchRangeMask(wordFromLanes(lanes), lo, hi)
				want := uint64(0)
				if byte(c) >= lo && byte(c) <= hi && c < 0x80 {
					want = 0x80 << (8 * lane)
				}
				require.Equal(t, want, m, "range [0x%02x,0x%02x] byte 0x%02x lane %d", lo, hi, c, lane)
			}
		}
	}
}

func Test_FirstLane_LastLane(t *testing.T) {
	t.Parallel()
	require.Equal(t, 8, FirstLane(0))
	require.Equal(t, -1, LastLane(0))
	for first := range 8 {
		for last := first; last < 8; last++ {
			mask := uint64(0x80)<<(8*first) | uint64(0x80)<<(8*last)
			require.Equal(t, first, FirstLane(mask))
			require.Equal(t, last, LastLane(mask))
		}
	}
}

// Benchmark_Load8_Fusion guards the load-combining requirement: an unfused
// generic instantiation shows a ~4-8x slower word loop. If this benchmark
// ever collapses to per-byte speed, the reslice idiom in Load8 has stopped
// fusing and Load8 needs concrete string/[]byte helpers behind the generic
// front.
func Benchmark_Load8_Fusion(b *testing.B) {
	buf := make([]byte, 4096)
	for i := range buf {
		buf[i] = byte(i)
	}
	s := string(buf)
	b.Run("bytes", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		var acc uint64
		for b.Loop() {
			for i := 0; i+8 <= len(buf); i += 8 {
				acc += Load8(buf, i)
			}
		}
		_ = acc
	})
	b.Run("string", func(b *testing.B) {
		b.SetBytes(int64(len(s)))
		var acc uint64
		for b.Loop() {
			for i := 0; i+8 <= len(s); i += 8 {
				acc += Load8(s, i)
			}
		}
		_ = acc
	})
}
