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

// randRead is a package-level indirection for crypto/rand.Read so tests
// can override it to simulate failures.
var randRead = rand.Read

const (
	toLowerTable = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@abcdefghijklmnopqrstuvwxyz[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
	toUpperTable = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`ABCDEFGHIJKLMNOPQRSTUVWXYZ{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
)

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

// GenerateSecureToken generates a cryptographically secure random token encoded in base64.
// It uses crypto/rand for randomness and base64.RawURLEncoding for URL-safe output.
// If length is less than or equal to 0, it defaults to 32 bytes (256 bits of entropy).
// Panics if the random source fails.
func GenerateSecureToken(length int) string {
	if length <= 0 {
		length = 32
	}
	bytes := make([]byte, length)
	if _, err := randRead(bytes); err != nil {
		// On Go 1.24+, crypto/rand.Read panics internally and never returns an error.
		// On Go 1.23 and earlier, we panic for the same reasons: RNG failures indicate
		// a broken system state (uninitialized entropy pool, misconfigured VM, etc.)
		// that is almost certainly permanent rather than transient.
		// See: https://cs.opensource.google/go/go/+/refs/tags/go1.24.0:src/crypto/rand/rand.go
		//      https://go.dev/issue/66821
		panic(fmt.Errorf("utils: failed to read random bytes for token: %w", err))
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// SecureToken generates a secure token with 32 bytes of entropy.
// Panics if the random source fails. See GenerateSecureToken for details.
func SecureToken() string {
	return GenerateSecureToken(32)
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
