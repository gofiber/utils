package utils

type byteSeq interface {
	~string | ~[]byte
}

// EqualFold tests ascii strings or bytes for equality case-insensitively
func EqualFold[S byteSeq](b, s S) bool {
	if len(b) != len(s) {
		return false
	}

	table := toUpperTable
	n := len(b)
	i := 0

	// Unroll by 4 to match other hot paths and drive instruction-level parallelism.
	limit := n &^ 3
	for i < limit {
		b0 := b[i+0]
		s0 := s[i+0]
		if table[b0] != table[s0] {
			return false
		}

		b1 := b[i+1]
		s1 := s[i+1]
		if table[b1] != table[s1] {
			return false
		}

		b2 := b[i+2]
		s2 := s[i+2]
		if table[b2] != table[s2] {
			return false
		}

		b3 := b[i+3]
		s3 := s[i+3]
		if table[b3] != table[s3] {
			return false
		}

		i += 4
	}

	for i < n {
		if table[b[i]] != table[s[i]] {
			return false
		}
		i++
	}
	return true
}

// TrimLeft is the equivalent of strings/bytes.TrimLeft
func TrimLeft[S byteSeq](s S, cutset byte) S {
	lenStr, start := len(s), 0
	for start < lenStr && s[start] == cutset {
		start++
	}
	return s[start:]
}

// Trim is the equivalent of strings/bytes.Trim
func Trim[S byteSeq](s S, cutset byte) S {
	i, j := 0, len(s)-1
	for ; i <= j; i++ {
		if s[i] != cutset {
			break
		}
	}
	for ; i < j; j-- {
		if s[j] != cutset {
			break
		}
	}

	return s[i : j+1]
}

// TrimRight is the equivalent of strings/bytes.TrimRight
func TrimRight[S byteSeq](s S, cutset byte) S {
	lenStr := len(s)
	for lenStr > 0 && s[lenStr-1] == cutset {
		lenStr--
	}
	return s[:lenStr]
}

// AddTrailingSlash appends a trailing '/' to v if it does not already end with one.
//
// For string inputs, the result is always a new string.
//
// For []byte inputs, the result is always a new slice with the trailing '/' appended.
// The original slice is never modified, even if it has sufficient capacity.
//
// This function guarantees that the returned value is independent of the input.
func AddTrailingSlash[T byteSeq](v T) T {
	if len(v) > 0 && v[len(v)-1] == '/' {
		return v
	}

	switch x := any(v).(type) {
	case string:
		return any(x + "/").(T) //nolint:forcetypeassert,errcheck // can't fail
	case []byte:
		result := make([]byte, len(x)+1)
		copy(result, x)
		result[len(x)] = '/'
		return any(result).(T) //nolint:forcetypeassert,errcheck // can't fail
	default:
		// impossible because of the constraint
		return v
	}
}
