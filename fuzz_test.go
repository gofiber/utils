package utils

import (
	"bytes"
	"strconv"
	"testing"
)

// Fuzz targets comparing the SWAR implementations against the scalar
// references from the unit tests. The seeds cover the known traps (HTAB,
// CR|0x20, overlap tails, 19-20 digit boundaries); go test runs the seeds
// as regular tests, go test -fuzz explores beyond them.

func FuzzIndexFold(f *testing.F) {
	f.Add("max-age=0,No-Cache,private", "no-cache")
	f.Add("no\rcache", "no-cache")
	f.Add("aaaaaaaaaaaaNO-CACHE", "no-cache")
	f.Add("aaaaaaagz", "gzip")
	f.Add("__Proxy-Authorization: basic", "proxy-authorization")
	f.Add("\xe9abc", "\xc9abc")
	f.Add("", "")
	f.Fuzz(func(t *testing.T, s, needle string) {
		want := refIndexFold(s, needle)
		if got := IndexFold(s, needle); got != want {
			t.Fatalf("IndexFold(%q, %q) = %d, want %d", s, needle, got, want)
		}
		if got := IndexFold([]byte(s), needle); got != want {
			t.Fatalf("IndexFold(bytes %q, %q) = %d, want %d", s, needle, got, want)
		}
		if got := ContainsFold(s, needle); got != (want >= 0) {
			t.Fatalf("ContainsFold(%q, %q) = %v, want %v", s, needle, got, want >= 0)
		}
	})
}

func FuzzEqualFold(f *testing.F) {
	f.Add("Content-Type", "content-type")
	f.Add("no\rcache", "no-cache")
	f.Add("\xc9abc", "\xe9abc")
	f.Add("aaaaaaaaaaaaX", "aaaaaaaaaaaax")
	f.Fuzz(func(t *testing.T, a, b string) {
		want := bytes.Equal(asciiFoldUpperB([]byte(a)), asciiFoldUpperB([]byte(b)))
		if got := EqualFold(a, b); got != want {
			t.Fatalf("EqualFold(%q, %q) = %v, want %v", a, b, got, want)
		}
		if got := EqualFold([]byte(a), []byte(b)); got != want {
			t.Fatalf("EqualFold(bytes %q, %q) = %v, want %v", a, b, got, want)
		}
	})
}

func FuzzParse(f *testing.F) {
	f.Add("9223372036854775807")
	f.Add("9223372036854775808")
	f.Add("18446744073709551616")
	f.Add("00000000000000000001")
	f.Add("-9223372036854775808")
	f.Add("12345678\n2345678")
	f.Add("+0")
	f.Fuzz(func(t *testing.T, s string) {
		gotU, errU := ParseUint(s)
		wantU, wantErrU := strconv.ParseUint(s, 10, 64)
		if (errU == nil) != (wantErrU == nil) {
			t.Fatalf("ParseUint(%q) err = %v, strconv err = %v", s, errU, wantErrU)
		}
		switch {
		case wantErrU == nil && gotU != wantU:
			t.Fatalf("ParseUint(%q) = %d, want %d", s, gotU, wantU)
		case wantErrU != nil && gotU != 0:
			// Unlike strconv, this library returns 0 (not a clamped value)
			// alongside every error.
			t.Fatalf("ParseUint(%q) returned %d with error", s, gotU)
		}

		gotI, errI := ParseInt(s)
		wantI, wantErrI := strconv.ParseInt(s, 10, 64)
		if (errI == nil) != (wantErrI == nil) {
			t.Fatalf("ParseInt(%q) err = %v, strconv err = %v", s, errI, wantErrI)
		}
		switch {
		case wantErrI == nil && gotI != wantI:
			t.Fatalf("ParseInt(%q) = %d, want %d", s, gotI, wantI)
		case wantErrI != nil && gotI != 0:
			t.Fatalf("ParseInt(%q) returned %d with error", s, gotI)
		}
	})
}

func FuzzIndexNonQuotable(f *testing.F) {
	f.Add("a\tb")
	f.Add("aaaaaaaa\t\x1f")
	f.Add("\xff\xff\xff\xff\xff\xff\xff\"")
	f.Add("caf\xc3\xa9 r\xc3\xa9sum\xc3\xa9.pdf")
	f.Add("name\\quoted")
	f.Fuzz(func(t *testing.T, s string) {
		want := refIndexNonQuotable(s)
		if got := IndexNonQuotable(s); got != want {
			t.Fatalf("IndexNonQuotable(%q) = %d, want %d", s, got, want)
		}
		if got, wantASCII := IsASCII(s), isASCIIRef(s); got != wantASCII {
			t.Fatalf("IsASCII(%q) = %v, want %v", s, got, wantASCII)
		}
	})
}
