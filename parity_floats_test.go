package utils

import (
	"testing"
)

// The float parsers promise strconv.ParseFloat's value, bit for bit, for the
// decimal grammar and a syntax error for everything else, whatever strconv
// makes of it (inf, nan, hex floats, underscores). The random and generated
// inputs of parsefloat_test.go cover the large numbers; this walks every short
// string over the bytes the grammar is made of, where its edges are:
// a lone sign, a point with no digits on a side, an exponent with no digits or
// a sign of its own.
func Test_ParseFloat_ExhaustiveSmall(t *testing.T) {
	t.Parallel()
	forEachString("09.e+-", 5, func(s string) {
		requireFloatMatches(t, s)
	})
	forEachString("019.eE+-", 4, func(s string) {
		requireFloatMatches(t, s)
	})
	// The bytes strconv reads as something else, mixed into the grammar.
	forEachString("09.e+-_xpInfa", 3, func(s string) {
		requireFloatMatches(t, s)
	})
}
