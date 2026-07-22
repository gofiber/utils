package utils

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Scalar references the SWAR implementations are checked against.

// Tight scalar loops mirroring what fiber's byte-at-a-time call sites do
// today — the honest baseline for the SWAR versions.
func refIndexAny2(s string, a, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == a || s[i] == b {
			return i
		}
	}
	return -1
}

func refIndexAny3(s string, a, b, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == a || s[i] == b || s[i] == c {
			return i
		}
	}
	return -1
}

func refIndexFold(s, needle string) int {
	if len(needle) == 0 {
		return 0
	}
	fold := func(c byte) byte {
		if c >= 'A' && c <= 'Z' {
			return c + 'a' - 'A'
		}
		return c
	}
	for i := 0; i+len(needle) <= len(s); i++ {
		ok := true
		for j := 0; j < len(needle); j++ {
			if fold(s[i+j]) != fold(needle[j]) {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func refIndexNonQuotable(s string) int {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c == '\\' || c == '"' || (c < 0x20 && c != '\t') || c == 0x7f {
			return i
		}
	}
	return -1
}

func Test_IndexAny_Fixed(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, IndexAny2("", ',', '"'))
	require.Equal(t, 0, IndexAny2(",", ',', '"'))
	require.Equal(t, 3, IndexAny2([]byte("abc,def"), ',', '"'))
	require.Equal(t, 5, IndexAny2("abcde\"fgh,ij", '"', ','))
	require.Equal(t, -1, IndexAny3("abcdefgh", ',', '"', '\\'))
	require.Equal(t, 7, IndexAny3("abcdefg\\", ',', '"', '\\'))
	// Needle equal to 0x00 and >= 0x80 must work exactly.
	require.Equal(t, 2, IndexAny2("ab\x00cd", 0x00, 0x01))
	require.Equal(t, 4, IndexAny2("abcd\xe9f", 0xE9, 0x01))
}

func Test_Search_Randomized(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(7)) //nolint:gosec // deterministic test data

	for range 30000 {
		// Lengths 0..100 cover all %8 tails, the overlap reads, and 2+
		// iterations of IsASCII's 32-byte unrolled loop.
		n := rng.Intn(101)
		buf := make([]byte, n)
		for i := range buf {
			// Small alphabet so needles hit often, plus occasional raw bytes.
			if rng.Intn(8) == 0 {
				buf[i] = byte(rng.Intn(256))
			} else {
				buf[i] = "abcDEF,\"\\\x00\x7f\xe9\t :."[rng.Intn(15)]
			}
		}
		s := string(buf)

		a, b, c := byte(rng.Intn(256)), byte(rng.Intn(256)), byte(rng.Intn(256))
		require.Equal(t, refIndexAny2(s, a, b), IndexAny2(s, a, b), "IndexAny2 %q %q %q", s, a, b)
		require.Equal(t, refIndexAny2(s, a, b), IndexAny2(buf, a, b))
		require.Equal(t, refIndexAny3(s, a, b, c), IndexAny3(s, a, b, c), "IndexAny3 %q", s)
		require.Equal(t, refIndexAny3(s, a, b, c), IndexAny3(buf, a, b, c))

		require.Equal(t, refIndexNonQuotable(s), IndexNonQuotable(s), "IndexNonQuotable %q", s)
		require.Equal(t, refIndexNonQuotable(s), IndexNonQuotable(buf))
		require.Equal(t, isASCIIRef(s), IsASCII(s), "IsASCII %q", s)
		require.Equal(t, isASCIIRef(s), IsASCII(buf))

		// Needle: random slice of s (guaranteed present) or random bytes.
		var needle string
		if n > 0 && rng.Intn(2) == 0 {
			start := rng.Intn(n)
			end := start + rng.Intn(min(n-start, 12)+1)
			needle = s[start:end]
			if rng.Intn(2) == 0 {
				// Random-case the needle: must still be found.
				nb := []byte(needle)
				for i := range nb {
					if rng.Intn(2) == 0 {
						if nb[i] >= 'a' && nb[i] <= 'z' {
							nb[i] -= 'a' - 'A'
						} else if nb[i] >= 'A' && nb[i] <= 'Z' {
							nb[i] += 'a' - 'A'
						}
					}
				}
				needle = string(nb)
			}
		} else {
			nb := make([]byte, rng.Intn(12))
			for i := range nb {
				nb[i] = byte(rng.Intn(256))
			}
			needle = string(nb)
		}
		require.Equal(t, refIndexFold(s, needle), IndexFold(s, needle), "IndexFold %q %q", s, needle)
		require.Equal(t, refIndexFold(s, needle), IndexFold(buf, needle))
		require.Equal(t, refIndexFold(s, needle) >= 0, ContainsFold(s, needle))
	}
}

func isASCIIRef(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func Test_IndexFold_Fixed(t *testing.T) {
	t.Parallel()
	require.Equal(t, 0, IndexFold("", ""))
	require.Equal(t, 0, IndexFold("abc", ""))
	require.Equal(t, -1, IndexFold("", "a"))
	require.Equal(t, 0, IndexFold("NO-CACHE", "no-cache"))
	require.Equal(t, 10, IndexFold("max-age=0,No-Cache,private", "no-cache"))
	require.Equal(t, 4, IndexFold([]byte("xyz BeArEr token"), "bearer"))
	// Fold must only touch A-Z/a-z lanes: CR|0x20 == '-' would falsely match
	// with a blanket |0x20 fold.
	require.Equal(t, -1, IndexFold("no\rcache", "no-cache"))
	require.Equal(t, -1, IndexFold("NO\rCACHE", "no-cache"))
	// Non-ASCII bytes must match exactly, never fold.
	require.Equal(t, 0, IndexFold("\xc9abc", "\xc9ABC"))
	require.Equal(t, -1, IndexFold("\xe9abc", "\xc9abc"))
	// Needles longer than 8 bytes take the scalar path.
	require.Equal(t, 2, IndexFold("__Proxy-Authorization: basic", "proxy-authorization"))
	require.Equal(t, -1, IndexFold("__Proxy-Authorizatio", "proxy-authorization!"))
	// Match at the very end, crossing the overlapping tail window.
	require.Equal(t, 12, IndexFold("aaaaaaaaaaaaNO-CACHE", "no-cache"))
	require.Equal(t, 5, IndexFold("aaaaaGZip", "gzip"))
	// First-byte candidates that overrun the last valid start must not match.
	require.Equal(t, -1, IndexFold("aaaaaaagz", "gzip"))
	// Same folded word, different case mixes.
	require.Equal(t, 0, IndexFold("gZiP, deflate", "GzIp"))
	require.True(t, ContainsFold("Cache-Control: NO-CACHE", "no-cache"))
	require.False(t, ContainsFold("Cache-Control: no\rcache", "no-cache"))
	// Long needles (> 8 bytes) use the SWAR candidate scan with a word-wise
	// fold-compare verify, including matches ending exactly at the tail.
	require.Equal(t, 21, IndexFold("x-real-ip, x-req-id, Proxy-Authorization", "proxy-authorization"))
	require.Equal(t, 7, IndexFold("cookie CONTENT-type: text", "content-type"))
	require.Equal(t, -1, IndexFold("cookie CONTENT-typ: text.", "content-type"))
	require.Equal(t, 12, IndexFold("abcdefghijklTRANSFER-ENCODING", "transfer-encoding"))
	require.Equal(t, 0, IndexFold("authorizatioN", "authorization"))
	require.Equal(t, -1, IndexFold("authorizatio", "authorization"))
	// The first candidate appears only in the scalar remainder after the
	// word loop, so the folded needle is built there.
	require.Equal(t, 10, IndexFold("aaaaaaaaaaGZ", "gz"))
	require.Equal(t, -1, IndexFold("aaaaaaaaaaGx", "gz"))
	// Needle length an exact multiple of the word size: the fold-compare
	// covers it with full words only, matching and mismatching.
	require.Equal(t, 2, IndexFold("__proxy-authorizatXYZ", "PROXY-AUTHORIZAT"))
	require.Equal(t, -1, IndexFold("__proxy-authorizaQXYZ", "PROXY-AUTHORIZAT"))
}

// refHasPrefixFold/refHasSuffixFold reduce to refIndexFold on the clipped
// window: with len(window) == len(needle) a 0 return means "equal folded".
func refHasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && refIndexFold(s[:len(prefix)], prefix) == 0
}

func refHasSuffixFold(s, suffix string) bool {
	return len(s) >= len(suffix) && refIndexFold(s[len(s)-len(suffix):], suffix) == 0
}

func Test_HasPrefixFold_HasSuffixFold_Fixed(t *testing.T) {
	t.Parallel()
	// Empty needles match anything, oversized needles nothing.
	require.True(t, HasPrefixFold("", ""))
	require.True(t, HasSuffixFold("", ""))
	require.True(t, HasPrefixFold("x", ""))
	require.True(t, HasSuffixFold("x", ""))
	require.False(t, HasPrefixFold("", "a"))
	require.False(t, HasSuffixFold("", "a"))
	require.False(t, HasPrefixFold("bear", "bearer "))
	require.False(t, HasSuffixFold("zip", "gzip"))

	// Scalar path (needle < 8 bytes), both hit and miss, string and []byte.
	require.True(t, HasPrefixFold("Bearer abc.def", "bearer "))
	require.True(t, HasPrefixFold([]byte("BEARER abc.def"), "bearer "))
	require.False(t, HasPrefixFold("Basic QWxhZGRpbg==", "bearer "))
	require.True(t, HasSuffixFold("application/vnd.api+JSON", "+json"))
	require.True(t, HasSuffixFold([]byte("deflate, GZIP"), "gzip"))
	require.False(t, HasSuffixFold("deflate, GZIPx", "gzip"))

	// Word path (needle >= 8 bytes), including the overlapping tail compare.
	require.True(t, HasPrefixFold("MULTIPART/FORM-DATA; boundary=X", "multipart/form-data"))
	require.False(t, HasPrefixFold("multipart/form-datX; boundary=X", "multipart/form-data"))
	require.True(t, HasPrefixFold("multipart/form-data", "multipart/form-data"))
	require.True(t, HasSuffixFold("application/ATOM+XML", "atom+xml"))
	require.False(t, HasSuffixFold("application/atomxxml", "atom+xml"))
	require.True(t, HasSuffixFold("PROXY-AUTHORIZATION", "proxy-authorization"))

	// Exact-length inputs on the scalar path.
	require.True(t, HasPrefixFold("GZip", "gzip"))
	require.True(t, HasSuffixFold("GZip", "gzip"))

	// Only A-Z/a-z fold: CR|0x20 == '-' style false folds must not match,
	// and non-ASCII bytes must match exactly.
	require.False(t, HasPrefixFold("no\rcache", "no-cache"))
	require.False(t, HasSuffixFold("xno\rcache", "no-cache"))
	require.True(t, HasPrefixFold("\xc9abcdefgh", "\xc9ABCDEFGH"))
	require.False(t, HasPrefixFold("\xe9abcdefgh", "\xc9abcdefgh"))
	require.True(t, HasSuffixFold("x\xc9abcdefg", "\xc9ABCDEFG"))
	require.False(t, HasSuffixFold("x\xe9abcdefg", "\xc9abcdefg"))
}

func Test_PrefixSuffixFold_Randomized(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(13)) //nolint:gosec // deterministic test data

	for range 30000 {
		n := rng.Intn(41)
		buf := make([]byte, n)
		for i := range buf {
			if rng.Intn(8) == 0 {
				buf[i] = byte(rng.Intn(256))
			} else {
				buf[i] = "abcDEF,\"\\\x00\x7f\xe9\t :."[rng.Intn(15)]
			}
		}
		s := string(buf)

		// Needle: prefix/suffix of s (random-cased, guaranteed hit) or
		// random bytes (usually a miss), lengths crossing the 8-byte split.
		k := rng.Intn(min(n, 12) + 1)
		var needle []byte
		switch rng.Intn(3) {
		case 0:
			needle = []byte(s[:k])
		case 1:
			needle = []byte(s[n-k:])
		default:
			needle = make([]byte, rng.Intn(12))
			for i := range needle {
				needle[i] = byte(rng.Intn(256))
			}
		}
		for i := range needle {
			if rng.Intn(2) == 0 {
				if needle[i] >= 'a' && needle[i] <= 'z' {
					needle[i] -= 'a' - 'A'
				} else if needle[i] >= 'A' && needle[i] <= 'Z' {
					needle[i] += 'a' - 'A'
				}
			}
		}
		nd := string(needle)

		require.Equal(t, refHasPrefixFold(s, nd), HasPrefixFold(s, nd), "HasPrefixFold %q %q", s, nd)
		require.Equal(t, refHasPrefixFold(s, nd), HasPrefixFold(buf, nd))
		require.Equal(t, refHasSuffixFold(s, nd), HasSuffixFold(s, nd), "HasSuffixFold %q %q", s, nd)
		require.Equal(t, refHasSuffixFold(s, nd), HasSuffixFold(buf, nd))
	}
}

func Test_IndexNonQuotable_Fixed(t *testing.T) {
	t.Parallel()
	require.Equal(t, -1, IndexNonQuotable(""))
	require.Equal(t, -1, IndexNonQuotable("simple filename.txt"))
	require.Equal(t, 4, IndexNonQuotable(`name"quoted`))
	require.Equal(t, 4, IndexNonQuotable(`name\quoted`))
	require.Equal(t, 0, IndexNonQuotable("\x1ftail"))
	require.Equal(t, 7, IndexNonQuotable("abcdefg\x7f"))
	// obs-text (>= 0x80) is quotable.
	require.Equal(t, -1, IndexNonQuotable("caf\xc3\xa9 r\xc3\xa9sum\xc3\xa9.pdf"))
	// HTAB is valid qdtext (RFC 9110) and must NOT be flagged, in the word
	// path, the scalar path, and directly before a real control byte.
	require.Equal(t, -1, IndexNonQuotable([]byte("aaaaaaaaa\ttab")))
	require.Equal(t, -1, IndexNonQuotable("a\tb"))
	require.Equal(t, 9, IndexNonQuotable("aaaaaaaa\t\x1f"))
	require.Equal(t, 3, IndexNonQuotable("ab\t\x08cdefgh"))
	// Borrow-stress shapes: obs-text lanes before a quote, and a quote
	// directly ahead of a control byte.
	require.Equal(t, 7, IndexNonQuotable("\xff\xff\xff\xff\xff\xff\xff\""))
	require.Equal(t, 1, IndexNonQuotable("a\"\x00defgh"))
}

func Test_IsASCII_Fixed(t *testing.T) {
	t.Parallel()
	require.True(t, IsASCII(""))
	require.True(t, IsASCII("plain ascii filename.txt"))
	require.False(t, IsASCII("caf\xc3\xa9"))
	require.False(t, IsASCII([]byte{0x80}))
	require.True(t, IsASCII([]byte{0x7f, 0x00, 'a'}))
	require.False(t, IsASCII("aaaaaaaaaaaaaaa\xff")) // high byte only in the tail
}

func Test_Search_NamedTypes_SIMDDispatch(t *testing.T) {
	t.Parallel()
	// Named byteSeq instantiations at 32+ bytes take the amd64 simd dispatch,
	// whose zero-copy []byte view assumes the header layouts named types
	// share with plain string/[]byte.
	type headerString string
	type headerBytes []byte
	long := strings.Repeat("x", 40) + ","
	require.Equal(t, 40, IndexAny2(headerString(long), ',', '"'))
	require.Equal(t, 40, IndexAny2(headerBytes(long), ',', '"'))
	require.Equal(t, 40, IndexAny3(headerString(long), ',', '"', '\\'))
	require.Equal(t, 40, IndexAny3(headerBytes(long), ',', '"', '\\'))
	require.True(t, IsASCII(headerString(long)))
	require.True(t, IsASCII(headerBytes(long)))
	require.False(t, IsASCII(headerString(long+"\xe9")))
	require.False(t, IsASCII(headerBytes(long+"\xe9")))
}

// --- Benchmarks: SWAR vs scalar reference at HTTP-typical sizes ---

var benchSizes = []int{7, 8, 16, 32, 64, 512}

func benchInput(size, needleAt int, needle byte) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = 'a' + byte(i%23)
	}
	if needleAt >= 0 && needleAt < size {
		b[needleAt] = needle
	}
	return string(b)
}

func Benchmark_IndexAny2(b *testing.B) {
	for _, size := range benchSizes {
		s := benchInput(size, size-1, ',') // worst case: match at the very end
		b.Run(FormatInt(int64(size))+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = IndexAny2(s, ',', '"')
			}
			_ = r
		})
		b.Run(FormatInt(int64(size))+"B/scalar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = refIndexAny2(s, ',', '"')
			}
			_ = r
		})
		b.Run(FormatInt(int64(size))+"B/stdlib-indexany", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = strings.IndexAny(s, ",\"")
			}
			_ = r
		})
	}
}

func Benchmark_IndexAny3(b *testing.B) {
	for _, size := range benchSizes {
		s := benchInput(size, size-1, ';')
		b.Run(FormatInt(int64(size))+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = IndexAny3(s, ';', '"', '\\')
			}
			_ = r
		})
		b.Run(FormatInt(int64(size))+"B/scalar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = refIndexAny3(s, ';', '"', '\\')
			}
			_ = r
		})
	}
}

func Benchmark_IndexFold(b *testing.B) {
	headers := map[string]string{
		"miss-8B":   "max-age=",
		"miss-32B":  "max-age=31536000, must-revalidat",
		"hit-32B":   "max-age=0, private, No-CaChE, mu",
		"miss-64B":  strings.Repeat("max-age=31536000", 4),
		"hit-tail":  "max-age=31536000, immutable, private, stale-if-error, NO-CACHE",
		"miss-512B": strings.Repeat("max-age=31536000", 32),
	}
	longNeedle := map[string]string{
		"long-miss-64B": strings.Repeat("max-age=31536000", 4),
		"long-hit-64B":  "max-age=31536000, immutable, Proxy-Authorization: basic aaaa",
	}
	for name, s := range headers {
		b.Run(name+"/swar", func(b *testing.B) {
			b.SetBytes(int64(len(s)))
			var r int
			for b.Loop() {
				r = IndexFold(s, "no-cache")
			}
			_ = r
		})
		b.Run(name+"/scalar", func(b *testing.B) {
			b.SetBytes(int64(len(s)))
			var r int
			for b.Loop() {
				r = refIndexFold(s, "no-cache")
			}
			_ = r
		})
	}
	for name, s := range longNeedle {
		b.Run(name+"/swar", func(b *testing.B) {
			b.SetBytes(int64(len(s)))
			var r int
			for b.Loop() {
				r = IndexFold(s, "proxy-authorization")
			}
			_ = r
		})
		b.Run(name+"/scalar", func(b *testing.B) {
			b.SetBytes(int64(len(s)))
			var r int
			for b.Loop() {
				r = refIndexFold(s, "proxy-authorization")
			}
			_ = r
		})
	}
}

func Benchmark_HasPrefixFold(b *testing.B) {
	cases := []struct {
		name   string
		s      string
		prefix string
	}{
		{"bearer-hit", "Bearer abc.def.ghi.jkl.mno.pqr", "bearer "},
		{"bearer-miss", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "bearer "},
		{"text-hit", "TEXT/HTML; charset=utf-8", "text/"},
		{"multipart-hit", "MULTIPART/FORM-DATA; boundary=X-BOUNDARY", "multipart/form-data"},
		{"multipart-miss", "application/x-www-form-urlencoded", "multipart/form-data"},
	}
	for _, tc := range cases {
		b.Run(tc.name+"/fiber", func(b *testing.B) {
			b.SetBytes(int64(len(tc.prefix)))
			var r bool
			for b.Loop() {
				r = HasPrefixFold(tc.s, tc.prefix)
			}
			_ = r
		})
		b.Run(tc.name+"/default", func(b *testing.B) {
			b.SetBytes(int64(len(tc.prefix)))
			var r bool
			for b.Loop() {
				r = len(tc.s) >= len(tc.prefix) && strings.EqualFold(tc.s[:len(tc.prefix)], tc.prefix)
			}
			_ = r
		})
	}
}

func Benchmark_HasSuffixFold(b *testing.B) {
	cases := []struct {
		name   string
		s      string
		suffix string
	}{
		{"gzip-hit", "deflate, br, GZIP", "gzip"},
		{"gzip-miss", "deflate, br, zstd", "gzip"},
		{"json-hit", "application/vnd.api+JSON", "+json"},
		{"atomxml-hit", "application/ATOM+XML", "atom+xml"},
		{"charset-miss", "text/html; charset=utf-8", "charset=iso-8859-1"},
	}
	for _, tc := range cases {
		b.Run(tc.name+"/fiber", func(b *testing.B) {
			b.SetBytes(int64(len(tc.suffix)))
			var r bool
			for b.Loop() {
				r = HasSuffixFold(tc.s, tc.suffix)
			}
			_ = r
		})
		b.Run(tc.name+"/default", func(b *testing.B) {
			b.SetBytes(int64(len(tc.suffix)))
			var r bool
			for b.Loop() {
				r = len(tc.s) >= len(tc.suffix) && strings.EqualFold(tc.s[len(tc.s)-len(tc.suffix):], tc.suffix)
			}
			_ = r
		})
	}
}

func Benchmark_IsASCII(b *testing.B) {
	for _, size := range benchSizes {
		s := benchInput(size, -1, 0)
		b.Run(FormatInt(int64(size))+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r bool
			for b.Loop() {
				r = IsASCII(s)
			}
			_ = r
		})
		b.Run(FormatInt(int64(size))+"B/scalar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r bool
			for b.Loop() {
				r = isASCIIRef(s)
			}
			_ = r
		})
	}
}

func Benchmark_IndexNonQuotable(b *testing.B) {
	for _, size := range benchSizes {
		s := benchInput(size, -1, 0) // clean input: full scan
		b.Run(FormatInt(int64(size))+"B/swar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = IndexNonQuotable(s)
			}
			_ = r
		})
		b.Run(FormatInt(int64(size))+"B/scalar", func(b *testing.B) {
			b.SetBytes(int64(size))
			var r int
			for b.Loop() {
				r = refIndexNonQuotable(s)
			}
			_ = r
		})
	}
}
