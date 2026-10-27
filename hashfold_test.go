package utils

import (
	"hash/maphash"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_HashFold pins the contract: keys EqualFold reports equal hash equal, at every length across the word
// boundaries; keys it tells apart (non-letter case pairs, bytes past ASCII, length alone) hash apart.
func Test_HashFold(t *testing.T) {
	t.Parallel()

	for _, pair := range [][2]string{
		{"", ""},
		{"a", "A"},
		{"id", "ID"},
		{"Content-Type", "content-type"},
		{"X-Forwarded-For", "x-FORWARDED-for"},
		{"hobby", "HoBbY"},
		{"nested.field.name", "NESTED.Field.NAME"},
		{"caf\xc3\xa9", "CAF\xc3\xa9"},
	} {
		require.True(t, EqualFold(pair[0], pair[1]))
		require.Equal(t, HashFold(pair[0]), HashFold(pair[1]), "%q and %q", pair[0], pair[1])
		require.Equal(t, HashFold(pair[0]), HashFold([]byte(pair[1])), "%q as bytes", pair[1])
	}

	for _, pair := range [][2]string{
		{"a", "b"},
		{"ab", "ba"},
		{"a", "a\x00"},
		{"", "\x00"},
		{"@", "`"},
		{"[", "{"},
		{"caf\xc3\xa9", "caf\xc3\x89"},
		{"abcd", "abce"},
		{"abcdefg", "abcdefh"},
		{"abcdefgh", "abcdefgi"},
		{"abcdefghi", "abcdefghj"},
		{"abcdefghijklmnop", "abcdefghijklmnoq"},
		// keys packing into the same words, told apart by length alone
		{"b", "ab"},
		{"b", "bb"},
		{"bb", "bbb"},
		{"aaaa", "aaaaa"},
		{"````", "a````"},
		{"aaaaaaaa", "aaaaaaaaa"},
		{"~bcdefghijklmno", "abcdefghijklmno\x00"},
	} {
		require.False(t, EqualFold(pair[0], pair[1]))
		require.NotEqual(t, HashFold(pair[0]), HashFold(pair[1]), "%q and %q", pair[0], pair[1])
	}

	for _, c := range "a-\x00" {
		seen := make(map[uint64]int)
		for n := range 64 {
			h := HashFold(strings.Repeat(string(c), n))
			other, ok := seen[h]
			require.False(t, ok, "%q repeated %d and %d times share hash %#x", c, other, n, h)
			seen[h] = n
		}
	}

	for n := range 40 {
		key := []byte(strings.Repeat("aB3-", n/4+1)[:n])
		want := HashFold(string(key))
		for i := range key {
			flipped := append([]byte(nil), key...)
			if c := flipped[i]; c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				flipped[i] ^= 0x20
			}
			require.Equal(t, want, HashFold(flipped), "length %d, byte %d", n, i)
		}
	}
}

func Test_HashFold_NoCollisions(t *testing.T) {
	t.Parallel()

	seen := make(map[uint64]string)
	add := func(key string) {
		h := HashFold(key)
		if other, ok := seen[h]; ok && !EqualFold(other, key) {
			t.Fatalf("%q and %q share hash %#x", other, key, h)
		}
		seen[h] = key
	}
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789-_."
	add("")
	for _, a := range alphabet {
		add(string(a))
		for _, b := range alphabet {
			add(string(a) + string(b))
			for _, c := range alphabet[:16] {
				add(string(a) + string(b) + string(c))
			}
		}
	}
	for i := range 40000 {
		add("field" + strconv.Itoa(i))
		add("X-Header-" + strconv.Itoa(i) + "-Name")
		if i%4 == 0 {
			add("user.addresses." + strconv.Itoa(i) + ".street.line.number")
		}
	}
	require.Greater(t, len(seen), 110000)
}

func FuzzHashFold(f *testing.F) {
	f.Add("Content-Type", "content-type")
	f.Add("abcdefghi", "ABCDEFGHI")
	f.Add("a", "b")
	f.Add("@[`{", "`{@[")
	f.Add("\xc9abc", "\xe9abc")
	f.Fuzz(func(t *testing.T, a, b string) {
		ha := HashFold(a)
		if HashFold([]byte(a)) != ha {
			t.Fatalf("HashFold(%q) differs between string and bytes", a)
		}
		if EqualFold(a, b) && HashFold(b) != ha {
			t.Fatalf("HashFold(%q) != HashFold(%q) though EqualFold holds", a, b)
		}
		swapped := []byte(a)
		for i, c := range swapped {
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				swapped[i] = c ^ 0x20
			}
		}
		if HashFold(swapped) != ha {
			t.Fatalf("HashFold(%q) != HashFold(%q)", a, swapped)
		}
	})
}

func Benchmark_HashFold(b *testing.B) {
	seed := maphash.MakeSeed()
	for _, key := range []string{"id", "Name", "Content-Type", "X-Forwarded-Proto", strings.Repeat("Header-Name-", 4)} {
		b.Run(strconv.Itoa(len(key))+"B/fiber", func(b *testing.B) {
			var sink uint64
			for b.Loop() {
				sink += HashFold(key)
			}
			_ = sink
		})
		b.Run(strconv.Itoa(len(key))+"B/default", func(b *testing.B) {
			// Baseline: fold a copy, then hash.
			var sink uint64
			for b.Loop() {
				sink += maphash.String(seed, strings.ToLower(key))
			}
			_ = sink
		})
	}
}
