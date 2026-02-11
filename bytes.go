// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLowerBytes converts an ASCII byte slice to lower-case without modifying the input.
func ToLowerBytes(b []byte) []byte {
	n := len(b)
	if n == 0 {
		return b
	}

	table := toLowerTable
	dst := make([]byte, n)
	i := 0
	// Unroll by 4 to balance instruction-level parallelism with cache pressure.
	limit := n &^ 3
	for i < limit {
		dst[i+0] = table[b[i+0]]
		dst[i+1] = table[b[i+1]]
		dst[i+2] = table[b[i+2]]
		dst[i+3] = table[b[i+3]]
		i += 4
	}

	for i < n {
		dst[i] = table[b[i]]
		i++
	}

	return dst
}

// ToLowerBytesMut converts an ASCII byte slice to lower-case in-place.
// The passed slice content is modified and the same slice is returned.
func ToLowerBytesMut(b []byte) []byte {
	table := toLowerTable
	n := len(b)
	i := 0

	// Unroll by 4 to balance instruction-level parallelism with cache pressure.
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		b1 := b[i+1]
		b2 := b[i+2]
		b3 := b[i+3]

		b[i+0] = table[b0]
		b[i+1] = table[b1]
		b[i+2] = table[b2]
		b[i+3] = table[b3]

		i += 4
	}

	for i < n {
		b[i] = table[b[i]]
		i++
	}

	return b
}

// ToUpperBytes converts an ASCII byte slice to upper-case without modifying the input.
func ToUpperBytes(b []byte) []byte {
	n := len(b)
	if n == 0 {
		return b
	}

	table := toUpperTable
	dst := make([]byte, n)
	i := 0
	// Unroll by 4 to balance instruction-level parallelism with cache pressure.
	limit := n &^ 3
	for i < limit {
		dst[i+0] = table[b[i+0]]
		dst[i+1] = table[b[i+1]]
		dst[i+2] = table[b[i+2]]
		dst[i+3] = table[b[i+3]]
		i += 4
	}

	for i < n {
		dst[i] = table[b[i]]
		i++
	}

	return dst
}

// ToUpperBytesMut converts an ASCII byte slice to upper-case in-place.
// The passed slice content is modified and the same slice is returned.
func ToUpperBytesMut(b []byte) []byte {
	table := toUpperTable
	n := len(b)
	i := 0

	// Unroll by 4 to balance instruction-level parallelism with cache pressure.
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		b1 := b[i+1]
		b2 := b[i+2]
		b3 := b[i+3]

		b[i+0] = table[b0]
		b[i+1] = table[b1]
		b[i+2] = table[b2]
		b[i+3] = table[b3]

		i += 4
	}

	for i < n {
		b[i] = table[b[i]]
		i++
	}

	return b
}

// AddTrailingSlashBytes appends a trailing '/' to b if it does not already end with one.
// If the input already ends with '/', the original slice is returned.
// A new slice is returned when a '/' is appended. The original slice is never modified.
func AddTrailingSlashBytes(b []byte) []byte {
	return UnsafeBytes(AddTrailingSlashString(UnsafeString(b)))
}
