package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"math/rand/v2"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ParseHTTPDate, FormatHTTPDate

func Test_HTTPDateLayout_IsTimeFormat(t *testing.T) {
	t.Parallel()
	require.Equal(t, http.TimeFormat, httpDateLayout)
	require.Len(t, http.TimeFormat, httpDateLen)
}

// assertHTTPDateStdlib requires ParseHTTPDate to answer as net/http.ParseTime
// does: the same instant in the same location, or the same *time.ParseError.
// The one thing it adds is leniency about ASCII whitespace around the date, so
// in is handed to the stdlib trimmed, and the stdlib is shown to refuse it
// padded.
func assertHTTPDateStdlib(tb testing.TB, in string) {
	tb.Helper()
	trimmed := strings.Trim(in, " \t\r\n")
	want, wantErr := http.ParseTime(trimmed)
	for _, form := range []string{"string", "bytes"} {
		var got time.Time
		var err error
		if form == "string" {
			got, err = ParseHTTPDate(in)
		} else {
			got, err = ParseHTTPDate([]byte(in))
		}
		if wantErr != nil {
			if err == nil {
				tb.Fatalf("ParseHTTPDate(%s %q) = %v, nil; http.ParseTime fails: %v", form, in, got, wantErr)
			}
			var wantParse, gotParse *time.ParseError
			if !errors.As(wantErr, &wantParse) || !errors.As(err, &gotParse) || *wantParse != *gotParse {
				tb.Fatalf("ParseHTTPDate(%s %q) error = %#v, http.ParseTime gives %#v", form, in, err, wantErr)
			}
			continue
		}
		if err != nil {
			tb.Fatalf("ParseHTTPDate(%s %q) error = %v; http.ParseTime gives %v", form, in, err, want)
		}
		if !got.Equal(want) || got.Location().String() != want.Location().String() {
			tb.Fatalf("ParseHTTPDate(%s %q) = %v (%v), http.ParseTime gives %v (%v)", form, in, got, got.Location(), want, want.Location())
		}
	}
}

func Test_ParseHTTPDate_StdlibParity(t *testing.T) {
	t.Parallel()
	bases := []string{
		"Sun, 06 Nov 1994 08:49:37 GMT",  // RFC 9110 preferred
		"Sunday, 06-Nov-94 08:49:37 GMT", // RFC 850
		"Sun Nov  6 08:49:37 1994",       // asctime
		"Sun Nov 16 08:49:37 1994",
	}
	replacements := []byte{' ', ':', ',', '-', '.', '0', '5', '9', 'A', 'M', 'a', 'z', 0x00, 0x7f, 0xff}
	for _, base := range bases {
		assertHTTPDateStdlib(t, base)
		for pos := range len(base) {
			for _, r := range replacements {
				assertHTTPDateStdlib(t, base[:pos]+string([]byte{r})+base[pos+1:])
			}
			assertHTTPDateStdlib(t, base[:pos]) // every truncation
			assertHTTPDateStdlib(t, base[:pos]+base[pos+1:])
		}
		assertHTTPDateStdlib(t, base+" ")
		assertHTTPDateStdlib(t, base+"x")
		assertHTTPDateStdlib(t, "\t "+base+"\r\n")
	}

	// Every month name and weekday name in each layout, and the case and
	// spelling variants time.Parse accepts or refuses.
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	longDays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, month := range months {
		for _, name := range []string{month, strings.ToUpper(month), strings.ToLower(month), month[:2] + "x", "Sept"} {
			assertHTTPDateStdlib(t, fmt.Sprintf("%s, 15 %s 2021 01:02:03 GMT", days[i%7], name))
			assertHTTPDateStdlib(t, fmt.Sprintf("%s, 15-%s-21 01:02:03 GMT", longDays[i%7], name))
			assertHTTPDateStdlib(t, fmt.Sprintf("%s %s 15 01:02:03 2021", days[i%7], name))
		}
	}
	for _, name := range append(append([]string{}, days...), "sun", "SUN", "Sunday", "Su", "") {
		assertHTTPDateStdlib(t, name+", 06 Nov 1994 08:49:37 GMT")
	}

	// Zone names, offsets and clock readings that are or are not accepted.
	for _, zone := range []string{"GMT", "UTC", "gmt", "Z", "+0000", "-0500", "EST", "GMT+1", ""} {
		assertHTTPDateStdlib(t, "Sun, 06 Nov 1994 08:49:37 "+zone)
		assertHTTPDateStdlib(t, "Sunday, 06-Nov-94 08:49:37 "+zone)
	}
	for _, clock := range []string{"00:00:00", "23:59:59", "24:00:00", "23:60:00", "23:59:60", "99:99:99", "1:2:3", "01:02", "01:02:03:04", "-1:00:00"} {
		assertHTTPDateStdlib(t, "Sat, 15 Jul 2023 "+clock+" GMT")
	}

	// The calendar: leap years and the days around the end of every month.
	for _, year := range []int{0, 1, 4, 99, 100, 400, 1600, 1900, 1969, 1970, 1996, 2000, 2023, 2100, 2400, 9999} {
		for mi, month := range months {
			for _, day := range []int{0, 1, 28, 29, 30, 31, 32, 99} {
				assertHTTPDateStdlib(t, fmt.Sprintf("%s, %02d %s %04d 12:00:00 GMT", days[(year+mi+day)%7], day, month, year))
			}
		}
	}
	// RFC 850 has two-digit years, which time.Parse maps to 1969..2068.
	for yy := range 100 {
		assertHTTPDateStdlib(t, fmt.Sprintf("Monday, 02-Jan-%02d 15:04:05 GMT", yy))
	}

	// Years that do not have four digits, signs, and other numbers that are not one.
	for _, year := range []string{"0000", "0001", "10000", "99999", "-001", "+001", "2x21", "\u0662\u0660\u0662\u0661", " 2021"} {
		assertHTTPDateStdlib(t, "Sun, 06 Nov "+year+" 08:49:37 GMT")
	}
	for _, in := range []string{"", " ", "GMT", "0", "not a date", "Sun, 06 Nov 1994 08:49:37 GMT, Sun, 06 Nov 1994 08:49:37 GMT", "Sun, 06 Nov 1994 08:49:37 GMT\x00", "\x00Sun, 06 Nov 1994 08:49:37 GMT"} {
		assertHTTPDateStdlib(t, in)
	}
}

// Test_ParseHTTPDate_WhitespaceLeniency pins the one difference from
// http.ParseTime, which the doc comment states: ASCII whitespace around the
// date is ignored, where the stdlib rejects it.
func Test_ParseHTTPDate_WhitespaceLeniency(t *testing.T) {
	t.Parallel()
	const date = "Sun, 06 Nov 1994 08:49:37 GMT"
	want, err := http.ParseTime(date)
	require.NoError(t, err)
	for _, pad := range []string{" ", "\t", "\r", "\n", " \t\r\n", "\r\n"} {
		for _, in := range []string{pad + date, date + pad, pad + date + pad} {
			_, stdErr := http.ParseTime(in)
			require.Error(t, stdErr, "http.ParseTime no longer rejects padding: the documented difference is gone")
			got, err := ParseHTTPDate(in)
			require.NoError(t, err, "%q", in)
			require.True(t, got.Equal(want), "%q", in)
		}
	}
	// Anything else around it is as much an error as it is to the stdlib.
	for _, pad := range []string{"\v", "\f", "\u00a0", "\x00", "x"} {
		_, err := ParseHTTPDate(pad + date)
		require.Error(t, err, "%q", pad)
		_, err = ParseHTTPDate(date + pad)
		require.Error(t, err, "%q", pad)
	}
}

func FuzzParseHTTPDateStdlib(f *testing.F) {
	for _, seed := range []string{
		"Sun, 06 Nov 1994 08:49:37 GMT", "Sunday, 06-Nov-94 08:49:37 GMT", "Sun Nov  6 08:49:37 1994",
		"  Sun, 06 Nov 1994 08:49:37 GMT\r\n", "Tue, 29 Feb 1900 12:00:00 GMT", "sun, 06 nov 1994 08:49:37 GMT",
		"Mon, 01 Jan 0001 00:00:00 GMT", "", "not a date",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		assertHTTPDateStdlib(t, s)
	})
}

func Test_FormatHTTPDate_StdlibParity(t *testing.T) {
	t.Parallel()
	zones := []*time.Location{
		time.UTC, time.FixedZone("PST", -8*3600), time.FixedZone("IST", 5*3600+1800),
		time.FixedZone("", 13*3600), time.FixedZone("far", -23*3600-59*60-59),
	}
	r := rand.New(rand.NewPCG(0x4d17e, 0x5eed)) //nolint:gosec // deterministic test data
	check := func(tm time.Time) {
		t.Helper()
		want := tm.UTC().Format(http.TimeFormat)
		if got := FormatHTTPDate(tm); got != want {
			t.Fatalf("FormatHTTPDate(%v) = %q, want %q", tm, got, want)
		}
		if got := string(AppendHTTPDate([]byte("Date: "), tm)); got != "Date: "+want {
			t.Fatalf("AppendHTTPDate(%v) = %q, want %q", tm, got, "Date: "+want)
		}
	}
	for range 20000 {
		// Every magnitude of Unix second, including the years HTTP cannot say.
		sec := r.Int64N(1<<40) - 1<<39
		if r.IntN(8) == 0 {
			sec = r.Int64N(1<<62) - 1<<61
		}
		tm := time.Unix(sec, r.Int64N(1e9)).In(zones[r.IntN(len(zones))])
		check(tm)
	}
	for _, tm := range []time.Time{
		{},
		time.Unix(0, 0), time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
		time.Date(9999, 12, 31, 23, 59, 59, 0, time.FixedZone("early", -3600)), // 10000-01-01 UTC
		time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(0, 1, 1, 0, 0, 0, 0, time.FixedZone("late", 3600)), // -0001 UTC
		time.Date(-1, 12, 31, 23, 59, 59, 0, time.UTC),
		time.Date(12345, 6, 7, 8, 9, 10, 0, time.UTC),
		time.Date(-12345, 6, 7, 8, 9, 10, 0, time.UTC),
	} {
		check(tm)
	}
}

// ---------------------------------------------------------------------------
// AppendQueryUnescape, AppendPathUnescape in place

// Test_AppendUnescape_InPlaceParity decodes with the destination being the
// input itself, which the doc comments allow because decoding never grows the
// input, and requires what net/url gives: the same bytes, or the same error.
// After an error the input is consumed, so only the error is compared then.
func Test_AppendUnescape_InPlaceParity(t *testing.T) {
	t.Parallel()
	forEachString("%+aF1\x00\xff", 6, func(s string) {
		for _, mode := range []struct {
			name     string
			unescape func(string) (string, error)
			decode   func(dst, src []byte) ([]byte, error)
		}{
			{"query", url.QueryUnescape, AppendQueryUnescape[[]byte]},
			{"path", url.PathUnescape, AppendPathUnescape[[]byte]},
		} {
			want, wantErr := mode.unescape(s)
			buf := []byte(s)
			got, err := mode.decode(buf[:0], buf)
			if wantErr != nil {
				if err == nil || err.Error() != wantErr.Error() {
					t.Fatalf("%s in-place unescape of %q: error = %v, net/url gives %v", mode.name, s, err, wantErr)
				}
				continue
			}
			if err != nil || string(got) != want {
				t.Fatalf("%s in-place unescape of %q = %q, %v; net/url gives %q", mode.name, s, got, err, want)
			}
			if len(got) > 0 && &got[0] != &buf[0] {
				t.Fatalf("%s in-place unescape of %q did not decode into the input", mode.name, s)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// StatusMessage

func Test_StatusMessage_StdlibParity(t *testing.T) {
	t.Parallel()
	// Codes net/http does not name but the status table does, because Fiber
	// and fasthttp spell them.
	extra := map[int]string{
		306: "Switch Proxy",
	}
	for code := -20; code <= 1100; code++ {
		got, want := StatusMessage(code), http.StatusText(code)
		if extraText, ok := extra[code]; ok {
			require.Empty(t, want, "net/http now names %d: drop it from the documented differences", code)
			require.Equal(t, extraText, got, "code %d", code)
			continue
		}
		require.Equal(t, want, got, "code %d", code)
	}
}

// ---------------------------------------------------------------------------
// SplitHostPort, IsIPv4, IsIPv6, ParseIPv4, ParseIPv6

// leanSplitHostPort is checkSplitHostPort for the loops that run it a million
// times, where the assertion library would dominate: the same comparisons, as
// plain conditions.
func leanSplitHostPort(tb testing.TB, in string) {
	tb.Helper()
	wantHost, wantPort, err := net.SplitHostPort(in)
	host, port, ok := SplitHostPort(in)
	if ok != (err == nil) || (ok && (host != wantHost || port != wantPort)) || (!ok && (host != "" || port != "")) {
		tb.Fatalf("SplitHostPort(%q) = (%q, %q, %v), net.SplitHostPort gives (%q, %q, %v)", in, host, port, ok, wantHost, wantPort, err)
	}
	bh, bp, bok := SplitHostPort([]byte(in))
	if bok != ok || string(bh) != host || string(bp) != port || (!bok && (bh != nil || bp != nil)) {
		tb.Fatalf("SplitHostPort(bytes %q) = (%q, %q, %v), string form gives (%q, %q, %v)", in, bh, bp, bok, host, port, ok)
	}
}

func Test_SplitHostPort_StdlibParity_Exhaustive(t *testing.T) {
	t.Parallel()
	// Every string of at most eight bytes over the bytes that make the grammar:
	// both brackets, the colon, the percent sign, and a letter.
	forEachString("[]:%a", 8, func(s string) {
		leanSplitHostPort(t, s)
	})
	// And of at most five over a wider set, with a space, the dot and a NUL.
	forEachString("[]:%a1. \x00", 5, func(s string) {
		leanSplitHostPort(t, s)
	})
}

// leanIPParity is assertIPGrammarParity and assertParseIPParity as plain
// conditions: IsIPv4/IsIPv6 against net.ParseIP, ParseIPv4/ParseIPv6 against
// netip.ParseAddr.
func leanIPParity(tb testing.TB, in string) {
	tb.Helper()
	std := net.ParseIP(in)
	stdIsV4 := std != nil && !strings.Contains(in, ":")
	stdIsV6 := std != nil && strings.Contains(in, ":")
	if IsIPv4(in) != stdIsV4 || IsIPv6(in) != stdIsV6 {
		tb.Fatalf("IsIPv4/IsIPv6(%q) = %v/%v, net.ParseIP gives %v", in, IsIPv4(in), IsIPv6(in), std)
	}

	want, wantErr := netip.ParseAddr(in)
	hasColon := strings.Contains(in, ":")
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

func Test_IP_StdlibParity_Exhaustive(t *testing.T) {
	t.Parallel()
	forEachString("0:.1f", 8, func(s string) {
		leanIPParity(t, s)
	})
	forEachString("0:.%1fg", 5, func(s string) {
		leanIPParity(t, s)
	})
}

// ---------------------------------------------------------------------------
// IncrementIPRange

func Test_IncrementIPRange_Parity(t *testing.T) {
	t.Parallel()
	// The oracle is arithmetic: an address is a big-endian number, the next one
	// is the number plus one, and the last wraps to zero.
	bump := func(ip net.IP) net.IP {
		n := new(big.Int).SetBytes(ip)
		n.Add(n, big.NewInt(1))
		n.Mod(n, new(big.Int).Lsh(big.NewInt(1), uint(8*len(ip))))
		return n.FillBytes(make([]byte, len(ip)))
	}
	check := func(ip net.IP) {
		t.Helper()
		want := bump(ip)
		got := append(net.IP(nil), ip...)
		IncrementIPRange(got)
		if !got.Equal(want) || len(got) != len(ip) || string(got) != string(want) {
			t.Fatalf("IncrementIPRange(%v) = %v, want %v", []byte(ip), []byte(got), []byte(want))
		}
		// netip says the same wherever the successor exists.
		if addr, ok := netip.AddrFromSlice(ip); ok {
			if next := addr.Next(); next.IsValid() && string(next.AsSlice()) != string(got) {
				t.Fatalf("IncrementIPRange(%v) = %v, netip.Addr.Next gives %v", ip, got, next)
			}
		}
	}

	r := rand.New(rand.NewPCG(0x1c4, 0x7fe)) //nolint:gosec // deterministic test data
	for _, size := range []int{4, 16} {
		// Every run of 0xff at the end, which is how far the carry travels.
		for run := 0; run <= size; run++ {
			ip := make(net.IP, size)
			for i := range ip {
				ip[i] = byte(r.IntN(255)) // never 0xff, so the run is exactly what is set below
			}
			for i := size - run; i < size; i++ {
				ip[i] = 0xff
			}
			check(ip)
			check(make(net.IP, size))
		}
		for range 2000 {
			ip := make(net.IP, size)
			for i := range ip {
				ip[i] = byte(r.IntN(256))
			}
			check(ip)
		}
	}
	// The v4-in-v6 spelling is a 16-byte number to this function.
	check(net.IPv4(192, 168, 1, 255))
	check(net.IPv4(255, 255, 255, 255))
	check(net.IPv6zero)
	check(net.IPv6loopback)

	// Other lengths are numbers too, and nothing is something to leave alone.
	for _, ip := range []net.IP{{0}, {0xff}, {1, 0xff}, {0xff, 0xff, 0xff}, {0, 0, 0, 0, 0}} {
		check(ip)
	}
	require.NotPanics(t, func() {
		IncrementIPRange(nil)
		IncrementIPRange(net.IP{})
	})
}

// ---------------------------------------------------------------------------
// ReadFile

func Test_ReadFile_StdlibParity(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	contents := map[string][]byte{
		"empty.txt":  {},
		"small.txt":  []byte("hello\n"),
		"binary.bin": {0, 1, 2, 0xff, 0xfe, '\n', '\r'},
		"large.bin":  make([]byte, 1<<20+7),
	}
	for i := range contents["large.bin"] {
		contents["large.bin"][i] = byte(i * 31)
	}
	for name, data := range contents {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), data, 0o600))
	}
	require.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o750))

	fsys := http.Dir(dir)
	for name := range contents {
		want, err := os.ReadFile(filepath.Join(dir, name)) // #nosec G304
		require.NoError(t, err)

		fromOS, err := ReadFile(filepath.Join(dir, name), nil)
		require.NoError(t, err)
		require.Equal(t, want, fromOS, name)

		fromFS, err := ReadFile("/"+name, fsys)
		require.NoError(t, err)
		require.Equal(t, want, fromFS, name)
	}

	// A missing file is fs.ErrNotExist from either source, as the stdlib has it.
	_, stdErr := os.ReadFile(filepath.Join(dir, "missing")) // #nosec G304
	require.ErrorIs(t, stdErr, fs.ErrNotExist)
	_, err := ReadFile(filepath.Join(dir, "missing"), nil)
	require.ErrorIs(t, err, fs.ErrNotExist)
	_, err = ReadFile("/missing", fsys)
	require.ErrorIs(t, err, fs.ErrNotExist)

	// A directory cannot be read as a file by any of the three.
	_, stdErr = os.ReadFile(filepath.Join(dir, "sub")) // #nosec G304
	require.Error(t, stdErr)
	_, err = ReadFile(filepath.Join(dir, "sub"), nil)
	require.Error(t, err)
	_, err = ReadFile("/sub", fsys)
	require.Error(t, err)

	// An empty file is an empty slice, not an error.
	empty, err := ReadFile("/empty.txt", fsys)
	require.NoError(t, err)
	require.Empty(t, empty)
}

// ---------------------------------------------------------------------------
// UnsafeString, UnsafeBytes, CopyString, CopyBytes

func Test_UnsafeConversions_Parity(t *testing.T) {
	t.Parallel()
	for n := 0; n <= 70; n++ {
		src := make([]byte, n)
		for i := range src {
			src[i] = byte(i*37 + 1)
		}

		view := UnsafeString(src)
		require.Equal(t, string(src), view, n)
		require.Len(t, view, n)
		if n > 0 {
			// A view, not a copy: it sees what the slice is changed to.
			src[0] ^= 0xff
			require.Equal(t, string(src), view, n)
			require.Same(t, &src[0], &UnsafeBytes(view)[0])
			src[0] ^= 0xff
		}

		str := string(src)
		bs := UnsafeBytes(str)
		require.Equal(t, str, string(bs)) // an empty string may give a nil slice, which is still empty
		require.Len(t, bs, n)
		require.Equal(t, n, cap(bs))

		copied := CopyString(str)
		require.Equal(t, str, copied)
		if n > 1 { // one-byte strings are shared by the runtime
			require.NotSame(t, &UnsafeBytes(str)[0], &UnsafeBytes(copied)[0], "CopyString must not share the original's bytes")
		}

		copiedBytes := CopyBytes(src)
		require.Equal(t, src, copiedBytes)
		require.Len(t, copiedBytes, n)
		if n > 0 {
			require.NotSame(t, &src[0], &copiedBytes[0])
			src[0] ^= 0xff
			require.NotEqual(t, src[0], copiedBytes[0])
		}
	}
}
