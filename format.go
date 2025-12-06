package utils

// smallInts contains precomputed string representations for small integers 0-99
var smallInts [100]string

// smallNegInts contains precomputed string representations for small negative integers -1 to -99
var smallNegInts [100]string

func init() {
	for i := range 100 {
		smallInts[i] = formatUintSmall(uint64(i))
		if i > 0 {
			smallNegInts[i] = "-" + smallInts[i]
		}
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
	if n >= 0 && n < 100 {
		return smallInts[n]
	}
	if n < 0 && n > -100 {
		return smallNegInts[-n]
	}
	if n >= 0 {
		var buf [20]byte
		i := formatUintBuf(&buf, uint64(n))
		return string(buf[i:])
	}
	var buf [20]byte
	i := formatUintBuf(&buf, uint64(-n))
	i--
	buf[i] = '-'
	return string(buf[i:])
}

// FormatUint32 formats a uint32 as a decimal string.
func FormatUint32(n uint32) string {
	if n < 100 {
		return smallInts[n]
	}
	var buf [10]byte // max 4294967295
	i := 10
	for n >= 10 {
		i--
		buf[i] = byte(n%10) + '0' //nolint:gosec // i is always in bounds: starts at 10, decrements max 10 times for uint32
		n /= 10
	}
	i--
	buf[i] = byte(n) + '0' //nolint:gosec // i is always >= 0 after loop
	return string(buf[i:])
}

// FormatInt32 formats an int32 as a decimal string.
func FormatInt32(n int32) string {
	if n >= 0 && n < 100 {
		return smallInts[n]
	}
	if n < 0 && n > -100 {
		return smallNegInts[-n]
	}
	if n >= 0 {
		return FormatUint32(uint32(n))
	}
	var buf [11]byte // max -2147483648
	un := uint32(-n)
	i := 11
	for un >= 10 {
		i--
		buf[i] = byte(un%10) + '0'
		un /= 10
	}
	i--
	buf[i] = byte(un) + '0'
	i--
	buf[i] = '-'
	return string(buf[i:])
}

// FormatUint16 formats a uint16 as a decimal string.
func FormatUint16(n uint16) string {
	if n < 100 {
		return smallInts[n]
	}
	var buf [5]byte // max 65535
	i := 5
	for n >= 10 {
		i--
		buf[i] = byte(n%10) + '0' //nolint:gosec // i is always in bounds: starts at 5, decrements max 5 times for uint16
		n /= 10
	}
	i--
	buf[i] = byte(n) + '0' //nolint:gosec // i is always >= 0 after loop
	return string(buf[i:])
}

// FormatInt16 formats an int16 as a decimal string.
func FormatInt16(n int16) string {
	if n >= 0 && n < 100 {
		return smallInts[n]
	}
	if n < 0 && n > -100 {
		return smallNegInts[-n]
	}
	if n >= 0 {
		return FormatUint16(uint16(n))
	}
	var buf [6]byte // max -32768
	un := uint16(-n)
	i := 6
	for un >= 10 {
		i--
		buf[i] = byte(un%10) + '0' //nolint:gosec // i is always in bounds
		un /= 10
	}
	i--
	buf[i] = byte(un) + '0' //nolint:gosec // i is always >= 1 after loop
	i--
	buf[i] = '-' //nolint:gosec // i is always >= 0 after decrement
	return string(buf[i:])
}

// FormatUint8 formats a uint8 as a decimal string.
func FormatUint8(n uint8) string {
	if n < 100 {
		return smallInts[n]
	}
	// uint8 max is 255, so max 3 digits
	return string([]byte{n/100 + '0', (n/10)%10 + '0', n%10 + '0'})
}

// FormatInt8 formats an int8 as a decimal string.
func FormatInt8(n int8) string {
	if n >= 0 && n < 100 {
		return smallInts[n]
	}
	if n < 0 && n > -100 {
		return smallNegInts[-n]
	}
	// Only -128 to -100 and 100 to 127 reach here
	if n >= 0 {
		un := uint8(n)
		return string([]byte{un/100 + '0', (un/10)%10 + '0', un%10 + '0'})
	}
	// n is -128 to -100
	un := uint8(-n)
	return string([]byte{'-', un/100 + '0', (un/10)%10 + '0', un%10 + '0'})
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
	if n > -100 {
		return append(dst, smallNegInts[-n]...)
	}
	var buf [20]byte
	i := formatUintBuf(&buf, uint64(-n))
	i--
	buf[i] = '-' //nolint:gosec // i is always >= 0: formatUintBuf returns at least 1 for any input, so i >= 0 after decrement
	return append(dst, buf[i:]...)
}
