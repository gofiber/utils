// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	casestrings "github.com/gofiber/utils/v2/strings"
)

// ToLower converts an ASCII string to lower-case.
//
// Deprecated: use [github.com/gofiber/utils/v2/strings.ToLower] instead.
//
//go:fix inline
func ToLower(b string) string {
	return casestrings.ToLower(b)
}

// ToUpper converts an ASCII string to upper-case.
//
// Deprecated: use [github.com/gofiber/utils/v2/strings.ToUpper] instead.
//
//go:fix inline
func ToUpper(b string) string {
	return casestrings.ToUpper(b)
}

// AddTrailingSlashString returns s if it ends with '/', else a new string with '/' appended.
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
