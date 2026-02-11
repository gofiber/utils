// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLower converts ascii string to lower-case
func ToLower(b string) string {
	return mapASCIIString(b, toLowerTable)
}

// ToLowerMut converts an ASCII string to lower-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func ToLowerMut(b string) string {
	return mapASCIIStringMut(b, toLowerTable)
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	return mapASCIIString(b, toUpperTable)
}

// ToUpperMut converts an ASCII string to upper-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func ToUpperMut(b string) string {
	return mapASCIIStringMut(b, toUpperTable)
}

// AddTrailingSlashString appends a trailing '/' to s if it does not already end with one.
// If the input already ends with '/', the original string is returned.
// A new string is returned only when a '/' needs to be appended.
func AddTrailingSlashString(s string) string {
	n := len(s)
	if n == 0 {
		return "/"
	}
	if s[n-1] == '/' {
		return s
	}
	buf := make([]byte, n+1)
	copy(buf, s)
	buf[n] = '/'
	return UnsafeString(buf)
}
