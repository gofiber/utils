package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
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

func FuzzPrefixSuffixFold(f *testing.F) {
	f.Add("Bearer abc.def", "bearer ")
	f.Add("no\rcache", "no-cache")
	f.Add("MULTIPART/FORM-DATA; boundary=X", "multipart/form-data")
	f.Add("application/ATOM+XML", "atom+xml")
	f.Add("\xe9abcdefgh", "\xc9abcdefgh")
	f.Add("", "")
	f.Fuzz(func(t *testing.T, s, needle string) {
		if want := refHasPrefixFold(s, needle); HasPrefixFold(s, needle) != want {
			t.Fatalf("HasPrefixFold(%q, %q) != %v", s, needle, want)
		} else if HasPrefixFold([]byte(s), needle) != want {
			t.Fatalf("HasPrefixFold(bytes %q, %q) != %v", s, needle, want)
		}
		if want := refHasSuffixFold(s, needle); HasSuffixFold(s, needle) != want {
			t.Fatalf("HasSuffixFold(%q, %q) != %v", s, needle, want)
		} else if HasSuffixFold([]byte(s), needle) != want {
			t.Fatalf("HasSuffixFold(bytes %q, %q) != %v", s, needle, want)
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

func FuzzURLEscape(f *testing.F) {
	f.Add("")
	f.Add("a b&c=d/e?f")
	f.Add("caf%C3%A9+r%C3%A9sum%C3%A9")
	f.Add("trailing%")
	f.Add("trailing%2")
	f.Add("%zz%1g%%41")
	f.Add("100%+tax")
	f.Add("\x00\x1f\x7f\xff")
	f.Fuzz(func(t *testing.T, s string) {
		if want := url.QueryEscape(s); want != string(AppendQueryEscape(nil, s)) ||
			want != string(AppendQueryEscape(nil, []byte(s))) {
			t.Fatalf("AppendQueryEscape(%q) diverges from net/url (%q)", s, want)
		}
		if want := url.PathEscape(s); want != string(AppendPathEscape(nil, s)) ||
			want != string(AppendPathEscape(nil, []byte(s))) {
			t.Fatalf("AppendPathEscape(%q) diverges from net/url (%q)", s, want)
		}

		checkUnescape := func(name, s string, got []byte, err error, want string, wantErr error) {
			t.Helper()
			if wantErr != nil {
				if !errors.Is(err, wantErr) {
					t.Fatalf("%s(%q) err = %v, want %v", name, s, err, wantErr)
				}
				if len(got) != 0 {
					t.Fatalf("%s(%q) produced output alongside error", name, s)
				}
				return
			}
			if err != nil {
				t.Fatalf("%s(%q) err = %v, want nil", name, s, err)
			}
			if string(got) != want {
				t.Fatalf("%s(%q) = %q, want %q", name, s, got, want)
			}
		}

		wantQ, wantQErr := url.QueryUnescape(s)
		gotQ, errQ := AppendQueryUnescape(nil, s)
		checkUnescape("AppendQueryUnescape", s, gotQ, errQ, wantQ, wantQErr)
		gotQB, errQB := AppendQueryUnescape(nil, []byte(s))
		checkUnescape("AppendQueryUnescape(bytes)", s, gotQB, errQB, wantQ, wantQErr)

		wantP, wantPErr := url.PathUnescape(s)
		gotP, errP := AppendPathUnescape(nil, s)
		checkUnescape("AppendPathUnescape", s, gotP, errP, wantP, wantPErr)
	})
}

func FuzzHTTPDate(f *testing.F) {
	f.Add("Sun, 06 Nov 1994 08:49:37 GMT")
	f.Add("Sunday, 06-Nov-94 08:49:37 GMT")
	f.Add("Sun Nov  6 08:49:37 1994")
	f.Add("  Tue, 29 Feb 2000 12:00:00 GMT\r\n")
	f.Add("Tue, 29 Feb 1900 12:00:00 GMT")
	f.Add("sun, 06 nov 1994 08:49:37 GMT")
	f.Add("Sun, 06 Nov 1994 08:49:37 UTC")
	f.Add("Mon, 01 Jan 0001 00:00:00 GMT")
	f.Add("not a date")
	f.Fuzz(func(t *testing.T, s string) {
		want, wantErr := refParseHTTPDate(s)
		got, err := ParseHTTPDate(s)
		if (err == nil) != (wantErr == nil) {
			t.Fatalf("ParseHTTPDate(%q) err = %v, reference err = %v", s, err, wantErr)
		}
		if wantErr == nil {
			if !got.Equal(want) {
				t.Fatalf("ParseHTTPDate(%q) = %v, want %v", s, got, want)
			}
			if got.Location().String() != want.Location().String() {
				t.Fatalf("ParseHTTPDate(%q) location = %v, want %v", s, got.Location(), want.Location())
			}
			// Canonical output re-parses through the fast path to the same
			// instant, and formatting agrees with the stdlib layout.
			if want.Year() >= 0 && want.Year() <= 9999 {
				formatted := FormatHTTPDate(want)
				if formatted != want.UTC().Format(httpDateLayout) {
					t.Fatalf("FormatHTTPDate(%v) = %q diverges from stdlib", want, formatted)
				}
				round, roundErr := ParseHTTPDate(formatted)
				if roundErr != nil || !round.Equal(want) {
					t.Fatalf("round trip of %q failed: %v, %v", formatted, round, roundErr)
				}
			}
		}
		if gotBytes, errBytes := ParseHTTPDate([]byte(s)); (errBytes == nil) != (err == nil) || (err == nil && !gotBytes.Equal(got)) {
			t.Fatalf("ParseHTTPDate(bytes %q) diverges from string form", s)
		}
	})
}

func FuzzAppendJSONString(f *testing.F) {
	f.Add("plain")
	f.Add("with \"quotes\" and \\slashes")
	f.Add("html <&> bits")
	f.Add("controls \x00\x1f\x7f")
	f.Add("caf\xc3\xa9 \xe2\x80\xa8 \xe2\x80\xa9")
	f.Add("broken \xff\xc3( utf8 \xc3")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		want, err := json.Marshal(s)
		if err != nil {
			t.Skipf("json.Marshal(%q) failed: %v", s, err)
		}
		if got := string(AppendJSONString(nil, s)); got != string(want) {
			t.Fatalf("AppendJSONString(%q) = %q, want %q", s, got, want)
		}
		if got := string(AppendJSONString(nil, []byte(s))); got != string(want) {
			t.Fatalf("AppendJSONString(bytes %q) = %q, want %q", s, got, want)
		}
	})
}

func FuzzParseIP(f *testing.F) {
	f.Add("127.0.0.1")
	f.Add("01.2.3.4")
	f.Add("255.255.255.256")
	f.Add("::ffff:1.2.3.4")
	f.Add("1:2:3:4:5:6:7::")
	f.Add("1::2::3")
	f.Add("fe80::1%eth0")
	f.Add("12345::")
	f.Add(":::")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		want, wantErr := netip.ParseAddr(s)
		hasColon := strings.ContainsRune(s, ':')

		got4, ok4 := ParseIPv4(s)
		if want4 := wantErr == nil && !hasColon; ok4 != want4 {
			t.Fatalf("ParseIPv4(%q) ok = %v, want %v (netip err: %v)", s, ok4, want4, wantErr)
		} else if want4 && got4 != want {
			t.Fatalf("ParseIPv4(%q) = %v, want %v", s, got4, want)
		}
		if _, okBytes := ParseIPv4([]byte(s)); okBytes != ok4 {
			t.Fatalf("ParseIPv4(bytes %q) diverges from string form", s)
		}

		got6, ok6 := ParseIPv6(s)
		if want6 := wantErr == nil && hasColon; ok6 != want6 {
			t.Fatalf("ParseIPv6(%q) ok = %v, want %v (netip err: %v)", s, ok6, want6, wantErr)
		} else if want6 && got6 != want {
			t.Fatalf("ParseIPv6(%q) = %v, want %v", s, got6, want)
		}
		if _, okBytes := ParseIPv6([]byte(s)); okBytes != ok6 {
			t.Fatalf("ParseIPv6(bytes %q) diverges from string form", s)
		}
	})
}

func FuzzCanonicalHeaderKey(f *testing.F) {
	f.Add("content-type")
	f.Add("Content-Type")
	f.Add("x--double--dash")
	f.Add("-leading")
	f.Add("bad key")
	f.Add("non-ascii-\xc3\xa9")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		want := http.CanonicalHeaderKey(s)
		if got := CanonicalHeaderKey(s); got != want {
			t.Fatalf("CanonicalHeaderKey(%q) = %q, want %q", s, got, want)
		}
		if got := string(CanonicalHeaderKey([]byte(s))); got != want {
			t.Fatalf("CanonicalHeaderKey(bytes %q) = %q, want %q", s, got, want)
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
