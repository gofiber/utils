// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	casebytes "github.com/gofiber/utils/v2/bytes"
)

// ToLowerBytes converts an ASCII byte slice to lower-case in-place.
//
// Deprecated: use [github.com/gofiber/utils/v2/bytes.UnsafeToLower] instead.
// This wrapper keeps backward compatibility by mutating the provided slice.
//
//go:fix inline
func ToLowerBytes(b []byte) []byte {
	return casebytes.UnsafeToLower(b)
}

// ToUpperBytes converts an ASCII byte slice to upper-case in-place.
//
// Deprecated: use [github.com/gofiber/utils/v2/bytes.UnsafeToUpper] instead.
// This wrapper keeps backward compatibility by mutating the provided slice.
//
//go:fix inline
func ToUpperBytes(b []byte) []byte {
	return casebytes.UnsafeToUpper(b)
}

// AddTrailingSlashBytes appends a trailing '/' to b if it does not already end with one.
// If the input already ends with '/', the original slice is returned.
// A new slice is returned when a '/' is appended, an empty b included, and the caller
// may write to it. The original slice is never modified.
func AddTrailingSlashBytes(b []byte) []byte {
	n := len(b)
	if n > 0 && b[n-1] == '/' {
		return b
	}
	// Not AddTrailingSlashString: for an empty input that returns the constant "/",
	// and a slice over it would fault on the first write.
	out := make([]byte, n+1)
	copy(out, b)
	out[n] = '/'
	return out
}
