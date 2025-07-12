package utils

import (
	"math"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// ParseUint parses a decimal ASCII string or byte slice into a uint64.
// It returns the parsed value and true on success.
// If the input contains non-digit characters, it returns 0 and false.
func ParseUint[S byteSeq](s S) (uint64, bool) {
	return parseUnsigned[S, uint64](s, uint64(math.MaxUint64))
}

// ParseInt parses a decimal ASCII string or byte slice into an int64.
// Returns the parsed value and true on success, else 0 and false.
func ParseInt[S byteSeq](s S) (int64, bool) {
	return parseSigned[S, int64](s, math.MinInt64, math.MaxInt64)
}

// ParseInt32 parses a decimal ASCII string or byte slice into an int32.
func ParseInt32[S byteSeq](s S) (int32, bool) {
	return parseSigned[S, int32](s, math.MinInt32, math.MaxInt32)
}

// ParseInt8 parses a decimal ASCII string or byte slice into an int8.
func ParseInt8[S byteSeq](s S) (int8, bool) {
	return parseSigned[S, int8](s, math.MinInt8, math.MaxInt8)
}

// ParseUint32 parses a decimal ASCII string or byte slice into a uint32.
func ParseUint32[S byteSeq](s S) (uint32, bool) {
	return parseUnsigned[S, uint32](s, uint32(math.MaxUint32))
}

// ParseUint8 parses a decimal ASCII string or byte slice into a uint8.
func ParseUint8[S byteSeq](s S) (uint8, bool) {
	return parseUnsigned[S, uint8](s, uint8(math.MaxUint8))
}

// parseDigits parses a sequence of digits and returns the uint64 value and success.
// Returns (0, false) if any non-digit is encountered or overflow happens.
func parseDigits[S byteSeq](s S, i int) (uint64, bool) {
	var n uint64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, false
		}
		nn := n*10 + uint64(c)
		if nn < n {
			return 0, false
		}
		n = nn
	}
	return n, true
}

// parseSigned parses a decimal ASCII string or byte slice into a signed integer type T.
// It supports optional '+' or '-' prefix, checks for overflow and underflow, and returns (0, false) on error.
func parseSigned[S byteSeq, T Signed](s S, minRange, maxRange T) (T, bool) {
	if len(s) == 0 {
		return 0, false
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
		return 0, false
	}

	// Parse digits
	n, ok := parseDigits(s, i)
	if !ok {
		return 0, false
	}

	if !neg {
		// Check for overflow
		if n > uint64(int64(maxRange)) {
			return 0, false
		}
		return T(n), true
	}

	// Check for underflow
	minAbs := uint64(-int64(minRange))
	if n > minAbs {
		return 0, false
	}

	return T(-int64(n)), true
}

// parseUnsigned parses a decimal ASCII string or byte slice into an unsigned integer type T.
// It does not support sign prefixes, checks for overflow, and returns (0, false) on error.
func parseUnsigned[S byteSeq, T Unsigned](s S, maxRange T) (T, bool) {
	if len(s) == 0 {
		return 0, false
	}

	i := 0

	// Parse digits
	n, ok := parseDigits(s, i)
	// Check for overflow
	if !ok || n > uint64(maxRange) {
		return 0, false
	}
	return T(n), true
}

// ParseFloat parses a decimal ASCII string or byte slice into a float64.
// Supports optional sign, fractional part and exponent. Returns (0, false) on error or overflow.
func ParseFloat[S byteSeq](s S) (float64, bool) {
	if len(s) == 0 {
		return 0, false
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
		return 0, false
	}

	var intPart uint64
	for i < len(s) {
		c := s[i] - '0'
		if c > 9 {
			break
		}
		nn := intPart*10 + uint64(c)
		if nn < intPart {
			return 0, false
		}
		intPart = nn
		i++
	}

	var fracPart uint64
	var fracDiv uint64 = 1
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				break
			}
			if fracDiv < 1e16 {
				fracPart = fracPart*10 + uint64(c)
				fracDiv *= 10
			}
			i++
		}
	}

	var expSign bool
	var exp int64
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i == len(s) {
			return 0, false
		}
		switch s[i] {
		case '-':
			expSign = true
			i++
		case '+':
			i++
		}
		if i == len(s) {
			return 0, false
		}
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				return 0, false
			}
			exp = exp*10 + int64(c)
			if exp > 308 {
				exp = 309
			}
			i++
		}
	}

	if i != len(s) {
		return 0, false
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
		return 0, false
	}
	return f, true
}
