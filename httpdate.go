package utils

import (
	"time"
)

// httpDateLayout is the RFC 9110 preferred date format (net/http.TimeFormat),
// a fixed-width 29-byte layout that doubles as the formatting template: every
// separator byte is already in place, so AppendHTTPDate only overwrites the
// fields.
const (
	httpDateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"
	httpDateLen    = len(httpDateLayout)
)

var (
	httpWeekdays = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	httpMonths   = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	// httpDaysInMonth is indexed by time.Month; February is adjusted for
	// leap years in daysInHTTPMonth.
	httpDaysInMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// AppendHTTPDate appends t in the RFC 9110 preferred HTTP date format
// ("Mon, 02 Jan 2006 15:04:05 GMT", net/http.TimeFormat) to dst and returns
// the extended slice. The output is byte-identical to
// t.UTC().AppendFormat(dst, http.TimeFormat) and always 29 bytes for the
// years 0..9999 that HTTP dates can represent; times outside that range
// delegate to time.AppendFormat.
func AppendHTTPDate(dst []byte, t time.Time) []byte {
	t = t.UTC()
	year, month, day := t.Date()
	if year < 0 || year > 9999 {
		// RFC 1123 assumes a four-digit year; keep stdlib behavior for the
		// unrepresentable rest instead of mis-padding it.
		return t.AppendFormat(dst, httpDateLayout)
	}
	hour, minute, sec := t.Clock()
	off := len(dst)
	dst = append(dst, httpDateLayout...)
	b := dst[off : off+httpDateLen : off+httpDateLen]
	copy(b, httpWeekdays[t.Weekday()])
	b[5] = byte(day/10) + '0'
	b[6] = byte(day%10) + '0'
	copy(b[8:11], httpMonths[month-1])
	b[12] = byte(year/1000) + '0'
	b[13] = byte(year/100%10) + '0'
	b[14] = byte(year/10%10) + '0'
	b[15] = byte(year%10) + '0'
	b[17] = byte(hour/10) + '0'
	b[18] = byte(hour%10) + '0'
	b[20] = byte(minute/10) + '0'
	b[21] = byte(minute%10) + '0'
	b[23] = byte(sec/10) + '0'
	b[24] = byte(sec%10) + '0'
	return dst
}

// FormatHTTPDate returns t in the RFC 9110 preferred HTTP date format, equal
// to t.UTC().Format(http.TimeFormat) with a single allocation for the result.
func FormatHTTPDate(t time.Time) string {
	var buf [httpDateLen]byte
	return string(AppendHTTPDate(buf[:0], t))
}

// ParseHTTPDate parses an HTTP date the way net/http.ParseTime does:
// the RFC 9110 preferred format ("Mon, 02 Jan 2006 15:04:05 GMT") plus the
// obsolete RFC 850 and ANSI C asctime forms, with surrounding ASCII
// whitespace tolerated. Canonical preferred-format input takes a fast scalar
// path that never calls time.Parse; every other input — legacy formats,
// unusual casing, non-GMT zone names, padding — falls back to time.Parse
// with byte-for-byte stdlib semantics, including its errors. The returned
// time is in time.UTC on the fast path and whatever time.Parse yields on the
// fallback; the instants agree in both cases.
func ParseHTTPDate[S byteSeq](s S) (time.Time, error) {
	if t, ok := parseRFC1123(s); ok {
		return t, nil
	}
	return parseHTTPDateSlow(string(s))
}

// parseRFC1123 is the strict fast path for canonical
// "Mon, 02 Jan 2006 15:04:05 GMT" input: exact length, exact separators,
// exact-case names. It reports ok=false for anything else — including
// out-of-range fields — so the slow path can produce stdlib-identical
// results and errors for the long tail.
func parseRFC1123[S byteSeq](s S) (time.Time, bool) {
	if len(s) != httpDateLen {
		return time.Time{}, false
	}
	if s[3] != ',' || s[4] != ' ' || s[7] != ' ' || s[11] != ' ' || s[16] != ' ' ||
		s[19] != ':' || s[22] != ':' || s[25] != ' ' ||
		s[26] != 'G' || s[27] != 'M' || s[28] != 'T' {
		return time.Time{}, false
	}
	if !isHTTPWeekday(pack3(s[0], s[1], s[2])) {
		return time.Time{}, false
	}
	month := httpMonthNum(pack3(s[8], s[9], s[10]))
	if month == 0 {
		return time.Time{}, false
	}
	day, okDay := twoDigits(s[5], s[6])
	yh, okYh := twoDigits(s[12], s[13])
	yl, okYl := twoDigits(s[14], s[15])
	hour, okHour := twoDigits(s[17], s[18])
	minute, okMin := twoDigits(s[20], s[21])
	sec, okSec := twoDigits(s[23], s[24])
	if !okDay || !okYh || !okYl || !okHour || !okMin || !okSec {
		return time.Time{}, false
	}
	year := yh*100 + yl
	if hour > 23 || minute > 59 || sec > 59 ||
		day < 1 || day > daysInHTTPMonth(month, year) {
		// time.Parse rejects these; let it produce the exact error.
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), day, hour, minute, sec, 0, time.UTC), true
}

// parseHTTPDateSlow mirrors net/http.ParseTime: trim the whitespace
// textproto would, then try the three allowed layouts in order, returning
// the last error when none matches.
func parseHTTPDateSlow(s string) (time.Time, error) {
	start, end := 0, len(s)
	for start < end && isHTTPDateSpace(s[start]) {
		start++
	}
	for start < end && isHTTPDateSpace(s[end-1]) {
		end--
	}
	s = s[start:end]
	t, err := time.Parse(httpDateLayout, s)
	if err == nil {
		return t, nil
	}
	if t, err850 := time.Parse(time.RFC850, s); err850 == nil {
		return t, nil
	}
	t, errANSIC := time.Parse(time.ANSIC, s)
	if errANSIC == nil {
		return t, nil
	}
	return time.Time{}, errANSIC
}

// isHTTPDateSpace matches textproto.TrimString's notion of whitespace.
func isHTTPDateSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// pack3 packs three bytes for the name switches below.
func pack3(a, b, c byte) uint32 {
	return uint32(a)<<16 | uint32(b)<<8 | uint32(c)
}

// isHTTPWeekday reports whether w is one of the seven exact-case short
// weekday names. Like time.Parse, callers do not check the name against the
// date; any valid name is accepted.
func isHTTPWeekday(w uint32) bool {
	switch w {
	case 'S'<<16 | 'u'<<8 | 'n', 'M'<<16 | 'o'<<8 | 'n', 'T'<<16 | 'u'<<8 | 'e',
		'W'<<16 | 'e'<<8 | 'd', 'T'<<16 | 'h'<<8 | 'u', 'F'<<16 | 'r'<<8 | 'i',
		'S'<<16 | 'a'<<8 | 't':
		return true
	}
	return false
}

// httpMonthNum maps an exact-case short month name to 1..12, or 0 if w is
// not one.
func httpMonthNum(w uint32) int {
	switch w {
	case 'J'<<16 | 'a'<<8 | 'n':
		return 1
	case 'F'<<16 | 'e'<<8 | 'b':
		return 2
	case 'M'<<16 | 'a'<<8 | 'r':
		return 3
	case 'A'<<16 | 'p'<<8 | 'r':
		return 4
	case 'M'<<16 | 'a'<<8 | 'y':
		return 5
	case 'J'<<16 | 'u'<<8 | 'n':
		return 6
	case 'J'<<16 | 'u'<<8 | 'l':
		return 7
	case 'A'<<16 | 'u'<<8 | 'g':
		return 8
	case 'S'<<16 | 'e'<<8 | 'p':
		return 9
	case 'O'<<16 | 'c'<<8 | 't':
		return 10
	case 'N'<<16 | 'o'<<8 | 'v':
		return 11
	case 'D'<<16 | 'e'<<8 | 'c':
		return 12
	}
	return 0
}

// twoDigits parses two ASCII digits into an int, reporting whether both
// bytes were digits.
func twoDigits(a, b byte) (int, bool) {
	a -= '0'
	b -= '0'
	if a > 9 || b > 9 {
		return 0, false
	}
	return int(a)*10 + int(b), true
}

// daysInHTTPMonth returns the number of days in month for year, using the
// same leap-year rule as the time package.
func daysInHTTPMonth(month, year int) int {
	if month == 2 && year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 29
	}
	return httpDaysInMonth[month]
}
