package simd

import (
	"bytes"
	"testing"
)

// Seeds hit known traps: SWAR borrow adjacency, high-bit bytes on signed compares, vector tails.

func FuzzMemchr2(f *testing.F) {
	f.Add([]byte("hello world"), byte('o'), byte('w'))
	f.Add([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01}, byte(0), byte(2))
	f.Add(bytes.Repeat([]byte{'x'}, 65), byte('a'), byte('x'))
	f.Add([]byte{}, byte(0), byte(0))
	f.Fuzz(func(t *testing.T, h []byte, a, b byte) {
		want := refMemchr2(h, a, b)
		if got := Memchr2(h, a, b); got != want {
			t.Fatalf("Memchr2(%q, %#x, %#x) = %d, want %d", h, a, b, got, want)
		}
		if got := memchr2Generic(h, a, b); got != want {
			t.Fatalf("memchr2Generic(%q, %#x, %#x) = %d, want %d", h, a, b, got, want)
		}
	})
}

func FuzzMemchr3(f *testing.F) {
	f.Add([]byte("hello\tworld\nfoo"), byte(' '), byte('\t'), byte('\n'))
	f.Add([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01}, byte(0), byte(2), byte(3))
	f.Add(bytes.Repeat([]byte{0xff}, 40), byte(0x80), byte(0xfe), byte(0xff))
	f.Fuzz(func(t *testing.T, h []byte, a, b, c byte) {
		want := refMemchr3(h, a, b, c)
		if got := Memchr3(h, a, b, c); got != want {
			t.Fatalf("Memchr3(%q, %#x, %#x, %#x) = %d, want %d", h, a, b, c, got, want)
		}
		if got := memchr3Generic(h, a, b, c); got != want {
			t.Fatalf("memchr3Generic(%q, %#x, %#x, %#x) = %d, want %d", h, a, b, c, got, want)
		}
	})
}

func FuzzMemchrPair(f *testing.F) {
	f.Add([]byte("hello example world"), byte('e'), byte('x'), 1)
	f.Add([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01}, byte(1), byte(0), 1)
	f.Add(bytes.Repeat([]byte{'a'}, 100), byte('a'), byte('a'), 33)
	f.Add([]byte("abc"), byte('a'), byte('b'), -1)
	f.Fuzz(func(t *testing.T, h []byte, b1, b2 byte, offset int) {
		if offset > len(h) {
			offset %= len(h) + 1
		}
		want := refMemchrPair(h, b1, b2, offset)
		if offset == 0 && b1 != b2 {
			want = -1
		}
		if got := MemchrPair(h, b1, b2, offset); got != want {
			t.Fatalf("MemchrPair(%q, %#x, %#x, %d) = %d, want %d", h, b1, b2, offset, got, want)
		}
		if offset >= 1 && len(h) > offset {
			if got := memchrPairGeneric(h, b1, b2, offset); got != want {
				t.Fatalf("memchrPairGeneric(%q, %#x, %#x, %d) = %d, want %d", h, b1, b2, offset, got, want)
			}
		}
	})
}

func FuzzMemmem(f *testing.F) {
	f.Add([]byte("hello world"), []byte("world"))
	f.Add(bytes.Repeat([]byte("ab"), 100), []byte("aba"))
	f.Add(bytes.Repeat([]byte{'a'}, 200), []byte("aab"))
	f.Add([]byte{}, []byte{})
	// Miss-budget seeds: memmemSingle (all-'z' needle) and memmemPaired ('zqe' over dense false pairs)
	// reach memmemFallback.
	f.Add(bytes.Repeat([]byte{'z'}, 200), append(bytes.Repeat([]byte{'z'}, 30), 'b'))
	f.Add(bytes.Repeat([]byte("zqf"), 40), []byte("zqe"))
	f.Add(append(bytes.Repeat([]byte("zqf"), 40), 'z', 'q', 'e'), []byte("zqe"))
	f.Fuzz(func(t *testing.T, h, needle []byte) {
		want := bytes.Index(h, needle)
		if got := Memmem(h, needle); got != want {
			t.Fatalf("Memmem(%q, %q) = %d, want %d", h, needle, got, want)
		}
		// Prefilter helpers called directly: Memmem's routing (no AVX2, sub-128B haystack) would skip them.
		if len(needle) >= 2 && len(h) >= len(needle) {
			info := SelectRareBytes(needle)
			if got := memmemSingle(h, needle, info.Byte1, info.Index1); got != want {
				t.Fatalf("memmemSingle(%q, %q) = %d, want %d", h, needle, got, want)
			}
			if info.Byte1 != info.Byte2 {
				if got := memmemPaired(h, needle, info); got != want {
					t.Fatalf("memmemPaired(%q, %q) = %d, want %d", h, needle, got, want)
				}
			}
		}
	})
}

func FuzzIsASCII(f *testing.F) {
	f.Add([]byte("hello world"))
	f.Add([]byte("h\xc3\xa9llo"))
	f.Add(bytes.Repeat([]byte{0x7f}, 64))
	f.Add(append(bytes.Repeat([]byte{'a'}, 63), 0x80))
	f.Fuzz(func(t *testing.T, h []byte) {
		want := refIsASCII(h)
		if got := IsASCII(h); got != want {
			t.Fatalf("IsASCII(%q) = %v, want %v", h, got, want)
		}
		if got := isASCIIGeneric(h); got != want {
			t.Fatalf("isASCIIGeneric(%q) = %v, want %v", h, got, want)
		}
		// Check the index, not just its sign: a transposed AVX2 mask arm gives a
		// wrong-but-positive index.
		wantFirst := refFirstNonASCII(h)
		if got := FirstNonASCII(h); got != wantFirst {
			t.Fatalf("FirstNonASCII(%q) = %d, want %d", h, got, wantFirst)
		}
		if got := firstNonASCIIGeneric(h); got != wantFirst {
			t.Fatalf("firstNonASCIIGeneric(%q) = %d, want %d", h, got, wantFirst)
		}
		wantCount := refCountNonASCII(h)
		if got := CountNonASCII(h); got != wantCount {
			t.Fatalf("CountNonASCII(%q) = %d, want %d", h, got, wantCount)
		}
		if got := countNonASCIIGeneric(h); got != wantCount {
			t.Fatalf("countNonASCIIGeneric(%q) = %d, want %d", h, got, wantCount)
		}
	})
}

// FuzzMemchrClass checks the \w classifier, which the AVX2 kernels (VPSHUFB nibble
// tables) and the fallbacks (SWAR range masks) derive independently.
func FuzzMemchrClass(f *testing.F) {
	f.Add([]byte("hello world"))
	f.Add([]byte("________"))
	f.Add([]byte("!\"#$%&'()*+,-./:;<=>?@[\\]^`{|}~ "))
	f.Add(append(bytes.Repeat([]byte{'a'}, 63), '.'))
	f.Add(append(bytes.Repeat([]byte{'.'}, 129), '_'))
	f.Add([]byte{0x2f, 0x30, 0x39, 0x3a, 0x40, 0x41, 0x5a, 0x5b, 0x5e, 0x5f, 0x60, 0x7a, 0x7b, 0x7f, 0x80, 0xff})
	f.Fuzz(func(t *testing.T, h []byte) {
		wantWord, wantNot, wantDigit := -1, -1, -1
		for i, b := range h {
			if wantWord < 0 && refIsWord(b) {
				wantWord = i
			}
			if wantNot < 0 && !refIsWord(b) {
				wantNot = i
			}
			if wantDigit < 0 && b >= '0' && b <= '9' {
				wantDigit = i
			}
		}
		if got := MemchrWord(h); got != wantWord {
			t.Fatalf("MemchrWord(%q) = %d, want %d", h, got, wantWord)
		}
		if got := memchrWordGeneric(h); got != wantWord {
			t.Fatalf("memchrWordGeneric(%q) = %d, want %d", h, got, wantWord)
		}
		if got := MemchrNotWord(h); got != wantNot {
			t.Fatalf("MemchrNotWord(%q) = %d, want %d", h, got, wantNot)
		}
		if got := memchrNotWordGeneric(h); got != wantNot {
			t.Fatalf("memchrNotWordGeneric(%q) = %d, want %d", h, got, wantNot)
		}
		if got := MemchrDigit(h); got != wantDigit {
			t.Fatalf("MemchrDigit(%q) = %d, want %d", h, got, wantDigit)
		}
		if got := memchrDigitGeneric(h); got != wantDigit {
			t.Fatalf("memchrDigitGeneric(%q) = %d, want %d", h, got, wantDigit)
		}
	})
}
