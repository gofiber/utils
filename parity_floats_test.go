package utils

import (
	"testing"
)

// Every short string over the grammar's bytes, where its edges are: a lone sign,
// a point or exponent with no digits. The third alphabet adds what strconv reads
// as inf, nan, hex floats and underscores, which must be syntax errors.
func Test_ParseFloat_ExhaustiveSmall(t *testing.T) {
	t.Parallel()
	forEachString("09.e+-", 5, func(s string) {
		requireFloatMatches(t, s)
	})
	forEachString("019.eE+-", 4, func(s string) {
		requireFloatMatches(t, s)
	})
	forEachString("09.e+-_xpInfa", 3, func(s string) {
		requireFloatMatches(t, s)
	})
}
