// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLower converts ascii string to lower-case
func ToLower(b string) string {
	if len(b) == 0 {
		return b
	}

	for i := 0; i < len(b); i++ {
		c := b[i]
		low := toLowerTable[c]
		if low != c {
			res := make([]byte, len(b))
			copy(res, b[:i])
			res[i] = low
			for j := i + 1; j < len(b); j++ {
				res[j] = toLowerTable[b[j]]
			}
			return UnsafeString(res)
		}
	}
	return b
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	if len(b) == 0 {
		return b
	}

	for i := 0; i < len(b); i++ {
		c := b[i]
		up := toUpperTable[c]
		if up != c {
			res := make([]byte, len(b))
			copy(res, b[:i])
			res[i] = up
			for j := i + 1; j < len(b); j++ {
				res[j] = toUpperTable[b[j]]
			}
			return UnsafeString(res)
		}
	}

	return b
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
