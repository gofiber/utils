package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/utils/v2/swar"
)

const (
	maxUint64Cutoff = math.MaxUint64 / 10
	maxUint64Cutlim = math.MaxUint64 % 10
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// ParseUint parses a decimal ASCII string or byte slice into a uint64.
// It returns the parsed value and nil on success.
// If the input contains non-digit characters, it returns 0 and an error.
func ParseUint[S byteSeq](s S) (uint64, error) {
	return parseUnsigned[S, uint64]("ParseUint", s, uint64(math.MaxUint64))
}

// ParseNativeUint parses a decimal ASCII string or byte slice into a uint.
func ParseNativeUint[S byteSeq](s S) (uint, error) {
	return parseUnsigned[S, uint]("ParseNativeUint", s, math.MaxUint)
}

// ParseInt parses a decimal ASCII string or byte slice into an int64.
// Returns the parsed value and nil on success, else 0 and an error.
func ParseInt[S byteSeq](s S) (int64, error) {
	if len(s) > 0 && s[0] != '-' && s[0] != '+' && len(s) <= 19 {
		// At most 19 digits fit here, so two 8-digit SWAR steps plus a
		// scalar remainder can never overflow before the final range check.
		// A word with a non-digit lane falls through to the scalar loop,
		// which reports the syntax error. The body deliberately duplicates
		// the parseDigitsBig structure: routing through parseDigitsBig
		// instead costs +22% on the 9-digit shape (benchstat -count=10,
		// Go 1.25, Apple M2 Pro), mostly digit bookkeeping and the extra
		// call frame.
		var n uint64
		i := 0
		for ; len(s)-i >= 8; i += 8 {
			w := swar.Load8(s, i)
			if !isEightDigits(w) {
				break
			}
			n = n*100000000 + parse8Digits(w)
		}
		for ; i < len(s); i++ {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseInt", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint64(c)
		}
		if n > uint64(math.MaxInt64) {
			return 0, &strconv.NumError{Func: "ParseInt", Num: string(s), Err: strconv.ErrRange}
		}
		return int64(n), nil
	}

	return parseSigned[S, int64]("ParseInt", s, math.MinInt64, math.MaxInt64)
}

// ParseNativeInt parses a decimal ASCII string or byte slice into an int.
func ParseNativeInt[S byteSeq](s S) (int, error) {
	return parseSigned[S, int]("ParseNativeInt", s, math.MinInt, math.MaxInt)
}

// ParseInt32 parses a decimal ASCII string or byte slice into an int32.
func ParseInt32[S byteSeq](s S) (int32, error) {
	return parseSigned[S, int32]("ParseInt32", s, math.MinInt32, math.MaxInt32)
}

// ParseInt16 parses a decimal ASCII string or byte slice into an int16.
func ParseInt16[S byteSeq](s S) (int16, error) {
	return parseSigned[S, int16]("ParseInt16", s, math.MinInt16, math.MaxInt16)
}

// ParseInt8 parses a decimal ASCII string or byte slice into an int8.
func ParseInt8[S byteSeq](s S) (int8, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: "ParseInt8", Num: "", Err: strconv.ErrSyntax}
	}

	neg := false
	i := 0
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrSyntax}
	}

	if len(s)-i <= 3 {
		var n uint16
		for ; i < len(s); i++ {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint16(c)
		}
		if neg {
			if n > 128 {
				return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrRange}
			}
			if n == 128 {
				return math.MinInt8, nil
			}
			return -int8(n), nil
		}
		if n > math.MaxInt8 {
			return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrRange}
		}
		return int8(n), nil
	}

	return parseSigned[S, int8]("ParseInt8", s, math.MinInt8, math.MaxInt8)
}

// ParseUint32 parses a decimal ASCII string or byte slice into a uint32.
func ParseUint32[S byteSeq](s S) (uint32, error) {
	return parseUnsigned[S, uint32]("ParseUint32", s, uint32(math.MaxUint32))
}

// ParseUint16 parses a decimal ASCII string or byte slice into a uint16.
func ParseUint16[S byteSeq](s S) (uint16, error) {
	return parseUnsigned[S, uint16]("ParseUint16", s, uint16(math.MaxUint16))
}

// ParseUint8 parses a decimal ASCII string or byte slice into a uint8.
func ParseUint8[S byteSeq](s S) (uint8, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: "ParseUint8", Num: "", Err: strconv.ErrSyntax}
	}

	if len(s) <= 3 {
		var n uint16
		for i := range len(s) {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseUint8", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint16(c)
		}
		if n > math.MaxUint8 {
			return 0, &strconv.NumError{Func: "ParseUint8", Num: string(s), Err: strconv.ErrRange}
		}
		return uint8(n), nil
	}

	return parseUnsigned[S, uint8]("ParseUint8", s, uint8(math.MaxUint8))
}

// isEightDigits reports whether every lane of w is an ASCII digit ('0'..'9').
// High nibbles must all be 3 and low nibbles must not carry past 9 when 6 is
// added. A lane >= 0xFA can carry into its neighbor's sum, but such a lane
// already fails its own high-nibble test, so the answer stays exact.
func isEightDigits(w uint64) bool {
	const (
		highNibbles = 0xF0F0F0F0F0F0F0F0
		sixes       = 0x0606060606060606
		threes      = 0x3333333333333333
	)
	return (w&highNibbles)|((w+sixes)&highNibbles)>>4 == threes
}

// parse8Digits converts a word of 8 ASCII digit bytes (first digit in lane 0,
// most significant) to its numeric value using two pairwise multiply-combine
// steps. The caller must have validated that every lane is '0'..'9'.
func parse8Digits(w uint64) uint64 {
	const (
		digitZeros = 0x3030303030303030 // '0' in every lane
		pairMask   = 0x000000FF000000FF
		mul1       = 100 + (1000000 << 32)
		mul2       = 1 + (10000 << 32)
	)
	w -= digitZeros
	w = w*10 + w>>8 // adjacent digit pairs -> 2-digit values in even lanes
	return ((w&pairMask)*mul1 + (w>>16&pairMask)*mul2) >> 32
}

// parseDigits parses a run of fewer than 8 digits and returns the uint64
// value, or an error on the first non-digit. Callers must route runs of 8+
// bytes to parseDigitsBig — that is what makes overflow impossible here (at
// most 7 digits) and keeps this body under the inlining budget.
func parseDigits[S byteSeq](s S, i int) (uint64, error) {
	var n uint64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, strconv.ErrSyntax
		}
		n = n*10 + uint64(c)
	}
	return n, nil
}

// parseDigitsBig parses digit runs of 8+ bytes, consuming 8 digits per word
// while the running value is guaranteed to still fit: two full steps reach 16
// digits, and any value with <= 19 digits fits in uint64. A word with any
// non-digit lane falls through to the scalar loop, which reports syntax
// errors byte-precisely and applies the overflow checks.
func parseDigitsBig[S byteSeq](s S, i int) (uint64, error) {
	var n uint64
	digits := 0
	for len(s)-i >= 8 && digits+8 <= 19 {
		w := swar.Load8(s, i)
		if !isEightDigits(w) {
			break
		}
		n = n*100000000 + parse8Digits(w)
		digits += 8
		i += 8
	}
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, strconv.ErrSyntax
		}
		d := uint64(c)
		// Any value with <= 19 digits is guaranteed to fit in uint64.
		if digits >= 19 && (n > maxUint64Cutoff || (n == maxUint64Cutoff && d > maxUint64Cutlim)) {
			return 0, strconv.ErrRange
		}
		n = n*10 + d
		digits++
	}
	return n, nil
}

// parseSigned parses a decimal ASCII string or byte slice into a signed integer type T.
// It supports optional '+' or '-' prefix, checks for overflow and underflow, and returns (0, error) on error.
func parseSigned[S byteSeq, T Signed](fn string, s S, minRange, maxRange T) (T, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}

	neg := false
	i := 0
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}

	// Parse digits, taking the 8-digits-per-word path for long runs.
	var n uint64
	var err error
	if len(s)-i >= 8 {
		n, err = parseDigitsBig(s, i)
	} else {
		n, err = parseDigits(s, i)
	}
	if err != nil {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: err}
	}

	if !neg {
		// Check for overflow
		if n > uint64(int64(maxRange)) {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
		}
		return T(n), nil
	}

	// Check for underflow
	minAbs := uint64(-int64(minRange))
	if n > minAbs {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}

	return T(-int64(n)), nil
}

// parseUnsigned parses a decimal ASCII string or byte slice into an unsigned integer type T.
// It does not support sign prefixes, checks for overflow, and returns (0, error) on error.
func parseUnsigned[S byteSeq, T Unsigned](fn string, s S, maxRange T) (T, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}

	// Parse digits directly from index 0, taking the 8-digits-per-word path
	// for long runs.
	var n uint64
	var err error
	if len(s) >= 8 {
		n, err = parseDigitsBig(s, 0)
	} else {
		n, err = parseDigits(s, 0)
	}
	if err != nil {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: err}
	}
	// Check for overflow
	if n > uint64(maxRange) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}
	return T(n), nil
}

// parseFloat parses a decimal ASCII string or byte slice into a float64.
// It supports optional sign, fractional part and exponent. It returns (0, error)
// on error or overflow.
func parseFloat[S byteSeq](fn string, s S) (float64, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}
	i := 0
	neg := false
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}

	// Collect integer and fractional digits into a single mantissa and track
	// the decimal exponent. Digits beyond uint64 precision are dropped: integer
	// digits shift the exponent up, fractional digits are simply ignored.
	var mantissa uint64
	var exp10 int
	digits := 0
	for i < len(s) {
		c := s[i] - '0'
		if c > 9 {
			break
		}
		if mantissa > maxUint64Cutoff || (mantissa == maxUint64Cutoff && uint64(c) > maxUint64Cutlim) {
			exp10++
		} else {
			mantissa = mantissa*10 + uint64(c)
		}
		digits++
		i++
	}

	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				break
			}
			if mantissa < maxUint64Cutoff || (mantissa == maxUint64Cutoff && uint64(c) <= maxUint64Cutlim) {
				mantissa = mantissa*10 + uint64(c)
				exp10--
			}
			digits++
			i++
		}
	}

	if digits == 0 {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}

	var expSign bool
	var exp int64
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i == len(s) {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
		}
		switch s[i] {
		case '-':
			expSign = true
			i++
		case '+':
			i++
		}
		if i == len(s) {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
		}
		// Saturate the parsed exponent far beyond anything exp10 (bounded by
		// the input length) could pull back into range. Clamping to the final
		// range must wait until both are combined below.
		const maxParsedExp = int64(1) << 50
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
			}
			if exp < maxParsedExp {
				exp = exp*10 + int64(c)
			}
			i++
		}
	}

	if i != len(s) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}
	if expSign {
		exp = -exp
	}
	// Clamp the combined exponent to math.Pow10's saturation range so the
	// int conversion below cannot wrap on 32-bit platforms for extremely
	// long digit strings; larger magnitudes overflow (Inf, caught below) or
	// underflow (Pow10 returns 0) anyway.
	exp += int64(exp10)
	if exp > 309 {
		exp = 309
	} else if exp < -325 {
		exp = -325
	}

	f := float64(mantissa)
	// Skip scaling for a zero mantissa: 0 * Pow10(309) would be 0 * Inf = NaN
	// and turn inputs like "0e400" into a spurious range error.
	if exp != 0 && f != 0 {
		f *= math.Pow10(int(exp))
	}
	if neg {
		f = -f
	}
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}
	return f, nil
}

// ParseFloat64 parses a decimal ASCII string or byte slice into a float64. It
// delegates the actual parsing to parseFloat.
func ParseFloat64[S byteSeq](s S) (float64, error) {
	return parseFloat[S]("ParseFloat64", s)
}

// ParseFloat32 parses a decimal ASCII string or byte slice into a float32. It
// returns (0, false) on error or if the parsed value overflows float32.
func ParseFloat32[S byteSeq](s S) (float32, error) {
	f, err := parseFloat[S]("ParseFloat32", s)
	if err != nil {
		return 0, err
	}
	if f > math.MaxFloat32 || f < -math.MaxFloat32 {
		return 0, &strconv.NumError{Func: "ParseFloat32", Num: string(s), Err: strconv.ErrRange}
	}
	return float32(f), nil
}
