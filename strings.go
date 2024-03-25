// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLower converts ascii string to lower-case
func ToLower(b string) string {
	res := make([]byte, len(b))
	copy(res, b)
	for i := 0; i < len(res); i++ {
		res[i] = toLowerTable[res[i]]
	}

	return UnsafeString(res)
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	res := make([]byte, len(b))
	copy(res, b)
	for i := 0; i < len(res); i++ {
		res[i] = toUpperTable[res[i]]
	}

	return UnsafeString(res)
}

// IfToLower returns an lowercase version of the input ASCII string.
//
// It first checks if the string contains any uppercase characters before converting it.
//
// For strings that are already lowercase,this function will be faster than `ToLower`.
//
// In the case of mixed-case or uppercase strings, this function will be slightly slower than `ToLower`.
func IfToLower(s string) string {
	hasUpper := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if toLowerTable[c] != c {
			hasUpper = true
			break
		}
	}

	if !hasUpper {
		return s
	}
	return ToLower(s)
}

// IfToUpper returns an uppercase version of the input ASCII string.
//
// It first checks if the string contains any lowercase characters before converting it.
//
// For strings that are already uppercase,this function will be faster than `ToUpper`.
//
// In the case of mixed-case or lowercase strings, this function will be slightly slower than `ToUpper`.
func IfToUpper(s string) string {
	hasLower := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if toUpperTable[c] != c {
			hasLower = true
			break
		}
	}

	if !hasLower {
		return s
	}
	return ToUpper(s)
}
