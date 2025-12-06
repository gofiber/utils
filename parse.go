package utils

import (
	"math"
	"strconv"
)

const maxFracDigits = 16

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

// ParseInt parses a decimal ASCII string or byte slice into an int64.
// Returns the parsed value and nil on success, else 0 and an error.
func ParseInt[S byteSeq](s S) (int64, error) {
	return parseSigned[S, int64]("ParseInt", s, math.MinInt64, math.MaxInt64)
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
	return parseUnsigned[S, uint8]("ParseUint8", s, uint8(math.MaxUint8))
}

// parseDigits parses a sequence of digits and returns the uint64 value.
// It returns an error if any non-digit is encountered or overflow happens.
func parseDigits[S byteSeq](s S, i int) (uint64, error) {
	var n uint64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, strconv.ErrSyntax
		}
		nn := n*10 + uint64(c)
		if nn < n {
			return 0, strconv.ErrRange
		}
		n = nn
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

	// Parse digits
	n, err := parseDigits(s, i)
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

	// Parse digits directly from index 0
	n, err := parseDigits(s, 0)
	// Check for overflow
	if err != nil {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: err}
	}
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

	var intPart uint64
	for i < len(s) {
		c := s[i] - '0'
		if c > 9 {
			break
		}
		nn := intPart*10 + uint64(c)
		if nn < intPart {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
		}
		intPart = nn
		i++
	}

	var fracPart uint64
	var fracDiv uint64 = 1
	var fracDigits int
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				break
			}
			if fracDigits >= maxFracDigits {
				return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
			}
			fracPart = fracPart*10 + uint64(c)
			fracDiv *= 10
			fracDigits++
			i++
		}
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
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
			}
			exp = exp*10 + int64(c)
			if !expSign && exp > 308 {
				exp = 309
			}
			if expSign && exp > 324 {
				exp = 325
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

	f := float64(intPart)
	if fracPart > 0 {
		f += float64(fracPart) / float64(fracDiv)
	}
	if exp != 0 {
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
