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
	for j := len(ip) - 1; j >= 0; j-- {
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
// The number may have a '+', a point and an exponent (1.5e3K), with any number of
// digits. The result is the exact product rounded down, capped at math.MaxInt. It is
// 0 if there is no number or it is negative or malformed; an unknown unit counts as bytes.
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
	for j := strLen - 1; j >= 0; j-- {
		c := humanReadableString[j]
		if c >= '0' && c <= '9' {
			lastNumberPos = j
			break
		}
		// Track the first letter position (unit prefix) from the end
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			unitPrefixPos = j
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

	// A whole number is the digits read above, already capped.
	if lastNumberPos == i-1 {
		hi, lo := bits.Mul64(sizeFast, mult)
		if hi != 0 || lo > maxInt {
			return math.MaxInt
		}
		return int(lo)
	}

	// Integers, not floats: 1.001 * 1000 is 1000.9999999999999.
	return scaleDecimal(numPart, mult)
}

// scaleDecimal returns floor(s * mult) capped at math.MaxInt, or 0 if s is not a
// number or is negative. s is digits with an optional '+', point and exponent.
func scaleDecimal(s string, mult uint64) int {
	var mantissa uint64
	var exp int64
	end := len(s)            // where the digits end and the exponent starts
	places, dropped := -1, 0 // places counts the digits after the point, -1 before it
	droppedFrac := 0
	inexact := false // a nonzero digit did not fit the mantissa
scan:
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c-'0' <= 9:
			switch {
			case mantissa < maxUint64Cutoff:
				mantissa = mantissa*10 + uint64(c-'0')
				if places >= 0 {
					places++
				}
			case places < 0:
				dropped++ // an integer digit still scales the number
				inexact = inexact || c != '0'
			default:
				droppedFrac++
				inexact = inexact || c != '0'
			}
		case c == '.' && places < 0:
			places = 0
		case c == 'e' || c == 'E':
			var ok bool
			if exp, ok = parseExp(s, i+1); !ok {
				return 0
			}
			end = i
			break scan
		case c == '+' && i == 0:
		default:
			return 0
		}
	}
	if mantissa == 0 {
		return 0
	}

	// The mantissa holds a prefix of the digits, so the number lies in
	// [mantissa, mantissa+1) * 10^shift. If both ends give the same result, so does it.
	shift := exp + int64(dropped) - int64(max(places, 0))
	n := scaleMantissa(mantissa, shift, mult)
	if inexact && n != scaleMantissa(mantissa+1, shift, mult) {
		return scaleDigits(s, end, int64(droppedFrac)-shift+int64(dropped), mult)
	}
	return n
}

// scaleMantissa returns floor(mantissa * 10^shift * mult) capped at math.MaxInt.
// A 128-bit division costs several times a 64-bit one, so each case does only the
// work it needs.
func scaleMantissa(mantissa uint64, shift int64, mult uint64) int {
	if shift > 0 {
		if shift >= int64(len(pow10)) {
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
	if down := -shift; down < int64(len(pow10)) {
		divisor := pow10[down]
		switch {
		case hi == 0:
			quotient = lo / divisor
		case hi >= divisor: // the quotient needs more than 64 bits
			return math.MaxInt
		default:
			quotient, _ = bits.Div64(hi, lo, divisor)
		}
	} else if down < int64(2*len(pow10)-1) {
		// The product is below 2^114, so one division by 10^19 cannot overflow
		// and leaves 51 bits.
		quotient, _ = bits.Div64(hi, lo, pow10[len(pow10)-1])
		quotient /= pow10[down-int64(len(pow10)-1)]
	} // else the product is below 10^down
	if quotient > uint64(math.MaxInt) {
		return math.MaxInt
	}
	return int(quotient)
}

// scaleDigits multiplies the digits of s[:end] by mult, one at a time from the
// last, and keeps the digits of the product from the down-th on: the integer part
// of the number, whose point is down places from the end.
func scaleDigits(s string, end int, down int64, mult uint64) int {
	if down <= 0 {
		return math.MaxInt // the digits are a whole number of 10^19 or more
	}
	var result, carry uint64
	pos := int64(0) // the product digit made next
	for i := end - 1; i >= 0 || carry != 0; i-- {
		var d uint64
		if i >= 0 {
			c := s[i] - '0'
			if c > 9 { // the point or the sign
				continue
			}
			d = uint64(c)
		}
		t := d*mult + carry
		carry = t / 10
		if digit := t % 10; digit != 0 && pos >= down {
			if pos-down >= int64(len(pow10)-1) { // 10^19 is more than math.MaxInt
				return math.MaxInt
			}
			result += digit * pow10[pos-down] // below 2^64 even after the last sum
			if result > uint64(math.MaxInt) {
				return math.MaxInt
			}
		}
		pos++
	}
	return int(result)
}
