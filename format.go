package utils

// smallInts contains precomputed string representations for small integers 0-99
var smallInts [100]string

func init() {
	for i := range 100 {
		smallInts[i] = formatUintSmall(uint64(i))
	}
}

func formatUintSmall(n uint64) string {
	if n < 10 {
		return string(byte(n) + '0')
	}
	return string([]byte{byte(n/10) + '0', byte(n%10) + '0'})
}

// formatUintBuf writes the digits of n into buf from the end and returns the start index.
// buf must be at least 20 bytes.
func formatUintBuf(buf *[20]byte, n uint64) int {
	i := 20
	for n >= 10 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	i--
	buf[i] = byte(n) + '0'
	return i
}

// FormatUint formats a uint64 as a decimal string.
// It is faster than strconv.FormatUint for most inputs.
func FormatUint(n uint64) string {
	if n < 100 {
		return smallInts[n]
	}
	var buf [20]byte
	i := formatUintBuf(&buf, n)
	return string(buf[i:])
}

// FormatInt formats an int64 as a decimal string.
// It is faster than strconv.FormatInt for most inputs.
func FormatInt(n int64) string {
	if n >= 0 {
		return FormatUint(uint64(n))
	}
	var buf [20]byte
	i := formatUintBuf(&buf, uint64(-n))
	i--
	buf[i] = '-'
	return string(buf[i:])
}

// FormatUint32 formats a uint32 as a decimal string.
func FormatUint32(n uint32) string {
	return FormatUint(uint64(n))
}

// FormatInt32 formats an int32 as a decimal string.
func FormatInt32(n int32) string {
	return FormatInt(int64(n))
}

// FormatUint16 formats a uint16 as a decimal string.
func FormatUint16(n uint16) string {
	return FormatUint(uint64(n))
}

// FormatInt16 formats an int16 as a decimal string.
func FormatInt16(n int16) string {
	return FormatInt(int64(n))
}

// FormatUint8 formats a uint8 as a decimal string.
func FormatUint8(n uint8) string {
	return FormatUint(uint64(n))
}

// FormatInt8 formats an int8 as a decimal string.
func FormatInt8(n int8) string {
	return FormatInt(int64(n))
}

// AppendUint appends the decimal string representation of n to dst.
func AppendUint(dst []byte, n uint64) []byte {
	if n < 100 {
		return append(dst, smallInts[n]...)
	}
	var buf [20]byte
	i := formatUintBuf(&buf, n)
	return append(dst, buf[i:]...)
}

// AppendInt appends the decimal string representation of n to dst.
func AppendInt(dst []byte, n int64) []byte {
	if n >= 0 {
		return AppendUint(dst, uint64(n))
	}
	var buf [20]byte
	i := formatUintBuf(&buf, uint64(-n))
	i--
	buf[i] = '-'
	return append(dst, buf[i:]...)
}
