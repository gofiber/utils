package swar

import (
	"encoding/binary"
	"math/rand"
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

func Test_Store8_LaneOrder_RoundTrip(t *testing.T) {
	t.Parallel()
	b := make([]byte, 12)
	Store8(b, 2, 0x0807060504030201)
	require.Equal(t, []byte{0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0}, b)
	// Load8/Store8 round trip at every offset, leaving neighbors untouched.
	src := []byte{0xFF, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xAA, 0xBB}
	for i := 0; i+8 <= len(src); i++ {
		dst := make([]byte, len(src))
		copy(dst, src)
		Store8(dst, i, Load8(src, i))
		require.Equal(t, src, dst, "offset %d", i)
	}
}

func Test_Broadcast(t *testing.T) {
	t.Parallel()
	for _, c := range []byte{0x00, 0x01, '0', 'a', 0x7F, 0x80, 0xFF} {
		want := wordFromLanes([8]byte{c, c, c, c, c, c, c, c})
		require.Equal(t, want, Broadcast(c), "byte 0x%02x", c)
	}
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
	// Plain checks keep the ~1M words fast, as in the range test below.
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
						if got := byte(m >> (8 * i)); got != want {
							t.Fatalf("needle 0x%02x lanes %v lane %d: got 0x%02x want 0x%02x",
								needle, lanes, i, got, want)
						}
					}
				}
			}
		}
	}
}

func Test_MatchRangeMask_Exhaustive(t *testing.T) {
	t.Parallel()
	// All valid (lo, hi) pairs, with every byte value rotating through the
	// lanes. The neighbor lanes hold the range boundaries and both extremes
	// rather than a constant fill: the carry-freedom claim is exactly about
	// adjacent-lane interaction. Plain checks keep the ~2M words fast.
	for lo := 0; lo <= 0x7F; lo++ {
		for hi := lo; hi <= 0x7F; hi++ {
			base := [8]byte{0x00, byte(lo), byte(hi), byte(lo) - 1, byte(hi) + 1, 0x7F, 0x80, 0xFF}
			for c := range 256 {
				lanes := base
				lanes[c&7] = byte(c)
				m := MatchRangeMask(wordFromLanes(lanes), byte(lo), byte(hi))
				for i, v := range lanes {
					want := uint64(0)
					if v >= byte(lo) && v <= byte(hi) && v < 0x80 {
						want = 0x80
					}
					if got := m >> (8 * i) & 0xFF; got != want {
						t.Fatalf("range [0x%02x,0x%02x] lanes %v lane %d: got 0x%02x want 0x%02x",
							lo, hi, lanes, i, got, want)
					}
				}
			}
		}
	}

	// Random valid ranges over fully random lanes, as a second angle on
	// mixed neighbor values.
	rng := rand.New(rand.NewSource(11)) //nolint:gosec // deterministic test data
	for range 200000 {
		lo := byte(rng.Intn(0x80))
		hi := lo + byte(rng.Intn(0x80-int(lo)))
		var lanes [8]byte
		for i := range lanes {
			lanes[i] = byte(rng.Intn(256))
		}
		m := MatchRangeMask(wordFromLanes(lanes), lo, hi)
		for i, v := range lanes {
			want := uint64(0)
			if v >= lo && v <= hi && v < 0x80 {
				want = 0x80
			}
			if got := m >> (8 * i) & 0xFF; got != want {
				t.Fatalf("range [0x%02x,0x%02x] lanes %v lane %d: got 0x%02x want 0x%02x",
					lo, hi, lanes, i, got, want)
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
// front. The guard is advisory only, it cannot fail CI; eyeball it (or
// benchstat it against the previous toolchain) when touching Load8/Store8
// or bumping the Go version. The store leg covers Store8's write fusion.
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
	b.Run("store", func(b *testing.B) {
		b.SetBytes(int64(len(buf)))
		w := Broadcast('x')
		for b.Loop() {
			for i := 0; i+8 <= len(buf); i += 8 {
				Store8(buf, i, w)
			}
		}
	})
}

func Test_ZeroLanes_FirstMatchProperty(t *testing.T) {
	t.Parallel()
	// For every word shape: the mask is zero iff no lane is zero, and the
	// lowest set lane is always a true zero lane (false positives may only
	// appear above it).
	for zeroAt := range 8 {
		for fill := range 256 {
			var lanes [8]byte
			for i := range lanes {
				lanes[i] = byte(fill)
			}
			lanes[zeroAt] = 0
			m := ZeroLanes(wordFromLanes(lanes))
			require.NotZero(t, m, "fill 0x%02x zero at %d", fill, zeroAt)
			lo := FirstLane(m)
			require.Zero(t, lanes[lo], "fill 0x%02x zero at %d: first lane %d", fill, zeroAt, lo)
			if byte(fill) != 0 {
				require.Equal(t, zeroAt, lo)
			}
		}
	}
	for fill := 1; fill < 256; fill++ {
		var lanes [8]byte
		for i := range lanes {
			lanes[i] = byte(fill)
		}
		require.Zero(t, ZeroLanes(wordFromLanes(lanes)), "fill 0x%02x", fill)
	}
}
