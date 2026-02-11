// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLowerBytes converts an ASCII byte slice to lower-case without modifying the input.
func ToLowerBytes(b []byte) []byte {
	return mapASCIICopy(b, toLowerTable)
}

// ToLowerBytesMut converts an ASCII byte slice to lower-case in-place.
// The passed slice content is modified and the same slice is returned.
func ToLowerBytesMut(b []byte) []byte {
	mapASCIIInPlace(b, toLowerTable)
	return b
}

// ToUpperBytes converts an ASCII byte slice to upper-case without modifying the input.
func ToUpperBytes(b []byte) []byte {
	return mapASCIICopy(b, toUpperTable)
}

// ToUpperBytesMut converts an ASCII byte slice to upper-case in-place.
// The passed slice content is modified and the same slice is returned.
func ToUpperBytesMut(b []byte) []byte {
	mapASCIIInPlace(b, toUpperTable)
	return b
}

// AddTrailingSlashBytes appends a trailing '/' to b if it does not already end with one.
// If the input already ends with '/', the original slice is returned.
// A new slice is returned when a '/' is appended. The original slice is never modified.
func AddTrailingSlashBytes(b []byte) []byte {
	return UnsafeBytes(AddTrailingSlashString(UnsafeString(b)))
}
