package utils

func mapASCII(dst, src []byte, table string) {
	n := len(src)
	i := 0
	// Unroll by 4 to balance instruction-level parallelism with cache pressure.
	for ; i+3 < n; i += 4 {
		dst[i+0] = table[src[i+0]]
		dst[i+1] = table[src[i+1]]
		dst[i+2] = table[src[i+2]]
		dst[i+3] = table[src[i+3]]
	}

	for ; i < n; i++ {
		dst[i] = table[src[i]]
	}
}

func mapASCIIInPlace(b []byte, table string) {
	mapASCII(b, b, table)
}

func mapASCIICopy(src []byte, table string) []byte {
	n := len(src)
	if n == 0 {
		if src == nil {
			return nil
		}
		return make([]byte, 0)
	}

	dst := make([]byte, n)
	mapASCII(dst, src, table)
	return dst
}

func mapASCIIString(s, table string) string {
	n := len(s)
	if n == 0 {
		return s
	}

	src := s
	i := 0
	for ; i < n; i++ {
		c := src[i]
		mapped := table[c]
		if mapped != c {
			dst := make([]byte, n)
			copy(dst, src[:i])
			dst[i] = mapped
			i++
			// Unroll by 4 to balance instruction-level parallelism with cache pressure.
			for ; i+3 < n; i += 4 {
				dst[i+0] = table[src[i+0]]
				dst[i+1] = table[src[i+1]]
				dst[i+2] = table[src[i+2]]
				dst[i+3] = table[src[i+3]]
			}
			for ; i < n; i++ {
				dst[i] = table[src[i]]
			}
			return UnsafeString(dst)
		}
	}

	return s
}

func mapASCIIStringMut(s, table string) string {
	mapASCIIInPlace(UnsafeBytes(s), table)
	return s
}
