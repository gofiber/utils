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

// Calendar fields are converted with Hinnant's civil_from_days and
// days_from_civil (400-year eras of 146097 days, origin -0400-03-01 so every
// intermediate is non-negative); the tests sweep every day of 0..9999.
const (
	httpDateMinUnix        = -62167219200 // 0000-01-01T00:00:00Z
	httpDateMaxUnix        = 253402300799 // 9999-12-31T23:59:59Z
	secondsPerDay          = 86400
	daysPerEra             = 146097
	eraDaysToMarch         = daysPerEra - 60     // 0000-01-01 to -0400-03-01
	daysFromMarchEraToUnix = 719468 + daysPerEra // -0400-03-01 to 1970-01-01
	weekdayOfYearZero      = 6                   // 0000-01-01 was a Saturday
)

// civilFromDays converts days since 0000-01-01 (at most 3652424) into a
// proleptic Gregorian year, month (1..12), and day (1..31).
func civilFromDays(days uint64) (year, month, day uint64) { //nolint:nonamedreturns // the three calendar fields are only readable named
	z := days + eraDaysToMarch
	era := z / daysPerEra
	doe := z - era*daysPerEra                              // day of era, [0, 146096]
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365 // year of era, [0, 399]
	doy := doe - (365*yoe + yoe/4 - yoe/100)               // day of March-based year, [0, 365]
	mp := (5*doy + 2) / 153                                // March-based month, [0, 11]
	day = doy - (153*mp+2)/5 + 1
	month = mp + 3
	if mp >= 10 {
		month = mp - 9
	}
	year = yoe + era*400 - 400
	if month <= 2 {
		year++
	}
	return year, month, day
}

// unixFromCivil returns the Unix time of valid UTC calendar fields with the
// year in 0..9999, as parseRFC1123 guarantees.
func unixFromCivil(year, month, day, hour, minute, sec int) int64 {
	y := year + 400 // one era back keeps January/February of year 0 non-negative
	mp := month + 9
	if month > 2 {
		mp = month - 3
	} else {
		y--
	}
	era := y / 400
	yoe := y - era*400
	doy := (153*mp+2)/5 + day - 1
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	days := era*daysPerEra + doe - daysFromMarchEraToUnix
	return int64(days)*secondsPerDay + int64(hour*3600+minute*60+sec)
}

// putPair writes n < 100 as two zero-padded digits at b[i] and b[i+1].
func putPair(b []byte, i int, n uint64) {
	_ = b[i+1]
	b[i] = decimalPairs[2*n]
	b[i+1] = decimalPairs[2*n+1]
}

// AppendHTTPDate appends t in the RFC 9110 preferred HTTP date format
// ("Mon, 02 Jan 2006 15:04:05 GMT", net/http.TimeFormat) to dst and returns
// the extended slice. The output is byte-identical to
// t.UTC().AppendFormat(dst, http.TimeFormat) and always 29 bytes for the
// years 0..9999 that HTTP dates can represent; times outside that range
// delegate to time.AppendFormat.
func AppendHTTPDate(dst []byte, t time.Time) []byte {
	sec := t.Unix()
	if sec < httpDateMinUnix || sec > httpDateMaxUnix {
		// RFC 1123 assumes a four-digit year; keep stdlib behavior for the
		// unrepresentable rest instead of mis-padding it.
		return t.UTC().AppendFormat(dst, httpDateLayout)
	}
	u := uint64(sec - httpDateMinUnix) // seconds since 0000-01-01, so every split is unsigned
	days := u / secondsPerDay
	sod := u - days*secondsPerDay
	hour := sod / 3600
	sod -= hour * 3600
	minute := sod / 60
	second := sod - minute*60
	year, month, day := civilFromDays(days)

	var b [httpDateLen]byte
	copy(b[:], httpDateLayout)
	weekday := httpWeekdays[(days+weekdayOfYearZero)%7] // byte stores: copy would be a memmove call
	b[0], b[1], b[2] = weekday[0], weekday[1], weekday[2]
	putPair(b[:], 5, day)
	monthName := httpMonths[month-1]
	b[8], b[9], b[10] = monthName[0], monthName[1], monthName[2]
	century := year / 100
	putPair(b[:], 12, century)
	putPair(b[:], 14, year-century*100)
	putPair(b[:], 17, hour)
	putPair(b[:], 20, minute)
	putPair(b[:], 23, second)
	return append(dst, b[:]...)
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
	return time.Unix(unixFromCivil(year, month, day, hour, minute, sec), 0).UTC(), true
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
