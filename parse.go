package utils

import "math"

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

func parseSigned[S byteSeq, T Signed](s S, min, max T) (T, bool) {
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

	if !neg {
		if n > uint64(int64(max)) {
			return 0, false
		}
		return T(n), true
	}

	minAbs := uint64(-int64(min))
	if n > minAbs {
		return 0, false
	}

	return T(-int64(n)), true
}

func parseUnsigned[S byteSeq, T Unsigned](s S, max T) (T, bool) {
	if len(s) == 0 {
		return 0, false
	}

	i := 0

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

	if n > uint64(max) {
		return 0, false
	}

	return T(n), true
}
