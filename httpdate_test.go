package utils

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// refParseHTTPDate mirrors net/http.ParseTime; it is the oracle for the unit and fuzz tests.
func refParseHTTPDate(s string) (time.Time, error) {
	start, end := 0, len(s)
	for start < end && isHTTPDateSpace(s[start]) {
		start++
	}
	for start < end && isHTTPDateSpace(s[end-1]) {
		end--
	}
	s = s[start:end]
	var t time.Time
	var err error
	for _, layout := range []string{httpDateLayout, time.RFC850, time.ANSIC} {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return t, err
}

func httpDateSamples() []time.Time {
	return []time.Time{
		{},
		time.Unix(0, 0),
		time.Unix(784111777, 0),
		time.Date(2000, time.February, 29, 23, 59, 59, 0, time.UTC),
		time.Date(1900, time.February, 28, 0, 0, 0, 0, time.UTC),
		time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
		time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.July, 23, 12, 34, 56, 789, time.UTC),
		time.Date(2026, time.July, 23, 20, 0, 0, 0, time.FixedZone("PST", -8*3600)),
	}
}

func Test_AppendHTTPDate(t *testing.T) {
	t.Parallel()
	for _, tm := range httpDateSamples() {
		want := tm.UTC().Format(httpDateLayout)
		require.Equal(t, want, string(AppendHTTPDate(nil, tm)))
		require.Equal(t, want, FormatHTTPDate(tm))
		require.Len(t, want, httpDateLen)
	}

	// Out-of-range years delegate to time.AppendFormat.
	for _, tm := range []time.Time{
		time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(-42, time.January, 1, 0, 0, 0, 0, time.UTC),
	} {
		want := tm.UTC().Format(httpDateLayout)
		require.Equal(t, want, string(AppendHTTPDate(nil, tm)))
	}

	got := AppendHTTPDate([]byte("Date: "), time.Unix(0, 0))
	require.Equal(t, "Date: Thu, 01 Jan 1970 00:00:00 GMT", string(got))
}

// Test_FormatHTTPDate_Allocs holds FormatHTTPDate to one allocation for every time, years outside
// 0..9999 included.
func Test_FormatHTTPDate_Allocs(t *testing.T) {
	// Not parallel: AllocsPerRun needs nothing else allocating.
	var sink string
	for _, tm := range append(httpDateSamples(),
		time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(-1, time.December, 31, 23, 59, 59, 0, time.UTC),
		time.Unix(httpDateMaxUnix+1, 0),
	) {
		require.Equal(t, tm.UTC().Format(httpDateLayout), FormatHTTPDate(tm))
		require.InDelta(t, 1, testing.AllocsPerRun(100, func() { sink = FormatHTTPDate(tm) }), 0, "%v", tm)
	}
	_ = sink
}

func Test_ParseHTTPDate(t *testing.T) {
	t.Parallel()

	for _, tm := range httpDateSamples() {
		s := FormatHTTPDate(tm)
		got, err := ParseHTTPDate(s)
		require.NoError(t, err, "input %q", s)
		require.True(t, got.Equal(tm.Truncate(time.Second)), "input %q: got %v", s, got)
		require.Same(t, time.UTC, got.Location(), "input %q", s)

		gotBytes, err := ParseHTTPDate([]byte(s))
		require.NoError(t, err)
		require.True(t, gotBytes.Equal(got))
	}

	inputs := []string{
		"Sun, 06 Nov 1994 08:49:37 GMT",
		"Sunday, 06-Nov-94 08:49:37 GMT",
		"Sun Nov  6 08:49:37 1994",
		"  Sun, 06 Nov 1994 08:49:37 GMT\r\n",
		"sun, 06 nov 1994 08:49:37 GMT",
		"Sun, 06 Nov 1994 08:49:37 UTC",
		"Mon, 06 Nov 1994 08:49:37 GMT", // wrong weekday: accepted
		"Tue, 29 Feb 2000 12:00:00 GMT",
		"Tue, 29 Feb 1900 12:00:00 GMT",
		"Tue, 30 Feb 2020 12:00:00 GMT",
		"Tue, 00 Feb 2020 12:00:00 GMT",
		"Tue, 32 Jan 2020 12:00:00 GMT",
		"Tue, 15 Jan 2020 24:00:00 GMT",
		"Tue, 15 Jan 2020 12:60:00 GMT",
		"Tue, 15 Jan 2020 12:00:60 GMT",
		"Foo, 15 Jan 2020 12:00:00 GMT",
		"Tue, 15 Foo 2020 12:00:00 GMT",
		"Tue, 15 Jan 2020 12:00:00 GM",
		"Tue, 15 Jan 2020 12:00:00",
		"",
		"not a date",
		"Tue, 1x Jan 2020 12:00:00 GMT",
		"Tue, 15 Jan 2x20 12:00:00 GMT",
	}
	for _, in := range inputs {
		want, wantErr := refParseHTTPDate(in)
		got, err := ParseHTTPDate(in)
		if wantErr != nil {
			require.Error(t, err, "input %q", in)
			continue
		}
		require.NoError(t, err, "input %q", in)
		require.True(t, got.Equal(want), "input %q: got %v, want %v", in, got, want)
		require.Equal(t, want.Location().String(), got.Location().String(), "input %q", in)
	}
}

func Benchmark_AppendHTTPDate(b *testing.B) {
	tm := time.Unix(784111777, 0)
	dst := make([]byte, 0, httpDateLen)
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			AppendHTTPDate(dst, tm)
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			tm.UTC().AppendFormat(dst, httpDateLayout)
		}
	})
}

func Benchmark_ParseHTTPDate(b *testing.B) {
	input := "Sun, 06 Nov 1994 08:49:37 GMT"
	b.Run("fiber", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := ParseHTTPDate(input); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := time.Parse(httpDateLayout, input); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Test_HTTPDate_DaySweep(t *testing.T) {
	t.Parallel()
	require.Equal(t, int64(httpDateMinUnix), time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Unix())
	require.Equal(t, int64(httpDateMaxUnix), time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC).Unix())

	var buf [httpDateLen]byte
	var want [httpDateLen]byte
	for day := int64(0); day <= (httpDateMaxUnix-httpDateMinUnix)/secondsPerDay; day++ {
		sod := (day * 7919) % secondsPerDay
		tm := time.Unix(httpDateMinUnix+day*secondsPerDay+sod, 0)
		got := AppendHTTPDate(buf[:0], tm)
		exp := tm.UTC().AppendFormat(want[:0], httpDateLayout)
		if string(got) != string(exp) {
			t.Fatalf("day %d: got %q, want %q", day, got, exp)
		}
		parsed, err := ParseHTTPDate(got)
		if err != nil {
			t.Fatalf("day %d: parse %q: %v", day, got, err)
		}
		if !parsed.Equal(tm) || parsed.Location() != time.UTC {
			t.Fatalf("day %d: parsed %q to %v (%#v), want %v (%#v)", day, got, parsed, parsed, tm.UTC(), tm.UTC())
		}
	}

	// Just outside the four-digit years, the stdlib path is taken.
	for _, sec := range []int64{httpDateMinUnix - 1, httpDateMaxUnix + 1} {
		tm := time.Unix(sec, 0)
		require.Equal(t, tm.UTC().Format(httpDateLayout), FormatHTTPDate(tm))
	}
}

func Test_HTTPDateLayout_IsTimeFormat(t *testing.T) {
	t.Parallel()
	require.Equal(t, http.TimeFormat, httpDateLayout)
	require.Len(t, http.TimeFormat, httpDateLen)
}

// assertHTTPDateStdlib holds ParseHTTPDate to http.ParseTime (same instant and location, or same
// *time.ParseError), which gets in trimmed since only ParseHTTPDate tolerates padding.
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
		"Sun, 06 Nov 1994 08:49:37 GMT",
		"Sunday, 06-Nov-94 08:49:37 GMT",
		"Sun Nov  6 08:49:37 1994",
		"Sun Nov 16 08:49:37 1994",
	}
	replacements := []byte{' ', ':', ',', '-', '.', '0', '5', '9', 'A', 'M', 'a', 'z', 0x00, 0x7f, 0xff}
	for _, base := range bases {
		assertHTTPDateStdlib(t, base)
		for pos := range len(base) {
			for _, r := range replacements {
				assertHTTPDateStdlib(t, base[:pos]+string([]byte{r})+base[pos+1:])
			}
			assertHTTPDateStdlib(t, base[:pos])
			assertHTTPDateStdlib(t, base[:pos]+base[pos+1:])
		}
		assertHTTPDateStdlib(t, base+" ")
		assertHTTPDateStdlib(t, base+"x")
		assertHTTPDateStdlib(t, "\t "+base+"\r\n")
	}

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

	for _, zone := range []string{"GMT", "UTC", "gmt", "Z", "+0000", "-0500", "EST", "GMT+1", ""} {
		assertHTTPDateStdlib(t, "Sun, 06 Nov 1994 08:49:37 "+zone)
		assertHTTPDateStdlib(t, "Sunday, 06-Nov-94 08:49:37 "+zone)
	}
	for _, clock := range []string{"00:00:00", "23:59:59", "24:00:00", "23:60:00", "23:59:60", "99:99:99", "1:2:3", "01:02", "01:02:03:04", "-1:00:00"} {
		assertHTTPDateStdlib(t, "Sat, 15 Jul 2023 "+clock+" GMT")
	}

	for _, year := range []int{0, 1, 4, 99, 100, 400, 1600, 1900, 1969, 1970, 1996, 2000, 2023, 2100, 2400, 9999} {
		for mi, month := range months {
			for _, day := range []int{0, 1, 28, 29, 30, 31, 32, 99} {
				assertHTTPDateStdlib(t, fmt.Sprintf("%s, %02d %s %04d 12:00:00 GMT", days[(year+mi+day)%7], day, month, year))
			}
		}
	}
	// RFC 850 two-digit years map to 1969..2068.
	for yy := range 100 {
		assertHTTPDateStdlib(t, fmt.Sprintf("Monday, 02-Jan-%02d 15:04:05 GMT", yy))
	}

	// Malformed and boundary years: length, sign, non-ASCII digits.
	for _, year := range []string{"0000", "0001", "10000", "99999", "-001", "+001", "2x21", "\u0662\u0660\u0662\u0661", " 2021"} {
		assertHTTPDateStdlib(t, "Sun, 06 Nov "+year+" 08:49:37 GMT")
	}
	for _, in := range []string{"", " ", "GMT", "0", "not a date", "Sun, 06 Nov 1994 08:49:37 GMT, Sun, 06 Nov 1994 08:49:37 GMT", "Sun, 06 Nov 1994 08:49:37 GMT\x00", "\x00Sun, 06 Nov 1994 08:49:37 GMT"} {
		assertHTTPDateStdlib(t, in)
	}
}

// Test_ParseHTTPDate_WhitespaceLeniency pins the documented difference: ASCII padding is ignored, unlike http.ParseTime.
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
	// \v, \f and the rest are not padding.
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
		// Every magnitude of Unix second, out-of-range years included.
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
