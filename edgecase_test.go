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

	"github.com/stretchr/testify/require"
)

// Differential edge-case suites for the helpers that promise byte-for-byte
// stdlib parity. Rather than hand-picking expectations, these enumerate
// structured input spaces — positional mutations, word-boundary offsets,
// malformed UTF-8, calendar boundaries — and require agreement with the
// stdlib reference on every input.

func Test_ParseHTTPDate_AllMonthsAndWeekdays(t *testing.T) {
	t.Parallel()
	// Twelve months exercise every httpMonthNum case; seven consecutive
	// days exercise every weekday name.
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
	// Length mutations: truncations and extensions.
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
	// Clock-field boundaries.
	for _, clock := range []string{"00:00:00", "23:59:59", "24:00:00", "23:60:00", "23:59:60", "99:99:99"} {
		assertHTTPDateParity(t, "Sat, 15 Jul 2023 "+clock+" GMT")
	}
}

// assertHTTPDateParity requires ParseHTTPDate to agree with the
// net/http.ParseTime reference on outcome, instant, and location.
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
		long[:31] + "%" + long[:1], // '%' exactly at index 31
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

	// A late error must still return the original dst content untouched.
	got, err := AppendQueryUnescape([]byte("keep"), long+"%zz")
	require.Error(t, err)
	require.Equal(t, "keep", string(got))
	gotPath, err := AppendPathUnescape([]byte("keep"), long+"%")
	require.Error(t, err)
	require.Equal(t, "keep", string(gotPath))
}

func Test_AppendUnescape_EveryEscapeByte(t *testing.T) {
	t.Parallel()
	// Every byte in both positions of an escape pair, plus that byte alone
	// after '%', pinned to net/url's accept set and error values.
	for i := range 256 {
		c := string([]byte{byte(i)})
		assertUnescapeParity(t, "%"+c+"0ab")
		assertUnescapeParity(t, "%4"+c+"ab")
		assertUnescapeParity(t, "ab%"+c)
	}
}

// assertUnescapeParity requires the two unescape helpers to agree with
// net/url on output and exact error values, for both instantiations.
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
	// Each special byte or sequence at every offset around the SWAR word
	// width, so escapes land in every lane and on both sides of the
	// word-loop/tail split.
	specials := []string{
		"\"", "\\", "\n", "\r", "\t", "\b", "\f", "\x00", "\x1f", "\x7f",
		"<", ">", "&", "\xff", "\xc3\xa9", "\u2028", "\u2029", "\xf0\x9f\x98\x80", "\xc3",
	}
	for _, sp := range specials {
		for pad := range 18 {
			in := strings.Repeat("a", pad) + sp + strings.Repeat("b", 17-pad)
			assertJSONStringParity(t, in)
			// And with the special at the very end, where the tail paths run.
			assertJSONStringParity(t, strings.Repeat("a", pad)+sp)
		}
	}
}

func Test_AppendJSONString_UTF8EdgeCases(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"\xc0\x80",                             // overlong NUL
		"\xc1\xbf",                             // overlong
		"\xe0\x80\xaf",                         // overlong 3-byte
		"\xed\xa0\x80",                         // UTF-16 surrogate half
		"\xed\xbf\xbf",                         // UTF-16 surrogate half
		"\xf4\x8f\xbf\xbf",                     // U+10FFFF, highest valid rune
		"\xf4\x90\x80\x80",                     // just above U+10FFFF
		"\xf8\x88\x80\x80\x80",                 // 5-byte form
		"\x80",                                 // lone continuation
		"\xbf\xbf",                             // lone continuations
		"\xc3",                                 // truncated 2-byte
		"\xe2\x80",                             // truncated 3-byte (prefix of U+2028)
		"\xf0\x9f\x98",                         // truncated 4-byte
		"\xef\xbb\xbfBOM",                      // BOM prefix
		"ok\xc3\xa9\xc3",                       // valid then truncated
		"\xe2\x80\xa8\xe2\x80\xa9\xe2\x80\xaa", // U+2028, U+2029, then U+202A (unescaped)
		strings.Repeat("\xff", 32),
		strings.Repeat("\xc3\xa9", 16),
	}
	for _, in := range inputs {
		assertJSONStringParity(t, in)
	}
}

// assertJSONStringParity requires AppendJSONString to agree with
// encoding/json.Marshal byte for byte, for both instantiations.
func assertJSONStringParity(t *testing.T, in string) {
	t.Helper()
	want, err := json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, string(want), string(AppendJSONString(nil, in)), "input %q", in)
	require.Equal(t, string(want), string(AppendJSONString(nil, []byte(in))), "input %q", in)
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
			assertParseIPParity(t, base[:pos]) // every truncation
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
		"1;2",       // field followed by a non-separator byte
		"1:2;3",     // same, after a valid group
		"fe80::1%%", // zone containing '%'
		"1::2%",     // empty zone after valid address
		"::%25eth0",
		"1:2:3:4:5:6:7:8%eth0", // full address with zone
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

// assertParseIPParity requires ParseIPv4/ParseIPv6 to accept exactly the
// per-family netip.ParseAddr inputs and produce identical addresses, for
// both instantiations.
func assertParseIPParity(t *testing.T, in string) {
	t.Helper()
	want, wantErr := netip.ParseAddr(in)
	hasColon := strings.ContainsRune(in, ':')

	got4, ok4 := ParseIPv4(in)
	_, ok4B := ParseIPv4([]byte(in))
	require.Equal(t, wantErr == nil && !hasColon, ok4, "ParseIPv4 input %q (netip err: %v)", in, wantErr)
	require.Equal(t, ok4, ok4B, "ParseIPv4 bytes input %q", in)
	if ok4 {
		require.Equal(t, want, got4, "ParseIPv4 input %q", in)
	} else {
		require.Equal(t, netip.Addr{}, got4, "ParseIPv4 input %q", in)
	}

	got6, ok6 := ParseIPv6(in)
	_, ok6B := ParseIPv6([]byte(in))
	require.Equal(t, wantErr == nil && hasColon, ok6, "ParseIPv6 input %q (netip err: %v)", in, wantErr)
	require.Equal(t, ok6, ok6B, "ParseIPv6 bytes input %q", in)
	if ok6 {
		require.Equal(t, want, got6, "ParseIPv6 input %q", in)
	} else {
		require.Equal(t, netip.Addr{}, got6, "ParseIPv6 input %q", in)
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
	// The rewrite path must copy, never mutate the caller's bytes.
	in := []byte("content-TYPE")
	out := CanonicalHeaderKey(in)
	require.Equal(t, "content-TYPE", string(in))
	require.Equal(t, "Content-Type", string(out))
}
