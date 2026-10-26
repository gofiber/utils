package utils

import (
	"encoding/binary"
	"math/bits"
)

var smallInts [100]string

var smallNegInts [100]string

var uint8Strs [256]string

var int8Strs [256]string

func init() {
	for i := range 100 {
		smallInts[i] = formatUintSmall(uint64(i))
		if i > 0 {
			smallNegInts[i] = "-" + smallInts[i]
		}
	}

	for i := range 256 {
		v := uint8(i)
		uint8Strs[i] = formatUint8Slow(v)

		sv := int8(i)
		if sv >= 0 {
			int8Strs[i] = uint8Strs[sv]
		} else {
			int8Strs[i] = "-" + uint8Strs[uint8(-sv)]
		}
	}
}

func formatUintSmall(n uint64) string {
	if n < 10 {
		return string(byte(n) + '0')
	}
	return string([]byte{byte(n/10) + '0', byte(n%10) + '0'})
}

func formatUint8Slow(n uint8) string {
	if n < 100 {
		return smallInts[n]
	}
	return string([]byte{n/100 + '0', (n/10)%10 + '0', n%10 + '0'})
}

const decimalPairs = "00010203040506070809101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899"

const (
	digitZeros   = 0x3030303030303030 // '0' per lane; lanes are below 16, so OR equals add
	digitsBufLen = 24                 // three 8-digit groups: 20 digits plus a sign
	digitsGroup  = 100000000
)

// digits8 returns n < 1e8 as eight zero-padded digit lanes; no step carries across lanes.
func digits8(n uint32) uint64 {
	hi := n / 10000
	x := uint64(hi) | uint64(n-hi*10000)<<32
	q := (x * 5243 >> 19) & 0x0000007F0000007F // lane/100, exact below 43699
	y := q | (x-q*100)<<16
	q = (y * 103 >> 10) & 0x000F000F000F000F // lane/10, exact below 100
	return q | (y-q*10)<<8
}

// uintToBuf right-aligns n >= 100 in buf and returns the index of its first digit.
func uintToBuf(buf *[digitsBufLen]byte, n uint64) int {
	if n < digitsGroup {
		z := digits8(uint32(n))
		binary.LittleEndian.PutUint64(buf[16:24], z|digitZeros)
		return 16 + bits.TrailingZeros64(z)>>3
	}
	hi := n / digitsGroup
	binary.LittleEndian.PutUint64(buf[16:24], digits8(uint32(n-hi*digitsGroup))|digitZeros)
	if hi < digitsGroup {
		z := digits8(uint32(hi))
		binary.LittleEndian.PutUint64(buf[8:16], z|digitZeros)
		return 8 + bits.TrailingZeros64(z)>>3
	}
	top := uint32(hi / digitsGroup) // at most 1844
	binary.LittleEndian.PutUint64(buf[8:16], digits8(uint32(hi-uint64(top)*digitsGroup))|digitZeros)
	z := digits8(top)
	binary.LittleEndian.PutUint64(buf[0:8], z|digitZeros)
	return bits.TrailingZeros64(z) >> 3
}

func uint32ToBuf(buf *[16]byte, n uint32) int {
	if n < digitsGroup {
		z := digits8(n)
		binary.LittleEndian.PutUint64(buf[8:16], z|digitZeros)
		return 8 + bits.TrailingZeros64(z)>>3
	}
	hi := n / digitsGroup
	binary.LittleEndian.PutUint64(buf[8:16], digits8(n-hi*digitsGroup)|digitZeros)
	z := digits8(hi)
	binary.LittleEndian.PutUint64(buf[0:8], z|digitZeros)
	return bits.TrailingZeros64(z) >> 3
}

// FormatUint formats n in decimal.
func FormatUint(n uint64) string {
	if n < 100 {
		return smallInts[n]
	}
	var buf [digitsBufLen]byte
	i := uintToBuf(&buf, n)
	return string(buf[i:])
}

// FormatInt formats n in decimal.
func FormatInt(n int64) string {
	if n >= 0 {
		if n < 100 {
			return smallInts[n]
		}
		return FormatUint(uint64(n))
	}
	if n > -100 {
		return smallNegInts[-n]
	}
	var buf [digitsBufLen]byte
	i := uintToBuf(&buf, uint64(-n)) - 1 // uint64(-n) is the magnitude, math.MinInt64 included
	buf[i] = '-'
	return string(buf[i:])
}

// FormatUint32 formats n in decimal.
func FormatUint32(n uint32) string {
	if n < 100 {
		return smallInts[n]
	}
	var buf [16]byte
	i := uint32ToBuf(&buf, n)
	return string(buf[i:])
}

// FormatInt32 formats n in decimal.
func FormatInt32(n int32) string {
	if n >= 0 {
		if n < 100 {
			return smallInts[n]
		}
		return FormatUint32(uint32(n))
	}
	if n > -100 {
		return smallNegInts[-n]
	}
	var buf [16]byte
	i := uint32ToBuf(&buf, uint32(-n)) - 1
	buf[i] = '-'
	return string(buf[i:])
}

// FormatUint16 formats n in decimal.
func FormatUint16(n uint16) string {
	return FormatUint32(uint32(n))
}

// FormatInt16 formats n in decimal.
func FormatInt16(n int16) string {
	return FormatInt32(int32(n))
}

// FormatUint8 formats n in decimal.
func FormatUint8(n uint8) string {
	return uint8Strs[n]
}

// FormatInt8 formats n in decimal.
func FormatInt8(n int8) string {
	return int8Strs[uint8(n)]
}

var pow10 = [20]uint64{
	1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
}

func uintDigits(n uint64) int {
	if n < 10 {
		return 1
	}
	// floor(log10(n)) from floor(log2(n)): 1233/4096 ~= log10(2).
	t := (bits.Len64(n) * 1233) >> 12
	if n < pow10[t] {
		return t
	}
	return t + 1
}

// intDigits returns the length of n's decimal string, including a '-' sign.
func intDigits(n int64) int {
	if n < 0 {
		return 1 + uintDigits(uint64(-n))
	}
	return uintDigits(uint64(n))
}

// AppendUint appends FormatUint(n) to dst.
func AppendUint(dst []byte, n uint64) []byte {
	if n < 100 {
		return append(dst, smallInts[n]...)
	}
	var buf [digitsBufLen]byte
	i := uintToBuf(&buf, n)
	return append(dst, buf[i:]...)
}

// AppendInt appends FormatInt(n) to dst.
func AppendInt(dst []byte, n int64) []byte {
	if n >= 0 {
		if n < 100 {
			return append(dst, smallInts[n]...)
		}
		return AppendUint(dst, uint64(n))
	}
	if n > -100 {
		return append(dst, smallNegInts[-n]...)
	}
	var buf [digitsBufLen]byte
	i := uintToBuf(&buf, uint64(-n)) - 1
	buf[i] = '-'
	return append(dst, buf[i:]...)
}
