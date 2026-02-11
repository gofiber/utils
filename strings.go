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

	table := toLowerTable
	for i := range n {
		c := b[i]
		low := table[c]
		if low != c {
			res := make([]byte, n)
			copy(res, b[:i])
			res[i] = low
			j := i + 1
			for ; j+3 < n; j += 4 {
				res[j+0] = table[b[j+0]]
				res[j+1] = table[b[j+1]]
				res[j+2] = table[b[j+2]]
				res[j+3] = table[b[j+3]]
			}
			for ; j < n; j++ {
				res[j] = table[b[j]]
			}
			return UnsafeString(res)
		}
	}
	return b
}

// ToLowerMut converts an ASCII string to lower-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func ToLowerMut(b string) string {
	ToLowerBytesMut(UnsafeBytes(b))
	return b
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	n := len(b)
	if n == 0 {
		return b
	}

	table := toUpperTable
	for i := range n {
		c := b[i]
		up := table[c]
		if up != c {
			res := make([]byte, n)
			copy(res, b[:i])
			res[i] = up
			j := i + 1
			for ; j+3 < n; j += 4 {
				res[j+0] = table[b[j+0]]
				res[j+1] = table[b[j+1]]
				res[j+2] = table[b[j+2]]
				res[j+3] = table[b[j+3]]
			}
			for ; j < n; j++ {
				res[j] = table[b[j]]
			}
			return UnsafeString(res)
		}
	}

	return b
}

// ToUpperMut converts an ASCII string to upper-case by mutating its backing bytes in-place.
// This function is unsafe: it breaks string immutability and must only be used when the
// string is known to reference mutable memory.
func ToUpperMut(b string) string {
	ToUpperBytesMut(UnsafeBytes(b))
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
