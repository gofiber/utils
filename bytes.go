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
// A new slice is returned when a '/' is appended. The original slice is never modified.
func AddTrailingSlashBytes(b []byte) []byte {
	return UnsafeBytes(AddTrailingSlashString(UnsafeString(b)))
}
