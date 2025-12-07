// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

// ToLowerBytes converts ascii slice to lower-case
func ToLowerBytes(b []byte) []byte {
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

// ToUpperBytes converts ascii slice to upper-case
func ToUpperBytes(b []byte) []byte {
	table := toUpperTable
	n := len(b)
	i := 0

	// Unroll by 4 to match ToLowerBytes and maximize throughput on amd64.
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
