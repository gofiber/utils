package simd

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// The tests compare every exported function AND its portable fallback
// against plain scalar references, sweeping sizes across the dispatch
// boundaries (sub-word, sub-vector, 32+ bytes, unaligned tails) so both the
// AVX2 kernels and the SWAR paths are exercised on amd64.

// testSizes crosses every dispatch boundary: scalar (<8), SWAR words with
// tails (8..31), one vector (32..63), several vectors plus tails, and sizes
// beyond the Memmem bytes.Index crossover (128+).
var testSizes = []int{
	0, 1, 2, 3, 7, 8, 9, 15, 16, 17, 23, 24, 31, 32, 33, 39, 47, 48, 63, 64,
	65, 95, 96, 100, 127, 128, 129, 130, 200, 255, 256, 257, 511, 512, 1000,
}

// fill produces deterministic pseudo-random bytes over the full 0..255
// range, seeded per call site, without math/rand (xorshift is plenty here).
func fill(n int, seed uint64) []byte {
	b := make([]byte, n)
	x := seed*2654435761 + 1
	for i := range b {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		b[i] = byte(x)
	}
	return b
}

func refMemchr2(h []byte, a, b byte) int {
	for i, c := range h {
		if c == a || c == b {
			return i
		}
	}
	return -1
}

func refMemchr3(h []byte, a, b, c byte) int {
	for i, x := range h {
		if x == a || x == b || x == c {
			return i
		}
	}
	return -1
}

func refMemchrPair(h []byte, b1, b2 byte, offset int) int {
	if offset < 0 {
		return -1
	}
	for i := 0; i+offset < len(h); i++ {
		if h[i] == b1 && h[i+offset] == b2 {
			return i
		}
	}
	return -1
}

func refIsASCII(h []byte) bool {
	for _, b := range h {
		if b >= 0x80 {
			return false
		}
	}
	return true
}

func checkMemchr2(t *testing.T, h []byte, a, b byte) {
	t.Helper()
	want := refMemchr2(h, a, b)
	require.Equal(t, want, Memchr2(h, a, b), "Memchr2(%d bytes, %#x, %#x)", len(h), a, b)
	require.Equal(t, want, memchr2Generic(h, a, b), "memchr2Generic(%d bytes, %#x, %#x)", len(h), a, b)
}

func checkMemchr3(t *testing.T, h []byte, a, b, c byte) {
	t.Helper()
	want := refMemchr3(h, a, b, c)
	require.Equal(t, want, Memchr3(h, a, b, c), "Memchr3(%d bytes)", len(h))
	require.Equal(t, want, memchr3Generic(h, a, b, c), "memchr3Generic(%d bytes)", len(h))
}

func checkMemchrPair(t *testing.T, h []byte, b1, b2 byte, offset int) {
	t.Helper()
	want := refMemchrPair(h, b1, b2, offset)
	require.Equal(t, want, MemchrPair(h, b1, b2, offset), "MemchrPair(%d bytes, %#x, %#x, %d)", len(h), b1, b2, offset)
	if offset >= 1 && len(h) > offset {
		require.Equal(t, want, memchrPairGeneric(h, b1, b2, offset), "memchrPairGeneric(%d bytes, %#x, %#x, %d)", len(h), b1, b2, offset)
	}
}

func Test_Memchr2(t *testing.T) {
	t.Parallel()
	for _, n := range testSizes {
		// Absent needles force a full scan of clean data.
		clean := bytes.Repeat([]byte{'x'}, n)
		checkMemchr2(t, clean, 'a', 'b')
		// Every single position, for either needle.
		for pos := range n {
			h := bytes.Repeat([]byte{'x'}, n)
			h[pos] = 'a'
			checkMemchr2(t, h, 'a', 'b')
			h[pos] = 'b'
			checkMemchr2(t, h, 'a', 'b')
			// Identical needles must behave like a single-byte search.
			checkMemchr2(t, h, 'b', 'b')
			h[pos] = 0x80 // high-bit byte must not confuse the SWAR masks
			checkMemchr2(t, h, 0x80, 0xff)
		}
		// First match wins when both needles are present.
		if n >= 2 {
			h := bytes.Repeat([]byte{'x'}, n)
			h[n/2] = 'b'
			h[n-1] = 'a'
			checkMemchr2(t, h, 'a', 'b')
		}
		// Random data over the full byte range.
		checkMemchr2(t, fill(n, uint64(n)), 0x00, 0x01)
	}
}

func Test_Memchr2_SWARBorrow(t *testing.T) {
	t.Parallel()
	// 0x00/0x01 adjacency is where the SWAR zero-lane formula produces
	// false positives above the first true match; the first-lane pick must
	// stay exact.
	patterns := [][]byte{
		{0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00},
		{0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01},
		bytes.Repeat([]byte{0x01}, 40),
		append(bytes.Repeat([]byte{0x01}, 39), 0x00),
	}
	for _, h := range patterns {
		checkMemchr2(t, h, 0x00, 0x02)
		checkMemchr2(t, h, 0x02, 0x00)
		checkMemchr3(t, h, 0x00, 0x02, 0x03)
		checkMemchrPair(t, h, 0x01, 0x00, 1)
		checkMemchrPair(t, h, 0x00, 0x00, 2)
	}
}

func Test_Memchr3(t *testing.T) {
	t.Parallel()
	for _, n := range testSizes {
		clean := bytes.Repeat([]byte{'x'}, n)
		checkMemchr3(t, clean, 'a', 'b', 'c')
		for pos := range n {
			for _, needle := range []byte{'a', 'b', 'c'} {
				h := bytes.Repeat([]byte{'x'}, n)
				h[pos] = needle
				checkMemchr3(t, h, 'a', 'b', 'c')
			}
		}
		if n >= 3 {
			h := bytes.Repeat([]byte{'x'}, n)
			h[n-1] = 'a'
			h[n/2] = 'b'
			h[n/3] = 'c'
			checkMemchr3(t, h, 'a', 'b', 'c')
		}
		checkMemchr3(t, fill(n, uint64(n)+1), 0x00, 0x80, 0xff)
	}
}

func Test_MemchrPair(t *testing.T) {
	t.Parallel()
	// Degenerate parameters.
	require.Equal(t, -1, MemchrPair([]byte("abc"), 'a', 'b', -1))
	require.Equal(t, -1, MemchrPair([]byte("abc"), 'a', 'b', 3))
	require.Equal(t, -1, MemchrPair([]byte("abc"), 'a', 'b', 7))
	require.Equal(t, -1, MemchrPair(nil, 'a', 'b', 0))
	// offset 0 needs both bytes at one position: only equal bytes can match.
	require.Equal(t, -1, MemchrPair([]byte("abc"), 'a', 'b', 0))
	require.Equal(t, 1, MemchrPair([]byte("abc"), 'b', 'b', 0))

	offsets := []int{1, 2, 3, 7, 8, 9, 15, 31, 32, 40}
	for _, n := range testSizes {
		for _, offset := range offsets {
			if offset >= n {
				continue
			}
			// Plant the pair at every valid start position.
			for pos := 0; pos+offset < n; pos++ {
				h := bytes.Repeat([]byte{'x'}, n)
				h[pos] = 'a'
				h[pos+offset] = 'b'
				checkMemchrPair(t, h, 'a', 'b', offset)
			}
			// Anchor byte present everywhere, partner absent.
			checkMemchrPair(t, bytes.Repeat([]byte{'a'}, n), 'a', 'b', offset)
			checkMemchrPair(t, fill(n, uint64(n*100+offset)), 'a', 'b', offset)
			checkMemchrPair(t, fill(n, uint64(n*100+offset+1)), 0x00, 0x01, offset)
		}
	}
}

func Test_IsASCII(t *testing.T) {
	t.Parallel()
	require.True(t, IsASCII(nil))
	for _, n := range testSizes {
		clean := make([]byte, n)
		for i := range clean {
			clean[i] = byte(i % 0x80)
		}
		require.True(t, IsASCII(clean), "size %d", n)
		require.True(t, isASCIIGeneric(clean), "generic size %d", n)
		for pos := range n {
			h := bytes.Clone(clean)
			h[pos] = 0x80
			require.False(t, IsASCII(h), "size %d pos %d", n, pos)
			require.False(t, isASCIIGeneric(h), "generic size %d pos %d", n, pos)
			h[pos] = 0xff
			require.False(t, IsASCII(h), "size %d pos %d", n, pos)
		}
	}
}

func Test_FirstNonASCII_CountNonASCII(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, FirstNonASCII(nil))
	require.Equal(t, 0, CountNonASCII(nil))
	require.Equal(t, -1, FirstNonASCII([]byte("hello")))
	require.Equal(t, 1, FirstNonASCII([]byte("h\xc3\xa9llo")))
	require.Equal(t, 2, CountNonASCII([]byte("h\xc3\xa9llo")))
	// The only non-ASCII byte sits in the overlapping final word (length
	// not a multiple of 8, clean words before it).
	require.Equal(t, 12, FirstNonASCII(append(bytes.Repeat([]byte{'a'}, 12), 0x80)))
	require.Equal(t, -1, FirstNonASCII(bytes.Repeat([]byte{'a'}, 13)))
	for _, n := range testSizes {
		h := fill(n, uint64(n)+7)
		wantFirst := -1
		wantCount := 0
		for i, b := range h {
			if b >= 0x80 {
				if wantFirst == -1 {
					wantFirst = i
				}
				wantCount++
			}
		}
		require.Equal(t, wantFirst, FirstNonASCII(h), "size %d", n)
		require.Equal(t, wantCount, CountNonASCII(h), "size %d", n)
		require.Equal(t, wantFirst == -1, IsASCII(h), "size %d", n)
	}
}

func Test_MemchrDigit(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, MemchrDigit(nil))
	// Every byte value alone: exactly '0'..'9' match.
	for b := range 256 {
		want := -1
		if b >= '0' && b <= '9' {
			want = 0
		}
		h := []byte{byte(b)}
		require.Equal(t, want, MemchrDigit(h), "byte %#x", b)
		require.Equal(t, want, memchrDigitGeneric(h), "generic byte %#x", b)
	}
	for _, n := range testSizes {
		clean := bytes.Repeat([]byte{'x'}, n)
		require.Equal(t, -1, MemchrDigit(clean), "size %d", n)
		for pos := range n {
			for _, d := range []byte{'0', '5', '9'} {
				h := bytes.Repeat([]byte{'x'}, n)
				h[pos] = d
				require.Equal(t, pos, MemchrDigit(h), "size %d pos %d", n, pos)
				require.Equal(t, pos, memchrDigitGeneric(h), "generic size %d pos %d", n, pos)
			}
			// Range neighbors must not match.
			h := bytes.Repeat([]byte{'x'}, n)
			h[pos] = '0' - 1 // '/'
			require.Equal(t, -1, MemchrDigit(h))
			h[pos] = '9' + 1 // ':'
			require.Equal(t, -1, MemchrDigit(h))
			h[pos] = 0xb5 // '5' with the high bit set: signed-compare trap
			require.Equal(t, -1, MemchrDigit(h))
		}
	}
}

func Test_MemchrDigitAt(t *testing.T) {
	t.Parallel()
	h := []byte("abc123def456")
	require.Equal(t, 3, MemchrDigitAt(h, 0))
	require.Equal(t, 3, MemchrDigitAt(h, 3))
	require.Equal(t, 5, MemchrDigitAt(h, 5))
	require.Equal(t, 9, MemchrDigitAt(h, 6))
	require.Equal(t, -1, MemchrDigitAt(h, 12))
	require.Equal(t, -1, MemchrDigitAt(h, 100))
	require.Equal(t, -1, MemchrDigitAt(h, -1))
	require.Equal(t, -1, MemchrDigitAt(nil, 0))
	// In-bounds start with no digit at or after it.
	require.Equal(t, -1, MemchrDigitAt([]byte("abc123xyz"), 6))
}

func refIsWord(b byte) bool {
	return b == '_' || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

func Test_MemchrWord_NotWord(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, MemchrWord(nil))
	require.Equal(t, -1, MemchrNotWord(nil))
	// Every byte value alone, including the range boundaries and their
	// unsigned-compare neighbors ('A'-1, 'Z'+1, 'a'-1, 'z'+1, '_'±1, 0x80+).
	for b := range 256 {
		h := []byte{byte(b)}
		wantWord := -1
		wantNot := 0
		if refIsWord(byte(b)) {
			wantWord = 0
			wantNot = -1
		}
		require.Equal(t, wantWord, MemchrWord(h), "byte %#x", b)
		require.Equal(t, wantNot, MemchrNotWord(h), "byte %#x", b)
	}
	for _, n := range testSizes {
		nonWord := bytes.Repeat([]byte{' '}, n)
		require.Equal(t, -1, MemchrWord(nonWord), "size %d", n)
		allWord := bytes.Repeat([]byte{'k'}, n)
		require.Equal(t, -1, MemchrNotWord(allWord), "size %d", n)
		for pos := range n {
			for _, w := range []byte{'A', 'Z', 'a', 'z', '0', '9', '_'} {
				h := bytes.Repeat([]byte{' '}, n)
				h[pos] = w
				require.Equal(t, pos, MemchrWord(h), "size %d pos %d byte %q", n, pos, w)
				require.Equal(t, pos, memchrWordGeneric(h), "generic size %d pos %d", n, pos)
			}
			for _, nw := range []byte{' ', '@', '[', '`', '{', '/', ':', 0x80, 0xff} {
				h := bytes.Repeat([]byte{'k'}, n)
				h[pos] = nw
				require.Equal(t, pos, MemchrNotWord(h), "size %d pos %d byte %#x", n, pos, nw)
				require.Equal(t, pos, memchrNotWordGeneric(h), "generic size %d pos %d", n, pos)
			}
		}
	}
}

func Test_MemchrInTable_NotInTable(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, MemchrInTable([]byte("abc"), nil))
	require.Equal(t, -1, MemchrNotInTable([]byte("abc"), nil))
	require.Equal(t, -1, MemchrInTable(nil, &[256]bool{}))
	require.Equal(t, -1, MemchrNotInTable(nil, &[256]bool{}))

	var vowels [256]bool
	for _, c := range []byte("aeiou") {
		vowels[c] = true
	}
	require.Equal(t, 1, MemchrInTable([]byte("halo"), &vowels))
	require.Equal(t, 0, MemchrNotInTable([]byte("halo"), &vowels))
	require.Equal(t, -1, MemchrInTable([]byte("xyz"), &vowels))
	require.Equal(t, -1, MemchrNotInTable([]byte("aeiou"), &vowels))

	for _, n := range testSizes {
		h := fill(n, uint64(n)+13)
		var table [256]bool
		inSet := make(map[byte]bool, 86)
		for i := range table {
			if i%3 == 0 {
				table[i] = true
				inSet[byte(i)] = true
			}
		}
		wantIn, wantNot := -1, -1
		for i, b := range h {
			if inSet[b] && wantIn == -1 {
				wantIn = i
			}
			if !inSet[b] && wantNot == -1 {
				wantNot = i
			}
		}
		require.Equal(t, wantIn, MemchrInTable(h, &table), "size %d", n)
		require.Equal(t, wantNot, MemchrNotInTable(h, &table), "size %d", n)
	}
}

func Test_SelectRareBytes(t *testing.T) {
	t.Parallel()
	require.Equal(t, RareByteInfo{}, SelectRareBytes(nil))
	require.Equal(t, RareByteInfo{Byte1: 'q', Byte2: 'q'}, SelectRareBytes([]byte("q")))

	// 'z' (rank 20) is rarer than 'e' (rank 245) and 'l' (rank 175).
	info := SelectRareBytes([]byte("zeel"))
	require.Equal(t, byte('z'), info.Byte1)
	require.Equal(t, 0, info.Index1)
	require.Equal(t, byte('l'), info.Byte2)
	require.Equal(t, 3, info.Index2)

	for _, n := range []int{2, 3, 5, 8, 16, 32, 100} {
		needle := fill(n, uint64(n)+29)
		info := SelectRareBytes(needle)
		// The reported bytes really sit at the reported indices.
		require.Equal(t, needle[info.Index1], info.Byte1, "needle %q", needle)
		require.Equal(t, needle[info.Index2], info.Byte2, "needle %q", needle)
		// Byte1 is a rarest byte of the whole needle.
		for _, b := range needle {
			require.LessOrEqual(t, ByteRank(info.Byte1), ByteRank(b), "needle %q", needle)
		}
		require.LessOrEqual(t, ByteRank(info.Byte1), ByteRank(info.Byte2))
	}

	// All-identical needle: no distinct second byte exists.
	info = SelectRareBytes(bytes.Repeat([]byte{'a'}, 10))
	require.Equal(t, byte('a'), info.Byte1)
	require.Equal(t, byte('a'), info.Byte2)
}

func Test_ByteRank(t *testing.T) {
	t.Parallel()
	require.Equal(t, byteFrequencies['e'], ByteRank('e'))
	require.Less(t, ByteRank('z'), ByteRank('e'))
	require.Equal(t, byte(255), ByteRank(' '))
}

func checkMemmem(t *testing.T, h, needle []byte) {
	t.Helper()
	require.Equal(t, bytes.Index(h, needle), Memmem(h, needle), "haystack %d bytes, needle %q", len(h), needle)
}

func Test_Memmem(t *testing.T) {
	t.Parallel()
	// bytes.Index contract on edge cases.
	require.Equal(t, 0, Memmem(nil, nil))
	require.Equal(t, 0, Memmem([]byte("abc"), nil))
	require.Equal(t, -1, Memmem(nil, []byte("a")))
	require.Equal(t, -1, Memmem([]byte("ab"), []byte("abc")))
	checkMemmem(t, []byte("hello world"), []byte("world"))
	checkMemmem(t, []byte("hello world"), []byte("xyz"))
	checkMemmem(t, []byte("aaaaaabaaaa"), []byte("aab"))

	needles := [][]byte{
		[]byte("z"), []byte("zq"), []byte("abc"), []byte("q7z"),
		[]byte("aaa"), []byte("aab"), []byte("the quick brown fox"),
		[]byte("\x00\x01"), []byte("\xff\xfe\xfd"),
		bytes.Repeat([]byte{'a'}, 33),
		append(bytes.Repeat([]byte{'a'}, 40), 'b'),
	}
	for _, n := range testSizes {
		// Text-like haystack (lowercase letters) and full-range haystack.
		text := make([]byte, n)
		for i := range text {
			text[i] = 'a' + byte((i*7)%26)
		}
		random := fill(n, uint64(n)+31)
		for _, needle := range needles {
			checkMemmem(t, text, needle)
			checkMemmem(t, random, needle)
			// Plant the needle at the start, middle, and end.
			if len(needle) <= n {
				for _, pos := range []int{0, (n - len(needle)) / 2, n - len(needle)} {
					h := bytes.Clone(text)
					copy(h[pos:], needle)
					checkMemmem(t, h, needle)
				}
			}
		}
		// Needle equals the whole haystack.
		if n > 0 {
			checkMemmem(t, text, text)
		}
	}

	// Repetitive data with many false candidates for the prefilter.
	big := bytes.Repeat([]byte("ab"), 300)
	checkMemmem(t, big, []byte("aba"))
	checkMemmem(t, big, []byte("abb"))
	checkMemmem(t, big, append(bytes.Repeat([]byte("ab"), 20), 'c'))
	almost := bytes.Repeat([]byte{'a'}, 500)
	almost[499] = 'b'
	checkMemmem(t, almost, []byte("aaab"))
	checkMemmem(t, almost, []byte("ab"))
}

// checkMemmemInternal drives the prefilter helpers directly, so they are
// covered on every platform regardless of the hasAVX2 / 128-byte routing
// inside Memmem.
func checkMemmemInternal(t *testing.T, h, needle []byte) {
	t.Helper()
	want := bytes.Index(h, needle)
	info := SelectRareBytes(needle)
	require.Equal(t, want, memmemSingle(h, needle, info.Byte1, info.Index1),
		"memmemSingle(%d bytes, %q)", len(h), needle)
	if info.Byte1 != info.Byte2 && info.Index1 != info.Index2 {
		require.Equal(t, want, memmemPaired(h, needle, info),
			"memmemPaired(%d bytes, %q)", len(h), needle)
	}
}

func Test_Memmem_InternalPaths(t *testing.T) {
	t.Parallel()
	// Regular haystacks, planted needles.
	for _, n := range testSizes {
		if n == 0 {
			continue
		}
		text := make([]byte, n)
		for i := range text {
			text[i] = 'a' + byte((i*7)%26)
		}
		for _, needle := range [][]byte{
			[]byte("zq"), []byte("q7z"), []byte("the quick brown fox"),
			append(bytes.Repeat([]byte{'a'}, 40), 'b'),
		} {
			checkMemmemInternal(t, text, needle)
			if len(needle) <= n {
				h := bytes.Clone(text)
				copy(h[n-len(needle):], needle)
				checkMemmemInternal(t, h, needle)
			}
		}
	}

	// Miss-budget fallback: every position is a rare-byte anchor, so the
	// loops exhaust memmemMaxMisses and must hand over to memmemFallback
	// without missing a late match or returning one out of order.
	anchors := bytes.Repeat([]byte{'z'}, 400)
	checkMemmemInternal(t, anchors, append(bytes.Repeat([]byte{'z'}, 50), 'b')) // absent
	late := bytes.Repeat([]byte{'z'}, 400)
	copy(late[350:], "zzzzb")
	checkMemmemInternal(t, late, []byte("zzzzb")) // found only after fallback
	// Match placed just before the fallback handover region, with a
	// nonzero anchor index (rare byte not at needle start).
	h := bytes.Repeat([]byte{'q'}, 300)
	for i := range 40 {
		h[i*2] = 'e' // 'e' is common; 'q' stays the rare anchor everywhere
	}
	copy(h[60:], "eeq")
	checkMemmemInternal(t, h, []byte("eeq"))
}

func Test_Memmem_AdversarialBounded(t *testing.T) {
	t.Parallel()
	// The O(n*m) trap from the code review: rare byte everywhere, long
	// almost-matching needle. The miss budget must keep this exact (equal
	// to bytes.Index) — timing is covered by Benchmark_Memmem_Adversarial.
	haystack := bytes.Repeat([]byte{'z'}, 1<<16)
	needle := append(bytes.Repeat([]byte{'z'}, 999), 'b')
	checkMemmem(t, haystack, needle)
	checkMemmemInternal(t, haystack, needle)
	// Same shape but with the needle present at the very end.
	present := bytes.Repeat([]byte{'z'}, 1<<16)
	present[len(present)-1] = 'b'
	checkMemmem(t, present, needle)
	checkMemmemInternal(t, present, needle)
}

func Test_Accelerated(t *testing.T) {
	t.Parallel()
	require.Equal(t, hasAVX2, Accelerated())
	t.Logf("simd.Accelerated() = %v", Accelerated())
	// CI runners known to have AVX2 can set this to prove the assembly
	// kernels are actually under test, instead of silently green-lighting
	// a fallback-only run (QEMU, GOAMD64-baseline VMs, detection bugs).
	if os.Getenv("GOFIBER_SIMD_REQUIRE_AVX2") == "1" {
		require.True(t, Accelerated(), "GOFIBER_SIMD_REQUIRE_AVX2=1 but AVX2 kernels are not active")
	}
}

func Test_TestSizesCoverDispatchBoundaries(t *testing.T) {
	t.Parallel()
	// Guard the sweep against accidental edits: it must include sub-word,
	// word-with-tail, exactly-one-vector, and beyond-crossover sizes.
	for _, must := range []int{0, 7, 8, 31, 32, 33, 63, 64, 127, 128, 512} {
		require.Contains(t, testSizes, must, strconv.Itoa(must))
	}
}
