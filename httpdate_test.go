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
