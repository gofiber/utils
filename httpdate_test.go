package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// refParseHTTPDate mirrors net/http.ParseTime: trim the whitespace textproto
// would, then try the three allowed layouts in order. It is the behavioral
// reference for ParseHTTPDate in the unit and fuzz tests.
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
		{}, // zero time: Mon, 01 Jan 0001
		time.Unix(0, 0),
		time.Unix(784111777, 0), // the RFC 9110 example date
		time.Date(2000, time.February, 29, 23, 59, 59, 0, time.UTC),                 // leap day
		time.Date(1900, time.February, 28, 0, 0, 0, 0, time.UTC),                    // century non-leap
		time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),                 // last 4-digit second
		time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),                         // year 0000
		time.Date(2026, time.July, 23, 12, 34, 56, 789, time.UTC),                   // sub-second truncated by format
		time.Date(2026, time.July, 23, 20, 0, 0, 0, time.FixedZone("PST", -8*3600)), // non-UTC input
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

	// Appending must preserve existing dst content.
	got := AppendHTTPDate([]byte("Date: "), time.Unix(0, 0))
	require.Equal(t, "Date: Thu, 01 Jan 1970 00:00:00 GMT", string(got))
}

func Test_ParseHTTPDate(t *testing.T) {
	t.Parallel()

	// Canonical fast-path inputs round-trip through the formatter.
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

	// Every input, valid or not, must agree with the net/http.ParseTime
	// reference on outcome, instant, and location.
	inputs := []string{
		"Sun, 06 Nov 1994 08:49:37 GMT",
		"Sunday, 06-Nov-94 08:49:37 GMT",      // RFC 850
		"Sun Nov  6 08:49:37 1994",            // ANSI C asctime
		"  Sun, 06 Nov 1994 08:49:37 GMT\r\n", // padded, per textproto trim
		"sun, 06 nov 1994 08:49:37 GMT",       // lowercase names: stdlib folds case
		"Sun, 06 Nov 1994 08:49:37 UTC",       // non-GMT zone literal
		"Mon, 06 Nov 1994 08:49:37 GMT",       // wrong weekday for the date: accepted
		"Tue, 29 Feb 2000 12:00:00 GMT",       // leap day
		"Tue, 29 Feb 1900 12:00:00 GMT",       // century non-leap: rejected
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

// Test_HTTPDate_DaySweep pins the calendar arithmetic to the time package
// over every day of the years 0..9999: the formatter must match
// AppendFormat byte for byte, and the parser must return the identical
// Time value (not merely the same instant) that time.Date builds, for a
// second of the day that varies with the day so hours, minutes, and
// seconds are exercised across their ranges as well.
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
		if parsed != tm.UTC() {
			t.Fatalf("day %d: parsed %q to %v (%#v), want %v (%#v)", day, got, parsed, parsed, tm.UTC(), tm.UTC())
		}
	}

	// The instants just outside the four-digit-year range take the
	// stdlib path and must still match it.
	for _, sec := range []int64{httpDateMinUnix - 1, httpDateMaxUnix + 1} {
		tm := time.Unix(sec, 0)
		require.Equal(t, tm.UTC().Format(httpDateLayout), FormatHTTPDate(tm))
	}
}
