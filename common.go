// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/bits"
	"net"
	"os"
	"reflect"
	"runtime"
	"slices"

	"github.com/google/uuid"
)

const (
	defaultSecureTokenLength  = 32
	maxFastTokenEncodedLength = 43 // base64.RawURLEncoding.EncodedLen(32)
)

func readRandomOrPanic(dst []byte) {
	if _, err := rand.Read(dst); err != nil {
		// On supported Go versions (1.24+), crypto/rand.Read panics internally and
		// does not return errors. This check preserves explicit panic semantics if
		// the behavior changes or an alternate implementation is used in the future.
		// See: https://cs.opensource.google/go/go/+/refs/tags/go1.24.0:src/crypto/rand/rand.go
		panic(fmt.Errorf("utils: failed to read random bytes for token: %w", err))
	}
}

// Lookup table for ASCII whitespace characters (true = whitespace, false = not whitespace)
var whitespaceTable = [256]bool{
	'\t': true, // 9 - horizontal tab
	'\n': true, // 10 - line feed
	'\v': true, // 11 - vertical tab
	'\f': true, // 12 - form feed
	'\r': true, // 13 - carriage return
	' ':  true, // 32 - space
}

// UUIDv4 returns a Random (Version 4) UUID.
// The strength of the UUIDs is based on the strength of the crypto/rand package.
func UUIDv4() string {
	token, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("utils: failed to generate secure UUID: %w", err))
	}
	return token.String()
}

// UUID generates an universally unique identifier (UUID).
// This is an alias for UUIDv4 for backward compatibility.
func UUID() string {
	return UUIDv4()
}

// GenerateSecureToken generates a cryptographically secure random token encoded in base64.
// It uses crypto/rand for randomness and base64.RawURLEncoding for URL-safe output.
// If length is less than or equal to 0, it defaults to 32 bytes (256 bits of entropy).
// Panics if the random source fails.
func GenerateSecureToken(length int) string {
	if length <= 0 {
		length = defaultSecureTokenLength
	}

	if length == defaultSecureTokenLength {
		var randomBuf [defaultSecureTokenLength]byte
		src := randomBuf[:]
		readRandomOrPanic(src)

		var encoded [maxFastTokenEncodedLength]byte
		encodedLen := base64.RawURLEncoding.EncodedLen(length)
		base64.RawURLEncoding.Encode(encoded[:encodedLen], src)
		return string(encoded[:encodedLen])
	}

	bytes := make([]byte, length)
	readRandomOrPanic(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// SecureToken generates a secure token with 32 bytes of entropy.
// Panics if the random source fails. See GenerateSecureToken for details.
func SecureToken() string {
	return GenerateSecureToken(defaultSecureTokenLength)
}

// FunctionName returns function name
func FunctionName(fn any) string {
	if fn == nil {
		return ""
	}
	v := reflect.ValueOf(fn)
	if v.Kind() == reflect.Func {
		if v.IsNil() {
			return ""
		}
		pc := v.Pointer()
		f := runtime.FuncForPC(pc)
		if f == nil {
			return ""
		}
		return f.Name()
	}
	return v.Type().String()
}

// GetArgument check if key is in arguments
func GetArgument(arg string) bool {
	return slices.Contains(os.Args[1:], arg)
}

// IncrementIPRange Find available next IP address
func IncrementIPRange(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- { //nolint:modernize // slicesbackward: false positive
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ConvertToBytes returns integer size of bytes from human-readable string, ex. 42kb, 42M
// Decimal units are powers of 1000 (42k = 42000); binary units with an 'i' infix
// are powers of 1024 (42Ki = 43008). Note that ByteSize formats with powers of 1024,
// so use binary suffixes for a lossless round-trip.
// The number may have a '+', a point and an exponent (1.5e3K). The result is its
// exact product with the unit, rounded down and capped at math.MaxInt; digits past
// the 19 or 20 a uint64 holds are dropped, so a longer number can come out a byte
// short. Returns 0 if the string is unrecognized or negative.
func ConvertToBytes(humanReadableString string) int {
	strLen := len(humanReadableString)
	if strLen == 0 {
		return 0
	}

	// Fast path for plain byte values (e.g. "42", "42B", "42b").
	var sizeFast uint64
	maxInt := uint64(math.MaxInt)
	i := 0
	for ; i < strLen; i++ {
		c := humanReadableString[i]
		if c < '0' || c > '9' {
			break
		}
		d := uint64(c - '0')
		if sizeFast > maxInt/10 || (sizeFast == maxInt/10 && d > maxInt%10) {
			sizeFast = maxInt
		} else if sizeFast < maxInt {
			sizeFast = sizeFast*10 + d
		}
	}
	if i > 0 {
		if i == strLen {
			return int(sizeFast)
		}
		if i+1 == strLen {
			last := humanReadableString[i]
			if last == 'b' || last == 'B' {
				return int(sizeFast)
			}
		}
	}

	// Find the last digit position by scanning backwards
	// Also identify the unit prefix position in the same pass
	lastNumberPos := -1
	unitPrefixPos := 0
	for i := strLen - 1; i >= 0; i-- {
		c := humanReadableString[i]
		if c >= '0' && c <= '9' {
			lastNumberPos = i
			break
		}
		// Track the first letter position (unit prefix) from the end
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			unitPrefixPos = i
		}
	}

	// No digits found
	if lastNumberPos < 0 {
		return 0
	}

	numPart := humanReadableString[:lastNumberPos+1]

	mult := uint64(1)
	if unitPrefixPos > 0 {
		// An 'i' after the unit prefix selects binary units (KiB = 1024).
		binary := unitPrefixPos+1 < strLen &&
			(humanReadableString[unitPrefixPos+1] == 'i' || humanReadableString[unitPrefixPos+1] == 'I')
		var decMul, binMul uint64
		switch humanReadableString[unitPrefixPos] {
		case 'k', 'K':
			decMul, binMul = 1e3, 1<<10
		case 'm', 'M':
			decMul, binMul = 1e6, 1<<20
		case 'g', 'G':
			decMul, binMul = 1e9, 1<<30
		case 't', 'T':
			decMul, binMul = 1e12, 1<<40
		case 'p', 'P':
			decMul, binMul = 1e15, 1<<50
		default:
			decMul, binMul = 1, 1
		}
		mult = decMul
		if binary {
			mult = binMul
		}
	}

	// Scale in integers: a float product truncates 1.001K to 1000.
	return scaleDecimal(numPart, mult)
}

// maxExponent bounds the exponent scaleDecimal reads: a power of ten this large
// already saturates or vanishes, and the bound keeps the arithmetic from overflowing.
const maxExponent = 1 << 20

// scaleDecimal returns floor(s * mult) exactly, capped at math.MaxInt, or 0 if s
// is not a number or is negative. s is digits with an optional '+', point and
// exponent. The mantissa is the longest run of leading digits a uint64 holds; the
// digits after it are dropped (integer digits still count towards the magnitude),
// so such a number can come out a byte short.
func scaleDecimal(s string, mult uint64) int {
	var mantissa, kept uint64
	places, dropped, exp := -1, 0, 0 // places counts the digits after the point, -1 before it
scan:
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c-'0' <= 9:
			d := uint64(c - '0')
			if mantissa > maxUint64Cutoff || (mantissa == maxUint64Cutoff && d > maxUint64Cutlim) {
				if kept == 0 {
					// Nothing may follow a dropped digit, not even one that would fit: park
					// the mantissa and let a maximal one turn every later digit away.
					kept, mantissa = mantissa, math.MaxUint64
				}
				if places < 0 {
					dropped++
				}
			} else {
				mantissa = mantissa*10 + d
				if places >= 0 {
					places++
				}
			}
		case c == '.' && places < 0:
			places = 0
		case c == 'e' || c == 'E':
			var ok bool
			if exp, ok = parseExponent(s[i+1:]); !ok {
				return 0
			}
			break scan
		case c == '+' && i == 0:
		default:
			return 0
		}
	}
	if kept != 0 {
		mantissa = kept
	}
	if mantissa == 0 {
		return 0
	}

	// The value is mantissa * 10^shift * mult. Past 10^19 a nonzero mantissa
	// overflows a uint64, and a 128-bit division costs several times a 64-bit
	// one, so each case does only the work it needs.
	shift := exp + dropped - max(places, 0)
	if shift > 0 {
		if shift >= len(pow10) {
			return math.MaxInt
		}
		hi, lo := bits.Mul64(mantissa, pow10[shift])
		if hi != 0 {
			return math.MaxInt
		}
		mantissa = lo
	}
	hi, lo := bits.Mul64(mantissa, mult)
	if shift >= 0 {
		if hi != 0 || lo > uint64(math.MaxInt) {
			return math.MaxInt
		}
		return int(lo)
	}

	var quotient uint64
	if places := -shift; places < len(pow10) {
		divisor := pow10[places]
		switch {
		case hi == 0:
			quotient = lo / divisor
		case hi >= divisor: // the quotient needs more than 64 bits
			return math.MaxInt
		default:
			quotient, _ = bits.Div64(hi, lo, divisor)
		}
	} else if places < 2*len(pow10)-1 {
		// The product is below 2^114, so hi is below 10^19 and one division by
		// 10^19 leaves 51 bits to divide by the rest.
		quotient, _ = bits.Div64(hi, lo, pow10[len(pow10)-1])
		quotient /= pow10[places-(len(pow10)-1)]
	} // else the product is below 10^places
	if quotient > uint64(math.MaxInt) {
		return math.MaxInt
	}
	return int(quotient)
}

// parseExponent reads the part of a number after its 'e': an optional sign and
// at least one digit, saturating at maxExponent.
func parseExponent(t string) (int, bool) {
	neg := false
	if len(t) > 0 && (t[0] == '+' || t[0] == '-') {
		neg = t[0] == '-'
		t = t[1:]
	}
	if len(t) == 0 {
		return 0, false
	}
	exp := 0
	for i := range len(t) {
		c := t[i] - '0'
		if c > 9 {
			return 0, false
		}
		if exp < maxExponent {
			exp = exp*10 + int(c)
		}
	}
	if neg {
		exp = -exp
	}
	return exp, true
}
