package utils

import (
	"bytes"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	casebytes "github.com/gofiber/utils/v2/bytes"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	casestrings "github.com/gofiber/utils/v2/strings"
)

func forEachString(alphabet string, maxLen int, fn func(s string)) {
	buf := make([]byte, 0, maxLen)
	var walk func(depth int)
	walk = func(depth int) {
		fn(string(buf))
		if depth == 0 {
			return
		}
		for i := range len(alphabet) {
			buf = append(buf, alphabet[i])
			walk(depth - 1)
			buf = buf[:len(buf)-1]
		}
	}
	walk(maxLen)
}

// requireSameView fails unless got is want as the stdlib returns it: same bytes, nil-ness and window.
func requireSameView(tb testing.TB, what string, in, got, want []byte) {
	tb.Helper()
	switch {
	case (got == nil) != (want == nil):
		tb.Fatalf("%s(%q): result nil = %v, stdlib nil = %v", what, in, got == nil, want == nil)
	case !bytes.Equal(got, want):
		tb.Fatalf("%s(%q) = %q, stdlib gives %q", what, in, got, want)
	case cap(got) != cap(want):
		tb.Fatalf("%s(%q): capacity %d, stdlib's %d", what, in, cap(got), cap(want))
	case len(got) > 0 && &got[0] != &want[0]:
		tb.Fatalf("%s(%q): not the stdlib's window of the input", what, in)
	}
}

// requireTrimmed is requireSameView except for an empty result: bytes.Trim, TrimLeft and TrimSpace
// return nil there (undocumented), these return an empty slice of the input, keeping TrimSpace inlinable.
func requireTrimmed(tb testing.TB, what string, in, got, want []byte) {
	tb.Helper()
	if len(want) > 0 {
		requireSameView(tb, what, in, got, want)
		return
	}
	if len(got) != 0 {
		tb.Fatalf("%s(%q) = %q, stdlib gives nothing", what, in, got)
	}
	if (got == nil) != (in == nil) {
		tb.Fatalf("%s(%q): an empty result must be a slice of the input, so nil only for nil; got nil = %v for input nil = %v", what, in, got == nil, in == nil)
	}
}

func refTrimLeft(s string, c byte) string {
	for len(s) > 0 && s[0] == c {
		s = s[1:]
	}
	return s
}

func refTrimRight(s string, c byte) string {
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

func Test_Trim_StdlibParity(t *testing.T) {
	t.Parallel()
	for cut := range 256 {
		c := byte(cut)
		cutset := string([]byte{c})
		if c < utf8.RuneSelf {
			// Empty and nil inputs: bytes.TrimRight returns its input, the others nil.
			for _, in := range [][]byte{{}, nil} {
				requireTrimmed(t, "TrimLeft", in, TrimLeft(in, c), bytes.TrimLeft(in, cutset))
				requireSameView(t, "TrimRight", in, TrimRight(in, c), bytes.TrimRight(in, cutset))
				requireTrimmed(t, "Trim", in, Trim(in, c), bytes.Trim(in, cutset))
			}
		}
		forEachString(string([]byte{c, 'a', ' ', 0xe9}), 5, func(s string) {
			wantLeft, wantRight, wantBoth := refTrimLeft(s, c), refTrimRight(s, c), refTrimRight(refTrimLeft(s, c), c)
			if c < utf8.RuneSelf {
				// Below 0x80 the stdlib's cutset is the byte; above it, malformed UTF-8 matching every malformed byte.
				if strings.TrimLeft(s, cutset) != wantLeft || strings.TrimRight(s, cutset) != wantRight || strings.Trim(s, cutset) != wantBoth {
					t.Fatalf("reference trims disagree with the stdlib for %q cutset %q", s, cutset)
				}
				b := []byte(s)
				requireTrimmed(t, "TrimLeft", b, TrimLeft(b, c), bytes.TrimLeft(b, cutset))
				requireSameView(t, "TrimRight", b, TrimRight(b, c), bytes.TrimRight(b, cutset))
				requireTrimmed(t, "Trim", b, Trim(b, c), bytes.Trim(b, cutset))
			}
			if got := TrimLeft(s, c); got != wantLeft {
				t.Fatalf("TrimLeft(%q, %#x) = %q, want %q", s, c, got, wantLeft)
			}
			if got := TrimRight(s, c); got != wantRight {
				t.Fatalf("TrimRight(%q, %#x) = %q, want %q", s, c, got, wantRight)
			}
			if got := Trim(s, c); got != wantBoth {
				t.Fatalf("Trim(%q, %#x) = %q, want %q", s, c, got, wantBoth)
			}
			b := []byte(s)
			if got := TrimLeft(b, c); string(got) != wantLeft || (len(got) > 0 && &got[0] != &b[len(b)-len(wantLeft)]) {
				t.Fatalf("TrimLeft(bytes %q, %#x) = %q, want the tail %q", s, c, got, wantLeft)
			}
			if got := TrimRight(b, c); string(got) != wantRight || cap(got) != cap(b) {
				t.Fatalf("TrimRight(bytes %q, %#x) = %q (cap %d), want the head %q", s, c, got, cap(got), wantRight)
			}
			if got := Trim(b, c); string(got) != wantBoth {
				t.Fatalf("Trim(bytes %q, %#x) = %q, want %q", s, c, got, wantBoth)
			}
		})
	}
}

func Test_TrimSpace_StdlibParity(t *testing.T) {
	t.Parallel()
	// Whitespace and its neighbors in a table or range test.
	forEachString("\t\n\v\f\r \x00\x08\x0e\x1f\x7f!a", 5, func(s string) {
		if got, want := TrimSpace(s), strings.TrimSpace(s); got != want {
			t.Fatalf("TrimSpace(%q) = %q, stdlib gives %q", s, got, want)
		}
		b := []byte(s)
		requireTrimmed(t, "TrimSpace", b, TrimSpace(b), bytes.TrimSpace(b))
	})
	requireTrimmed(t, "TrimSpace(empty)", []byte{}, TrimSpace([]byte{}), bytes.TrimSpace([]byte{}))
	requireTrimmed(t, "TrimSpace(nil)", nil, TrimSpace([]byte(nil)), bytes.TrimSpace(nil))

	type word string
	type blob []byte
	if got := TrimSpace(word(" \tx\n")); got != "x" {
		t.Fatalf("TrimSpace(word) = %q", got)
	}
	if got := TrimSpace(blob(" \tx\n")); string(got) != "x" {
		t.Fatalf("TrimSpace(blob) = %q", got)
	}
}

// Test_TrimSpace_OnlyASCIIWhitespace pins the deliberate difference: only the six ASCII whitespace
// bytes are trimmed, not Unicode spaces such as U+00A0.
func Test_TrimSpace_OnlyASCIIWhitespace(t *testing.T) {
	t.Parallel()
	for r := rune(utf8.RuneSelf); r <= unicode.MaxRune; r++ {
		if !unicode.IsSpace(r) {
			continue
		}
		in := " " + string(r) + "x" + string(r) + " "
		want := string(r) + "x" + string(r)
		if got := TrimSpace(in); got != want {
			t.Fatalf("TrimSpace(%q) = %q, want %q", in, got, want)
		}
		if got := TrimSpace([]byte(in)); string(got) != want {
			t.Fatalf("TrimSpace(bytes %q) = %q, want %q", in, got, want)
		}
		if strings.TrimSpace(in) != "x" {
			t.Fatalf("the stdlib no longer trims U+%04X: the documented difference is gone", r)
		}
	}
	// A lone byte above 0x7F is no space to either (0x85, 0xA0 too): the stdlib reads it as malformed UTF-8.
	for b := utf8.RuneSelf; b < 256; b++ {
		for _, in := range []string{string([]byte{byte(b)}), " " + string([]byte{byte(b)}) + " ", "x" + string([]byte{byte(b)})} {
			if got, want := TrimSpace(in), strings.TrimSpace(in); got != want {
				t.Fatalf("TrimSpace(%q) = %q, stdlib gives %q", in, got, want)
			}
		}
	}
	if got := TrimSpace("\u00a0 x \u00a0"); got != "\u00a0 x \u00a0" {
		t.Fatalf("TrimSpace left %q", got)
	}
}

func Test_CutByte_StdlibParity(t *testing.T) {
	t.Parallel()
	for sep := range 256 {
		c := byte(sep)
		needle := string([]byte{c}) // a raw byte: strings.Cut compares bytes, never runes
		forEachString(string([]byte{c, 'a', 0xe9}), 6, func(s string) {
			before, after, found := CutByte(s, c)
			wantBefore, wantAfter, wantFound := strings.Cut(s, needle)
			if before != wantBefore || after != wantAfter || found != wantFound {
				t.Fatalf("CutByte(%q, %#x) = (%q, %q, %v), strings.Cut gives (%q, %q, %v)", s, c, before, after, found, wantBefore, wantAfter, wantFound)
			}

			b := []byte(s)
			bb, ba, bfound := CutByte(b, c)
			wbb, wba, wfound := bytes.Cut(b, []byte(needle))
			if bfound != wfound {
				t.Fatalf("CutByte(bytes %q, %#x) found = %v, bytes.Cut gives %v", s, c, bfound, wfound)
			}
			requireSameView(t, "CutByte before", b, bb, wbb)
			requireSameView(t, "CutByte after", b, ba, wba)

			// LastCutByte has no stdlib twin: it is Cut around LastIndex.
			lbefore, lafter, lfound := LastCutByte(s, c)
			if i := strings.LastIndex(s, needle); i >= 0 {
				if lbefore != s[:i] || lafter != s[i+1:] || !lfound {
					t.Fatalf("LastCutByte(%q, %#x) = (%q, %q, %v)", s, c, lbefore, lafter, lfound)
				}
			} else if lbefore != s || lafter != "" || lfound {
				t.Fatalf("LastCutByte(%q, %#x) miss = (%q, %q, %v)", s, c, lbefore, lafter, lfound)
			}
			lbb, lba, lbfound := LastCutByte(b, c)
			if i := bytes.LastIndex(b, []byte(needle)); i >= 0 {
				requireSameView(t, "LastCutByte before", b, lbb, b[:i])
				requireSameView(t, "LastCutByte after", b, lba, b[i+1:])
				if !lbfound {
					t.Fatalf("LastCutByte(bytes %q, %#x) found = false", s, c)
				}
			} else {
				// A miss is the whole input and a nil after-part, as bytes.Cut's.
				requireSameView(t, "LastCutByte miss before", b, lbb, b)
				requireSameView(t, "LastCutByte miss after", b, lba, nil)
				if lbfound {
					t.Fatalf("LastCutByte(bytes %q, %#x) found = true on a miss", s, c)
				}
			}
		})
	}
}

func Test_EqualFold_StdlibParity(t *testing.T) {
	t.Parallel()
	const edges = "aAzZ@[`{1"
	forEachString(edges, 3, func(a string) {
		forEachString(edges, 3, func(b string) {
			want := strings.EqualFold(a, b)
			if got := EqualFold(a, b); got != want {
				t.Fatalf("EqualFold(%q, %q) = %v, stdlib gives %v", a, b, got, want)
			}
			if got := EqualFold([]byte(a), []byte(b)); got != want {
				t.Fatalf("EqualFold(bytes %q, %q) = %v, stdlib gives %v", a, b, got, want)
			}
		})
	})

	for n := 1; n <= 40; n++ {
		base := []byte(strings.Repeat("aBcDeFgH", 6)[:n])
		for pos := range n {
			for v := range 256 {
				other := bytes.Clone(base)
				other[pos] = byte(v)
				want := asciiFoldLower(string(base)) == asciiFoldLower(string(other))
				if v < utf8.RuneSelf {
					if std := bytes.EqualFold(base, other); std != want {
						t.Fatalf("the stdlib disagrees with the ASCII fold on %q and %q", base, other)
					}
				}
				if got := EqualFold(base, other); got != want {
					t.Fatalf("EqualFold(%q, %q) = %v, want %v", base, other, got, want)
				}
				if got := EqualFold(string(other), string(base)); got != want {
					t.Fatalf("EqualFold(%q, %q) = %v, want %v", other, base, got, want)
				}
			}
		}
	}
}

// Test_EqualFold_OnlyASCIIFolds pins the deliberate difference: the stdlib's Unicode folding relates letters of different bytes.
func Test_EqualFold_OnlyASCIIFolds(t *testing.T) {
	t.Parallel()
	for _, tc := range [][2]string{
		{"k", "\u212a"}, // KELVIN SIGN folds to k
		{"s", "\u017f"}, // LONG S folds to s
		{"\u00e9", "\u00c9"},
		{"\u03c3", "\u03c2"},
		{"\u01c5", "\u01c6"},
	} {
		if !strings.EqualFold(tc[0], tc[1]) {
			t.Fatalf("the stdlib no longer folds %q and %q: the documented difference is gone", tc[0], tc[1])
		}
		if EqualFold(tc[0], tc[1]) || EqualFold([]byte(tc[0]), []byte(tc[1])) {
			t.Fatalf("EqualFold(%q, %q) = true, want ASCII-only folding", tc[0], tc[1])
		}
		if !EqualFold(tc[0], tc[0]) {
			t.Fatalf("EqualFold(%q, %q) = false for identical bytes", tc[0], tc[0])
		}
	}
}

func Test_IndexAny_StdlibParity(t *testing.T) {
	t.Parallel()
	needles := [][3]byte{
		{'a', 'b', 'c'},
		{0x00, 0x01, 0x02},
		{0x7f, 0x80, 0x81},
		{0xff, 0xfe, 0x7e},
		{'"', '\\', '\t'},
		{0x80, 0x80, 0x80},
	}
	for _, nd := range needles {
		a, b, c := nd[0], nd[1], nd[2]
		// Fillers near the needles in value or bit pattern, which an inexact SWAR zero-lane test confuses.
		fillers := []byte{0x00, 0x7f, 0x80, 0xff, a ^ 0x80, a + 1, a - 1, b ^ 1, c + 1, 'x'}
		for n := 0; n <= 70; n++ {
			for _, f := range fillers {
				base := bytes.Repeat([]byte{f}, n)
				requireIndexAny(t, base, a, b, c)
				for pos := range n {
					for _, planted := range []byte{a, b, c} {
						h := bytes.Clone(base)
						h[pos] = planted
						requireIndexAny(t, h, a, b, c)
						for later := pos + 1; later < min(n, pos+3); later++ {
							h2 := bytes.Clone(h)
							h2[later] = a
							requireIndexAny(t, h2, a, b, c)
						}
					}
				}
			}
		}
	}
}

// requireIndexAny holds IndexAny2/3 to the earliest bytes.IndexByte hit, and ASCII needles to the stdlib IndexAny.
func requireIndexAny(tb testing.TB, h []byte, a, b, c byte) {
	tb.Helper()
	earliest := func(idx ...int) int {
		best := -1
		for _, i := range idx {
			if i >= 0 && (best < 0 || i < best) {
				best = i
			}
		}
		return best
	}
	hs := string(h)
	want2 := earliest(bytes.IndexByte(h, a), bytes.IndexByte(h, b))
	want3 := earliest(want2, bytes.IndexByte(h, c))
	if a < utf8.RuneSelf && b < utf8.RuneSelf && c < utf8.RuneSelf {
		if std := bytes.IndexAny(h, string([]byte{a, b})); std != want2 {
			tb.Fatalf("bytes.IndexAny(%q, %q%q) = %d, but the earliest IndexByte is %d", h, a, b, std, want2)
		}
		if std := strings.IndexAny(hs, string([]byte{a, b, c})); std != want3 {
			tb.Fatalf("strings.IndexAny(%q, %q%q%q) = %d, but the earliest IndexByte is %d", h, a, b, c, std, want3)
		}
	}
	if got := IndexAny2(h, a, b); got != want2 {
		tb.Fatalf("IndexAny2(len %d, %#x, %#x) = %d, want %d", len(h), a, b, got, want2)
	}
	if got := IndexAny2(hs, a, b); got != want2 {
		tb.Fatalf("IndexAny2(string len %d, %#x, %#x) = %d, want %d", len(h), a, b, got, want2)
	}
	if got := IndexAny3(h, a, b, c); got != want3 {
		tb.Fatalf("IndexAny3(len %d, %#x, %#x, %#x) = %d, want %d", len(h), a, b, c, got, want3)
	}
	if got := IndexAny3(hs, a, b, c); got != want3 {
		tb.Fatalf("IndexAny3(string len %d, %#x, %#x, %#x) = %d, want %d", len(h), a, b, c, got, want3)
	}
}

// deprecated holds the wrappers under test, called through it because their `//go:fix inline`
// makes vet flag direct calls.
var deprecated = struct {
	toLower, toUpper           func(string) string
	toLowerBytes, toUpperBytes func([]byte) []byte
}{ToLower, ToUpper, ToLowerBytes, ToUpperBytes}

func Test_Case_StdlibParity(t *testing.T) {
	t.Parallel()
	check := func(s string) {
		t.Helper()
		ascii := utf8.ValidString(s) && IsASCII(s)
		wantLower, wantUpper := asciiFoldLower(s), asciiFoldUpper(s)
		if ascii && (strings.ToLower(s) != wantLower || strings.ToUpper(s) != wantUpper) {
			t.Fatalf("the ASCII fold disagrees with the stdlib on ASCII input %q", s)
		}
		if got := casestrings.ToLower(s); got != wantLower {
			t.Fatalf("strings.ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := casestrings.ToUpper(s); got != wantUpper {
			t.Fatalf("strings.ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}
		if got := deprecated.toLower(s); got != wantLower {
			t.Fatalf("ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := deprecated.toUpper(s); got != wantUpper {
			t.Fatalf("ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}

		b := []byte(s)
		if got := casebytes.ToLower(b); string(got) != wantLower {
			t.Fatalf("bytes.ToLower(%q) = %q, want %q", s, got, wantLower)
		}
		if got := casebytes.ToUpper(b); string(got) != wantUpper {
			t.Fatalf("bytes.ToUpper(%q) = %q, want %q", s, got, wantUpper)
		}
		if string(b) != s {
			t.Fatalf("bytes.ToLower/ToUpper modified their input: %q became %q", s, b)
		}

		lower := []byte(s)
		if got := casebytes.UnsafeToLower(lower); string(got) != wantLower || string(lower) != wantLower || (len(got) > 0 && &got[0] != &lower[0]) {
			t.Fatalf("bytes.UnsafeToLower(%q) = %q, input %q, want %q in place", s, got, lower, wantLower)
		}
		upper := []byte(s)
		if got := casebytes.UnsafeToUpper(upper); string(got) != wantUpper || string(upper) != wantUpper || (len(got) > 0 && &got[0] != &upper[0]) {
			t.Fatalf("bytes.UnsafeToUpper(%q) = %q, input %q, want %q in place", s, got, upper, wantUpper)
		}
		lowerB, upperB := []byte(s), []byte(s)
		if got := deprecated.toLowerBytes(lowerB); string(got) != wantLower || string(lowerB) != wantLower {
			t.Fatalf("ToLowerBytes(%q) = %q, input %q, want %q in place", s, got, lowerB, wantLower)
		}
		if got := deprecated.toUpperBytes(upperB); string(got) != wantUpper || string(upperB) != wantUpper {
			t.Fatalf("ToUpperBytes(%q) = %q, input %q, want %q in place", s, got, upperB, wantUpper)
		}
		lowerS, upperS := []byte(s), []byte(s)
		if got := casestrings.UnsafeToLower(unsafeconv.UnsafeString(lowerS)); got != wantLower || string(lowerS) != wantLower {
			t.Fatalf("strings.UnsafeToLower(%q) = %q, backing %q, want %q", s, got, lowerS, wantLower)
		}
		if got := casestrings.UnsafeToUpper(unsafeconv.UnsafeString(upperS)); got != wantUpper || string(upperS) != wantUpper {
			t.Fatalf("strings.UnsafeToUpper(%q) = %q, backing %q, want %q", s, got, upperS, wantUpper)
		}
	}

	for c := range 256 {
		cs := string([]byte{byte(c)})
		check(cs)
		check("aB" + cs + "Cd")
		check(strings.Repeat("xY", 9) + cs + strings.Repeat("Zz", 9))
	}
	forEachString("aAzZ@[`{\xe9\xc9", 4, check)
	// Lengths around the word and 32-byte block sizes, the one letter to convert at every position.
	for n := range 100 {
		for pos := range n {
			plain := []byte(strings.Repeat("-4/", 40)[:n])
			up, low := bytes.Clone(plain), bytes.Clone(plain)
			up[pos], low[pos] = 'Q', 'q'
			check(string(up))
			check(string(low))
		}
		check(strings.Repeat("-4/", 40)[:n])
	}
}

// Test_Case_ASCIIOnly pins that case conversion stays in ASCII, where the stdlib maps every Unicode letter.
func Test_Case_ASCIIOnly(t *testing.T) {
	t.Parallel()
	for _, in := range []string{"\u00e9", "\u00c9", "\u0131", "\u212a", "\u00df", "\u01c5", "\u00e9 \u00c9 123 !", "\xff\xfe"} {
		if casestrings.ToLower(in) != in || casestrings.ToUpper(in) != in {
			t.Fatalf("case conversion changed the non-ASCII %q", in)
		}
	}
	if strings.ToUpper("\u00e9") == "\u00e9" {
		t.Fatal("the stdlib no longer upper-cases U+00E9: the documented difference is gone")
	}
	if got := casestrings.ToUpper("caf\u00e9 ok"); got != "CAF\u00e9 OK" {
		t.Fatalf("ToUpper left %q", got)
	}
}

// Test_CaseBytes_Aliasing pins that nothing to convert returns the argument itself, where bytes.ToLower copies.
func Test_CaseBytes_Aliasing(t *testing.T) {
	t.Parallel()
	for _, n := range []int{0, 1, 7, 8, 9, 31, 32, 33, 100} {
		plain := bytes.Repeat([]byte{'a'}, n)
		got := casebytes.ToLower(plain)
		if n > 0 && &got[0] != &plain[0] {
			t.Fatalf("ToLower(%d lower-case bytes) copied", n)
		}
		if std := bytes.ToLower(plain); n > 0 && &std[0] == &plain[0] {
			t.Fatalf("bytes.ToLower(%d lower-case bytes) no longer copies: the documented difference is gone", n)
		}

		upper := bytes.Repeat([]byte{'A'}, n)
		got = casebytes.ToUpper(upper)
		if n > 0 && &got[0] != &upper[0] {
			t.Fatalf("ToUpper(%d upper-case bytes) copied", n)
		}

		mixed := bytes.Repeat([]byte{'a'}, n)
		if n > 0 {
			mixed[n-1] = 'Q'
			if res := casebytes.ToLower(mixed); &res[0] == &mixed[0] || mixed[n-1] != 'Q' {
				t.Fatalf("ToLower(%d bytes ending in Q) did not return a fresh array", n)
			}
		}
	}
}

// Test_SplitTrimSeq_StdlibParity holds SplitTrimSeq to strings.SplitSeq with TrimSpace applied and empty elements dropped.
func Test_SplitTrimSeq_StdlibParity(t *testing.T) {
	t.Parallel()
	for _, sep := range []byte{',', ';', ' ', '\t', 'a', 0x00, 0xff} {
		alphabet := string([]byte{sep, ' ', '\t', '\n', 'x', 'a'})
		forEachString(alphabet, 6, func(s string) {
			var want []string
			for part := range strings.SplitSeq(s, string([]byte{sep})) {
				if part = strings.TrimSpace(part); part != "" {
					want = append(want, part)
				}
			}
			var got, gotBytes []string
			for elem := range SplitTrimSeq(s, sep) {
				got = append(got, elem)
			}
			for elem := range SplitTrimSeq([]byte(s), sep) {
				gotBytes = append(gotBytes, string(elem))
			}
			if strings.Join(got, "\x00") != strings.Join(want, "\x00") || len(got) != len(want) {
				t.Fatalf("SplitTrimSeq(%q, %#x) = %q, want %q", s, sep, got, want)
			}
			if strings.Join(gotBytes, "\x00") != strings.Join(want, "\x00") || len(gotBytes) != len(want) {
				t.Fatalf("SplitTrimSeq(bytes %q, %#x) = %q, want %q", s, sep, gotBytes, want)
			}
		})
	}
}

// Test_IndexControl_UnicodeParity holds the scans to unicode.IsControl on ASCII; the C1 range (U+0080-U+009F) is deliberately ignored.
func Test_IndexControl_UnicodeParity(t *testing.T) {
	t.Parallel()
	for b := range 256 {
		in := string([]byte{byte(b)})
		isControl := b < utf8.RuneSelf && unicode.IsControl(rune(b))
		want, wantExceptTab := -1, -1
		if isControl {
			want = 0
			if b != '\t' {
				wantExceptTab = 0
			}
		}
		if got := IndexControl(in); got != want {
			t.Fatalf("IndexControl(%#x) = %d, want %d", b, got, want)
		}
		if got := IndexControlExceptTab(in); got != wantExceptTab {
			t.Fatalf("IndexControlExceptTab(%#x) = %d, want %d", b, got, wantExceptTab)
		}

		for n := 0; n <= 40; n++ {
			s := strings.Repeat("a", n) + in
			if got := IndexControl(s); got != map[bool]int{true: n, false: -1}[isControl] {
				t.Fatalf("IndexControl(%d bytes then %#x) = %d", n, b, got)
			}
			if got := IndexControlExceptTab([]byte(s)); got != map[bool]int{true: n, false: -1}[isControl && b != '\t'] {
				t.Fatalf("IndexControlExceptTab(%d bytes then %#x) = %d", n, b, got)
			}
		}
	}

	for r := rune(0x80); r <= 0x9f; r++ {
		if !unicode.IsControl(r) {
			t.Fatalf("U+%04X is no longer a control character to the stdlib", r)
		}
		if got := IndexControl(string(r)); got != -1 {
			t.Fatalf("IndexControl(U+%04X) = %d, want it ignored", r, got)
		}
	}
}
