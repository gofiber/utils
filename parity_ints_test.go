package utils

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

// The integer parsers promise strconv's answers for the base-10 grammar: the
// same value for everything strconv accepts, and the same class of error, by
// errors.Is against strconv.ErrSyntax and strconv.ErrRange, for everything it
// rejects. Two things differ on purpose and are pinned here too: the value
// returned with an error is 0, where strconv clamps it, and the Func of the
// *strconv.NumError names the parser that failed.
//
// strconv fails on the first problem it reads, so a number that is out of
// range before a stray byte is a range error, whatever follows. The narrow
// widths used to report a syntax error there.

// intParser puts one parser on a common footing: both instantiations, a
// signed or unsigned result widened to 64 bits, and the strconv call that
// answers for it.
type intParser struct {
	name   string
	bits   int
	signed bool
	str    func(string) (int64, uint64, error)
	bytes  func([]byte) (int64, uint64, error)
}

func signedParser[T Signed](name string, bits int, str func(string) (T, error), bytes func([]byte) (T, error)) intParser {
	return intParser{
		name: name, bits: bits, signed: true,
		str: func(s string) (int64, uint64, error) {
			v, err := str(s)
			return int64(v), 0, err
		},
		bytes: func(b []byte) (int64, uint64, error) {
			v, err := bytes(b)
			return int64(v), 0, err
		},
	}
}

func unsignedParser[T Unsigned](name string, bits int, str func(string) (T, error), bytes func([]byte) (T, error)) intParser {
	return intParser{
		name: name, bits: bits,
		str: func(s string) (int64, uint64, error) {
			v, err := str(s)
			return 0, uint64(v), err
		},
		bytes: func(b []byte) (int64, uint64, error) {
			v, err := bytes(b)
			return 0, uint64(v), err
		},
	}
}

var intParsers = []intParser{
	unsignedParser("ParseUint8", 8, ParseUint8[string], ParseUint8[[]byte]),
	unsignedParser("ParseUint16", 16, ParseUint16[string], ParseUint16[[]byte]),
	unsignedParser("ParseUint32", 32, ParseUint32[string], ParseUint32[[]byte]),
	unsignedParser("ParseUint", 64, ParseUint[string], ParseUint[[]byte]),
	unsignedParser("ParseNativeUint", strconv.IntSize, ParseNativeUint[string], ParseNativeUint[[]byte]),
	signedParser("ParseInt8", 8, ParseInt8[string], ParseInt8[[]byte]),
	signedParser("ParseInt16", 16, ParseInt16[string], ParseInt16[[]byte]),
	signedParser("ParseInt32", 32, ParseInt32[string], ParseInt32[[]byte]),
	signedParser("ParseInt", 64, ParseInt[string], ParseInt[[]byte]),
	signedParser("ParseNativeInt", strconv.IntSize, ParseNativeInt[string], ParseNativeInt[[]byte]),
}

// assertIntParity requires p to answer in as strconv does.
func assertIntParity(tb testing.TB, p intParser, in string) {
	tb.Helper()
	var wantI int64
	var wantU uint64
	var wantErr error
	if p.signed {
		wantI, wantErr = strconv.ParseInt(in, 10, p.bits)
	} else {
		wantU, wantErr = strconv.ParseUint(in, 10, p.bits)
	}
	var wantNum *strconv.NumError
	if wantErr != nil && !errors.As(wantErr, &wantNum) {
		tb.Fatalf("strconv returned %v, not a *strconv.NumError", wantErr)
	}

	gotI, gotU, err := p.str(in)
	gotBI, gotBU, errB := p.bytes([]byte(in))
	if gotI != gotBI || gotU != gotBU || (err == nil) != (errB == nil) {
		tb.Fatalf("%s(%q): string form (%d, %d, %v) differs from bytes form (%d, %d, %v)", p.name, in, gotI, gotU, err, gotBI, gotBU, errB)
	}

	if wantErr == nil {
		if err != nil || gotI != wantI || gotU != wantU {
			tb.Fatalf("%s(%q) = (%d, %d, %v), strconv gives (%d, %d, nil)", p.name, in, gotI, gotU, err, wantI, wantU)
		}
		return
	}
	var num *strconv.NumError
	switch {
	case !errors.As(err, &num):
		tb.Fatalf("%s(%q) error = %v, strconv gives %v", p.name, in, err, wantErr)
	case !errors.Is(err, wantNum.Err):
		tb.Fatalf("%s(%q) error = %v, strconv gives %v", p.name, in, err, wantErr)
	case num.Func != p.name || num.Num != in:
		tb.Fatalf("%s(%q) error = %#v, want Func %q and Num %q", p.name, in, num, p.name, in)
	case gotI != 0 || gotU != 0:
		tb.Fatalf("%s(%q) returned (%d, %d) alongside %v", p.name, in, gotI, gotU, err)
	}
}

func assertAllIntParity(tb testing.TB, in string) {
	tb.Helper()
	for _, p := range intParsers {
		assertIntParity(tb, p, in)
	}
}

func Test_ParseInts_StrconvParity_Exhaustive(t *testing.T) {
	t.Parallel()
	// Every string up to five bytes over digits that straddle the 8-bit and
	// 16-bit limits, both signs, a stray letter, and a space.
	forEachString("0159+-x ", 5, func(s string) {
		assertAllIntParity(t, s)
	})
}

func Test_ParseInts_StrconvParity_EveryByte(t *testing.T) {
	t.Parallel()
	// A single stray byte, of all 256, before, inside, and after digits that are
	// in range, at the limit, and out of range, so every byte class meets every
	// position of the word-at-a-time and the scalar paths.
	runs := []string{"7", "255", "256", "65536", "4294967296", "12345678", "123456789", "9223372036854775808", "18446744073709551616", "99999999999999999999999"}
	for b := range 256 {
		c := string([]byte{byte(b)})
		for _, run := range runs {
			for cut := 0; cut <= len(run); cut++ {
				assertAllIntParity(t, run[:cut]+c+run[cut:])
			}
			assertAllIntParity(t, c+run)
			assertAllIntParity(t, run+c)
			assertAllIntParity(t, "-"+run+c)
		}
	}
}

func Test_ParseInts_StrconvParity_Boundaries(t *testing.T) {
	t.Parallel()
	// For every width: the limits of the signed and unsigned types, and one
	// either side of them, as they are, with leading zeros, with either sign,
	// and with a trailing byte that is not a digit.
	var limits []*big.Int
	for _, bits := range []uint{8, 16, 32, 63, 64, 65, 128} {
		hi := new(big.Int).Lsh(big.NewInt(1), bits)
		half := new(big.Int).Rsh(hi, 1)
		for _, base := range []*big.Int{hi, half} {
			for _, delta := range []int64{-2, -1, 0, 1, 2} {
				limits = append(limits, new(big.Int).Add(base, big.NewInt(delta)))
			}
		}
	}
	limits = append(limits, big.NewInt(0), big.NewInt(1), big.NewInt(9), big.NewInt(10), big.NewInt(99), big.NewInt(100), big.NewInt(100000000))
	for _, limit := range limits {
		digits := limit.String()
		for _, sign := range []string{"", "+", "-"} {
			for _, zeros := range []int{0, 1, 7, 8, 9, 19, 20, 25} {
				for _, tail := range []string{"", "x", "0", "9", " ", "+", "-", "_", ".", "e1", "\x00", "\xff"} {
					assertAllIntParity(t, sign+strings.Repeat("0", zeros)+digits+tail)
				}
			}
		}
	}
	for _, in := range []string{"", "+", "-", "++1", "--1", "+-1", "-+1", " 1", "1 ", "0x10", "1_000", "\u0661\u0662\u0663", "1e3", "1.0", "\uff11\uff12\uff13"} {
		assertAllIntParity(t, in)
	}
}

func Test_ParseInts_StrconvParity_Random(t *testing.T) {
	t.Parallel()
	r := rand.New(rand.NewPCG(0x1a7, 0x9a45e)) //nolint:gosec // deterministic test data
	const strays = "x+-_. \x00\xff"
	for range 60000 {
		b := make([]byte, r.IntN(42))
		for i := range b {
			b[i] = byte('0' + r.IntN(10))
		}
		switch r.IntN(4) {
		case 0: // one stray byte somewhere
			if len(b) > 0 {
				b[r.IntN(len(b))] = strays[r.IntN(len(strays))]
			}
		case 1: // a sign in front, then maybe a stray byte
			if len(b) > 0 {
				b[0] = "+-"[r.IntN(2)]
			}
			if len(b) > 1 && r.IntN(2) == 0 {
				b[1+r.IntN(len(b)-1)] = strays[r.IntN(len(strays))]
			}
		}
		assertAllIntParity(t, string(b))
	}
}

func FuzzParseIntWidths(f *testing.F) {
	for _, seed := range []string{
		"256x", "0500x", "65536+", "4294967296 ", "9223372036854775808x", "18446744073709551616x",
		"-129x", "-128", "128", "255", "-256x", "+", "-", "", "00000000000000000000000256x",
		"12345678x", "1234567890123456789012x",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		assertAllIntParity(t, s)
	})
}

// Test_ParseInts_RangeBeforeSyntax names the cases that used to differ: the
// digits are over the unsigned limit of the width before the stray byte, so
// strconv has already failed with a range error when it reaches the byte.
func Test_ParseInts_RangeBeforeSyntax(t *testing.T) {
	t.Parallel()
	check := func(in string, bits int, signedOnly bool, want error) {
		t.Helper()
		for _, p := range intParsers {
			if p.bits != bits || (signedOnly && !p.signed) {
				continue
			}
			if _, _, err := p.str(in); !errors.Is(err, want) {
				t.Errorf("%s(%q) error = %v, want %v", p.name, in, err, want)
			}
		}
	}
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"256x", 8},
		{"0500x", 8},
		{"2560 ", 8},
		{"256+", 8},
		{"65536x", 16},
		{"00065536x", 16},
		{"4294967296x", 32},
		{"99999999999x", 32},
		{"00000000000000000000000000256x", 8},
		{"18446744073709551616x", 64},
	} {
		check(tc.in, tc.bits, false, strconv.ErrRange)
	}
	// A sign is a syntax error to the unsigned parsers, which stop there.
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"-256x", 8}, {"+256-", 8}, {"-4294967296x", 32}, {"+65536x", 16},
	} {
		check(tc.in, tc.bits, true, strconv.ErrRange)
		for _, p := range intParsers {
			if p.bits == tc.bits && !p.signed {
				if _, _, err := p.str(tc.in); !errors.Is(err, strconv.ErrSyntax) {
					t.Errorf("%s(%q) error = %v, want a syntax error", p.name, tc.in, err)
				}
			}
		}
	}
	// The same stray bytes are a syntax error as soon as the digits before them
	// fit the width's unsigned limit, which is what strconv compares with, even
	// where the signed type is already out of range.
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"255x", 8},
		{"128x", 8},
		{"65535x", 16},
		{"32768x", 16},
		{"4294967295x", 32},
		{"2147483648x", 32},
		{"9223372036854775808x", 64},
		{"18446744073709551615x", 64},
	} {
		check(tc.in, tc.bits, false, strconv.ErrSyntax)
	}
	for _, tc := range []struct {
		in   string
		bits int
	}{
		{"-129x", 8}, {"-32769x", 16}, {"-2147483649x", 32}, {"-9223372036854775809x", 64},
	} {
		check(tc.in, tc.bits, true, strconv.ErrSyntax)
	}
}

// Test_Format_StrconvParity_Exhaustive formats every value of the 8-bit and
// 16-bit types with every formatter wide enough to take it, and compares each
// with strconv: the narrow formatters have tables and shortcuts of their own.
func Test_Format_StrconvParity_Exhaustive(t *testing.T) {
	t.Parallel()
	appended := func(prefix string, got []byte, want string) bool {
		return string(got) == prefix+want
	}
	checkSigned := func(n int64, narrowest int) {
		t.Helper()
		want := strconv.FormatInt(n, 10)
		if narrowest <= 8 {
			if got := FormatInt8(int8(n)); got != want {
				t.Fatalf("FormatInt8(%d) = %q, want %q", n, got, want)
			}
		}
		if narrowest <= 16 {
			if got := FormatInt16(int16(n)); got != want {
				t.Fatalf("FormatInt16(%d) = %q, want %q", n, got, want)
			}
		}
		if got := FormatInt32(int32(n)); got != want {
			t.Fatalf("FormatInt32(%d) = %q, want %q", n, got, want)
		}
		if got := FormatInt(n); got != want {
			t.Fatalf("FormatInt(%d) = %q, want %q", n, got, want)
		}
		if !appended("", AppendInt(nil, n), want) || !appended("x", AppendInt([]byte("x"), n), want) {
			t.Fatalf("AppendInt(%d) differs from %q", n, want)
		}
		if got := ToString(n); got != want {
			t.Fatalf("ToString(%d) = %q, want %q", n, got, want)
		}
	}
	checkUnsigned := func(n uint64, narrowest int) {
		t.Helper()
		want := strconv.FormatUint(n, 10)
		if narrowest <= 8 {
			if got := FormatUint8(uint8(n)); got != want {
				t.Fatalf("FormatUint8(%d) = %q, want %q", n, got, want)
			}
		}
		if narrowest <= 16 {
			if got := FormatUint16(uint16(n)); got != want {
				t.Fatalf("FormatUint16(%d) = %q, want %q", n, got, want)
			}
		}
		if got := FormatUint32(uint32(n)); got != want {
			t.Fatalf("FormatUint32(%d) = %q, want %q", n, got, want)
		}
		if got := FormatUint(n); got != want {
			t.Fatalf("FormatUint(%d) = %q, want %q", n, got, want)
		}
		if !appended("", AppendUint(nil, n), want) || !appended("x", AppendUint([]byte("x"), n), want) {
			t.Fatalf("AppendUint(%d) differs from %q", n, want)
		}
		if got := ToString(n); got != want {
			t.Fatalf("ToString(%d) = %q, want %q", n, got, want)
		}
	}

	for n := math.MinInt16; n <= math.MaxInt16; n++ {
		narrowest := 16
		if n >= math.MinInt8 && n <= math.MaxInt8 {
			narrowest = 8
		}
		checkSigned(int64(n), narrowest)
	}
	for n := 0; n <= math.MaxUint16; n++ {
		narrowest := 16
		if n <= math.MaxUint8 {
			narrowest = 8
		}
		checkUnsigned(uint64(n), narrowest)
	}
}
