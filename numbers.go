package utils

import (
	"strconv"
)

const (
	maxUint64      = ^uint64(0)
	maxUint64Div10 = maxUint64 / 10
	maxUint64Mod10 = maxUint64 % 10
)

// ParseUint parses a decimal ASCII representation of an unsigned integer.
// It returns an error if the input contains non-digit characters or
// the value overflows uint64.
func ParseUint[S ~string | ~[]byte](s S) (uint64, error) {
	if len(s) == 0 {
		return 0, strconv.ErrSyntax
	}
	var n uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, strconv.ErrSyntax
		}
		if n > maxUint64Div10 || (n == maxUint64Div10 && uint64(c-'0') > maxUint64Mod10) {
			return 0, strconv.ErrRange
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}
