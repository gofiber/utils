package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func Test_ParseHTTPDate_AllMonthsAndWeekdays(t *testing.T) {
	t.Parallel()
	for month := time.January; month <= time.December; month++ {
		tm := time.Date(2021, month, 15, 1, 2, 3, 0, time.UTC)
		got, err := ParseHTTPDate(FormatHTTPDate(tm))
		require.NoError(t, err, "month %v", month)
		require.True(t, got.Equal(tm), "month %v", month)
	}
	for day := 1; day <= 7; day++ {
		tm := time.Date(2021, time.August, day, 23, 59, 59, 0, time.UTC)
		got, err := ParseHTTPDate(FormatHTTPDate(tm))
		require.NoError(t, err, "day %d", day)
		require.True(t, got.Equal(tm), "day %d", day)
	}
}

func Test_ParseHTTPDate_PositionalMutations(t *testing.T) {
	t.Parallel()
	base := "Sun, 06 Nov 1994 08:49:37 GMT"
	replacements := []byte{' ', ':', ',', '-', '.', '0', '5', '9', 'A', 'M', 'a', 'z', 0x00, 0x7F, 0xFF}
	for pos := range len(base) {
		for _, r := range replacements {
			mutated := base[:pos] + string(r) + base[pos+1:]
			assertHTTPDateParity(t, mutated)
		}
	}
	for cut := range len(base) {
		assertHTTPDateParity(t, base[:cut])
	}
	assertHTTPDateParity(t, base+" ")
	assertHTTPDateParity(t, base+"X")
	assertHTTPDateParity(t, " "+base)
	assertHTTPDateParity(t, "\t"+base+"\r\n")
}

func Test_ParseHTTPDate_CalendarBoundaries(t *testing.T) {
	t.Parallel()
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	years := []int{0, 4, 100, 400, 1900, 1996, 2000, 2020, 2023, 9999}
	days := []int{0, 1, 28, 29, 30, 31, 32, 99}
	for _, year := range years {
		for _, month := range months {
			for _, day := range days {
				in := fmt.Sprintf("Fri, %02d %s %04d 12:00:00 GMT", day, month, year)
				assertHTTPDateParity(t, in)
			}
		}
	}
	for _, clock := range []string{"00:00:00", "23:59:59", "24:00:00", "23:60:00", "23:59:60", "99:99:99"} {
		assertHTTPDateParity(t, "Sat, 15 Jul 2023 "+clock+" GMT")
	}
}

func assertHTTPDateParity(t *testing.T, in string) {
	t.Helper()
	want, wantErr := refParseHTTPDate(in)
	got, err := ParseHTTPDate(in)
	if wantErr != nil {
		require.Error(t, err, "input %q", in)
		return
	}
	require.NoError(t, err, "input %q", in)
	require.True(t, got.Equal(want), "input %q: got %v, want %v", in, got, want)
	require.Equal(t, want.Location().String(), got.Location().String(), "input %q", in)
}

func Test_AppendUnescape_LongAndBoundaryInputs(t *testing.T) {
	t.Parallel()
	long := strings.Repeat("abcdefgh", 8) // 64B, engages the vectorized scans
	inputs := []string{
		long,
		long + "%",
		long + "%4",
		long + "%41",
		long + "%zz",
		"%41" + long,
		long[:31] + "%" + long[:1],
		strings.Repeat("+", 64),
		strings.Repeat("%25", 32),
		strings.Repeat("a+", 32),
		"%",
		"+",
		"%%%",
		"%+1",
		"+%41+",
	}
	for _, in := range inputs {
		assertUnescapeParity(t, in)
	}

	// A late error must leave dst's content untouched.
	got, err := AppendQueryUnescape([]byte("keep"), long+"%zz")
	require.Error(t, err)
	require.Equal(t, "keep", string(got))
	gotPath, err := AppendPathUnescape([]byte("keep"), long+"%")
	require.Error(t, err)
	require.Equal(t, "keep", string(gotPath))
}

func Test_AppendUnescape_EveryEscapeByte(t *testing.T) {
	t.Parallel()
	for i := range 256 {
		c := string([]byte{byte(i)})
		assertUnescapeParity(t, "%"+c+"0ab")
		assertUnescapeParity(t, "%4"+c+"ab")
		assertUnescapeParity(t, "ab%"+c)
	}
}

func assertUnescapeParity(t *testing.T, in string) {
	t.Helper()
	wantQ, wantQErr := url.QueryUnescape(in)
	gotQ, errQ := AppendQueryUnescape(nil, in)
	gotQB, errQB := AppendQueryUnescape([]byte(nil), []byte(in))
	wantP, wantPErr := url.PathUnescape(in)
	gotP, errP := AppendPathUnescape(nil, in)
	if wantQErr != nil {
		require.Equal(t, wantQErr, errQ, "query input %q", in)
		require.Equal(t, wantQErr, errQB, "query input %q", in)
	} else {
		require.NoError(t, errQ, "query input %q", in)
		require.Equal(t, wantQ, string(gotQ), "query input %q", in)
		require.Equal(t, wantQ, string(gotQB), "query input %q", in)
	}
	if wantPErr != nil {
		require.Equal(t, wantPErr, errP, "path input %q", in)
	} else {
		require.NoError(t, errP, "path input %q", in)
		require.Equal(t, wantP, string(gotP), "path input %q", in)
	}
}

func Test_AppendEscape_EdgeCases(t *testing.T) {
	t.Parallel()
	all := make([]byte, 256)
	for i := range all {
		all[i] = byte(i)
	}
	inputs := []string{
		string(all),
		strings.Repeat(" ", 64),
		strings.Repeat("%", 64),
		strings.Repeat("a ", 32),
		strings.Repeat("unreserved-only.~_", 16),
		"ends with unsafe/",
		"/starts with unsafe",
	}
	for _, in := range inputs {
		require.Equal(t, url.QueryEscape(in), string(AppendQueryEscape(nil, in)), "input %q", in)
		require.Equal(t, url.PathEscape(in), string(AppendPathEscape(nil, in)), "input %q", in)
	}
}

func Test_AppendJSONString_BoundaryOffsets(t *testing.T) {
	t.Parallel()
	// Each special at every offset around the SWAR word width: every lane, both sides of the tail split.
	specials := []string{
		"\"", "\\", "\n", "\r", "\t", "\b", "\f", "\x00", "\x1f", "\x7f",
		"<", ">", "&", "\xff", "\xc3\xa9", "\u2028", "\u2029", "\xf0\x9f\x98\x80", "\xc3",
	}
	for _, sp := range specials {
		for pad := range 18 {
			in := strings.Repeat("a", pad) + sp + strings.Repeat("b", 17-pad)
			assertJSONStringParity(t, in)
			assertJSONStringParity(t, strings.Repeat("a", pad)+sp)
		}
	}
}

func Test_AppendJSONString_UTF8EdgeCases(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"\xc0\x80",
		"\xc1\xbf",
		"\xe0\x80\xaf",
		"\xed\xa0\x80",
		"\xed\xbf\xbf",
		"\xf4\x8f\xbf\xbf", // U+10FFFF, highest valid rune
		"\xf4\x90\x80\x80",
		"\xf8\x88\x80\x80\x80",
		"\x80",
		"\xbf\xbf",
		"\xc3",
		"\xe2\x80", // truncated 3-byte (prefix of U+2028)
		"\xf0\x9f\x98",
		"\xef\xbb\xbfBOM",
		"ok\xc3\xa9\xc3",
		"\xe2\x80\xa8\xe2\x80\xa9\xe2\x80\xaa", // U+2028, U+2029, then U+202A (unescaped)
		strings.Repeat("\xff", 32),
		strings.Repeat("\xc3\xa9", 16),
	}
	for _, in := range inputs {
		assertJSONStringParity(t, in)
	}
}

// Test_AppendJSONString_InvalidUTF8Escapes pins one \ufffd escape per bad byte on every Go version
// (a valid U+FFFD stays raw), which assertJSONStringParity cannot see.
func Test_AppendJSONString_InvalidUTF8Escapes(t *testing.T) {
	t.Parallel()
	const esc = `\ufffd`
	cases := []struct{ in, want string }{
		{"a\xffb", `"a` + esc + `b"`},
		{"\xc3", `"` + esc + `"`},
		{"x\xe2\x80", `"x` + esc + esc + `"`},
		{"\xc0\x80", `"` + esc + esc + `"`},
		{"\xed\xa0\x80", `"` + esc + esc + esc + `"`},
		{"\xf4\x90\x80\x80", `"` + esc + esc + esc + esc + `"`},
		{strings.Repeat("\xff", 16), `"` + strings.Repeat(esc, 16) + `"`},
		{strings.Repeat("a", 20) + "\xff" + strings.Repeat("b", 20), `"` + strings.Repeat("a", 20) + esc + strings.Repeat("b", 20) + `"`},
		{"\xef\xbf\xbd", "\"\xef\xbf\xbd\""},
	}

	probe, err := json.Marshal("\xff")
	require.NoError(t, err)
	marshalEscapes := string(probe) == `"`+esc+`"`

	for _, c := range cases {
		for _, got := range [][]byte{AppendJSONString(nil, c.in), AppendJSONString(nil, []byte(c.in))} {
			require.Equal(t, c.want, string(got), "input %q", c.in)
		}
		if marshalEscapes || utf8.ValidString(c.in) {
			want, err := json.Marshal(c.in)
			require.NoError(t, err)
			require.Equal(t, c.want, string(want), "input %q", c.in)
		}
	}
}

// assertJSONStringParity holds AppendJSONString to json.Marshal; for invalid UTF-8 only decoded
// strings can agree, as Go 1.27's json/v2 writes U+FFFD raw.
func assertJSONStringParity(t *testing.T, in string) {
	t.Helper()
	want, err := json.Marshal(in)
	require.NoError(t, err)
	for _, got := range [][]byte{AppendJSONString(nil, in), AppendJSONString(nil, []byte(in))} {
		if utf8.ValidString(in) {
			require.Equal(t, string(want), string(got), "input %q", in)
			continue
		}
		var gotStr, wantStr string
		require.NoError(t, json.Unmarshal(got, &gotStr), "input %q", in)
		require.NoError(t, json.Unmarshal(want, &wantStr), "input %q", in)
		require.Equal(t, wantStr, gotStr, "input %q", in)
	}
}

func Test_ParseIP_PositionalMutations(t *testing.T) {
	t.Parallel()
	bases := []string{
		"203.0.113.195",
		"::1",
		"2001:db8::8a2e:370:7334",
		"1:2:3:4:5:6:7:8",
		"::ffff:1.2.3.4",
		"1:2:3:4:5:6:1.2.3.4",
		"fe80::1%eth0",
	}
	replacements := []byte{'.', ':', '%', '0', '9', 'a', 'f', 'g', 'G', '-', ' ', 0x00, 0xFF}
	for _, base := range bases {
		for pos := range len(base) {
			for _, r := range replacements {
				assertParseIPParity(t, base[:pos]+string(r)+base[pos+1:])
			}
			assertParseIPParity(t, base[:pos])
		}
		assertParseIPParity(t, base+":")
		assertParseIPParity(t, base+".")
		assertParseIPParity(t, base+"0")
		assertParseIPParity(t, base+"%")
	}
}

func Test_ParseIP_MoreEdgeCases(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"1;2",
		"1:2;3",
		"fe80::1%%",
		"1::2%",
		"::%25eth0",
		"1:2:3:4:5:6:7:8%eth0",
		"0.0.0.256",
		"0.0.0.00",
		"255.255.255.2555",
		"::0.0.0.0",
		"0::00:0:00:0.0.0.0",
		"1:2:3:4:5:1.2.3.4",     // v4 tail too early without ellipsis
		"1:2:3:4:5:6:7:1.2.3.4", // v4 tail too late
		"::ffff:016.1.2.3",
		"ABCD:EF01:2345:6789:abcd:ef01:2345:6789",
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"00000::",
		"::00000",
		"2:2:2:2:2:2:2:2:",
		"\x00::",
		"1.2.3.4%eth0", // zones are IPv6-only
	}
	for _, in := range inputs {
		assertParseIPParity(t, in)
	}
}

// assertParseIPParity holds ParseIPv4/ParseIPv6 to netip.ParseAddr per family. Plain conditions, not
// require: the exhaustive loops call it hundreds of thousands of times.
func assertParseIPParity(tb testing.TB, in string) {
	tb.Helper()
	want, wantErr := netip.ParseAddr(in)
	hasColon := strings.ContainsRune(in, ':')

	got4, ok4 := ParseIPv4(in)
	got6, ok6 := ParseIPv6(in)
	if ok4 != (wantErr == nil && !hasColon) || ok6 != (wantErr == nil && hasColon) {
		tb.Fatalf("ParseIPv4/ParseIPv6(%q) ok = %v/%v, netip.ParseAddr gives %v, %v", in, ok4, ok6, want, wantErr)
	}
	if ok4 && got4 != want || ok6 && got6 != want || !ok4 && got4 != (netip.Addr{}) || !ok6 && got6 != (netip.Addr{}) {
		tb.Fatalf("ParseIPv4/ParseIPv6(%q) = %v/%v, netip.ParseAddr gives %v", in, got4, got6, want)
	}
	if _, bok := ParseIPv4([]byte(in)); bok != ok4 {
		tb.Fatalf("ParseIPv4(bytes %q) diverges from the string form", in)
	}
	if _, bok := ParseIPv6([]byte(in)); bok != ok6 {
		tb.Fatalf("ParseIPv6(bytes %q) diverges from the string form", in)
	}
}

func Test_CanonicalHeaderKey_PositionalMutations(t *testing.T) {
	t.Parallel()
	bases := []string{"content-type", "X-Request-Id", "a", "-", "--", "a-", "-a", "A--b"}
	for _, base := range bases {
		for pos := range len(base) {
			for c := range 256 {
				in := base[:pos] + string(byte(c)) + base[pos+1:]
				want := http.CanonicalHeaderKey(in)
				require.Equal(t, want, CanonicalHeaderKey(in), "input %q", in)
				require.Equal(t, want, string(CanonicalHeaderKey([]byte(in))), "input %q", in)
			}
		}
	}
}

func Test_CanonicalHeaderKey_InputImmutability(t *testing.T) {
	t.Parallel()
	in := []byte("content-TYPE")
	out := CanonicalHeaderKey(in)
	require.Equal(t, "content-TYPE", string(in))
	require.Equal(t, "Content-Type", string(out))
}
