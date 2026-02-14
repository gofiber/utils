// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	casebytes "github.com/gofiber/utils/v2/bytes"
)

// ToLowerBytes converts an ASCII byte slice to lower-case in-place.
//
// Deprecated: use package "github.com/gofiber/utils/v2/bytes" and call bytes.UnsafeToLower.
// This wrapper keeps backward compatibility by mutating the provided slice.
func ToLowerBytes(b []byte) []byte {
	return casebytes.UnsafeToLower(b)
}

// ToUpperBytes converts an ASCII byte slice to upper-case in-place.
//
// Deprecated: use package "github.com/gofiber/utils/v2/bytes" and call bytes.UnsafeToUpper.
// This wrapper keeps backward compatibility by mutating the provided slice.
func ToUpperBytes(b []byte) []byte {
	return casebytes.UnsafeToUpper(b)
}

// AddTrailingSlashBytes appends a trailing '/' to b if it does not already end with one.
// If the input already ends with '/', the original slice is returned.
// A new slice is returned when a '/' is appended. The original slice is never modified.
func AddTrailingSlashBytes(b []byte) []byte {
	return UnsafeBytes(AddTrailingSlashString(UnsafeString(b)))
}
