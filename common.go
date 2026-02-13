// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"net"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"

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
// Returns 0 if string is unrecognized
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
	var size float64

	if strings.IndexByte(numPart, '.') >= 0 {
		var err error
		size, err = ParseFloat64(numPart)
		if err != nil {
			return 0
		}
	} else {
		i64, err := ParseUint(numPart)
		if err != nil {
			return 0
		}
		size = float64(i64)
	}

	if unitPrefixPos > 0 {
		switch humanReadableString[unitPrefixPos] {
		case 'k', 'K':
			size *= 1e3
		case 'm', 'M':
			size *= 1e6
		case 'g', 'G':
			size *= 1e9
		case 't', 'T':
			size *= 1e12
		case 'p', 'P':
			size *= 1e15
		}
	}

	if size > float64(math.MaxInt) {
		return math.MaxInt
	}

	return int(size)
}
