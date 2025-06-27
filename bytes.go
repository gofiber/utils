// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLowerBytes converts ascii slice to lower-case
func ToLowerBytes(b []byte) []byte {
	for i := range b {
		b[i] = toLowerTable[b[i]]
	}
	return b
}

// ToUpperBytes converts ascii slice to upper-case
func ToUpperBytes(b []byte) []byte {
	for i := range b {
		b[i] = toUpperTable[b[i]]
	}
	return b
}
