// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLower converts ascii string to lower-case
func ToLower(b string) string {
	n := len(b)
	if n == 0 {
		return b
	}

	// Find first character that needs conversion
	i := 0
	for ; i < n; i++ {
		if toLowerTable[b[i]] != b[i] {
			break
		}
	}

	// Already lowercase, return original
	if i == n {
		return b
	}

	// Allocate once and copy entire string, then modify in place
	res := make([]byte, n)
	copy(res, b)
	res[i] = toLowerTable[b[i]]

	// Process remaining characters with loop unrolling for better performance
	j := i + 1
	for ; j+3 < n; j += 4 {
		res[j] = toLowerTable[b[j]]
		res[j+1] = toLowerTable[b[j+1]]
		res[j+2] = toLowerTable[b[j+2]]
		res[j+3] = toLowerTable[b[j+3]]
	}
	for ; j < n; j++ {
		res[j] = toLowerTable[b[j]]
	}

	return UnsafeString(res)
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	n := len(b)
	if n == 0 {
		return b
	}

	// Find first character that needs conversion
	i := 0
	for ; i < n; i++ {
		if toUpperTable[b[i]] != b[i] {
			break
		}
	}

	// Already uppercase, return original
	if i == n {
		return b
	}

	// Allocate once and copy entire string, then modify in place
	res := make([]byte, n)
	copy(res, b)
	res[i] = toUpperTable[b[i]]

	// Process remaining characters with loop unrolling for better performance
	j := i + 1
	for ; j+3 < n; j += 4 {
		res[j] = toUpperTable[b[j]]
		res[j+1] = toUpperTable[b[j+1]]
		res[j+2] = toUpperTable[b[j+2]]
		res[j+3] = toUpperTable[b[j+3]]
	}
	for ; j < n; j++ {
		res[j] = toUpperTable[b[j]]
	}

	return UnsafeString(res)
}
