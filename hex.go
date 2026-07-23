package utils

import (
	"slices"

	"github.com/gofiber/utils/v2/swar"
)

// lowerHexDigits is the encoding alphabet used by AppendHexEncode, matching
// encoding/hex.
const lowerHexDigits = "0123456789abcdef"

// AppendHexEncode appends the lowercase hexadecimal encoding of src to dst
// and returns the extended slice, byte-identical to encoding/hex.AppendEncode.
// dst must not alias src: the output is longer than the input, so in-place
// encoding is impossible. Eight source bytes are expanded per SWAR word pair;
// the remainder is table-driven.
//
// There is deliberately no AppendHexDecode counterpart: encoding/hex's
// decoder is already table-driven and append-capable (hex.AppendDecode), and
// a SWAR kernel measured slower than it, so decoding stays with the stdlib.
func AppendHexEncode[S byteSeq](dst []byte, src S) []byte {
	n := len(src)
	if n == 0 {
		return dst
	}
	off := len(dst)
	dst = slices.Grow(dst, 2*n)[: off+2*n : off+2*n]
	out := dst[off:]
	i := 0
	// One 8-byte load feeds two independent expansions, so the halves
	// pipeline and each iteration retires 16 output bytes.
	for ; i+8 <= n; i += 8 {
		x := swar.Load8(src, i)
		swar.Store8(out, 2*i, hexEncodeWord(x&halfWord))
		swar.Store8(out, 2*i+8, hexEncodeWord(x>>32))
	}
	if i+4 <= n {
		swar.Store8(out, 2*i, hexEncodeWord(load4(src, i)))
		i += 4
	}
	for ; i < n; i++ {
		c := src[i]
		out[2*i] = lowerHexDigits[c>>4]
		out[2*i+1] = lowerHexDigits[c&0x0F]
	}
	return dst
}

// halfWord keeps the low four lanes of a SWAR word.
const halfWord = 0x00000000FFFFFFFF

// load4 assembles src[i:i+4] into the low four lanes of a uint64 in Load8's
// little-endian lane order. The caller must guarantee i+4 <= len(src).
func load4[S byteSeq](src S, i int) uint64 {
	w := src[i : i+4]
	return uint64(w[0]) | uint64(w[1])<<8 | uint64(w[2])<<16 | uint64(w[3])<<24
}

// hexEncodeWord expands the four bytes in the low lanes of x into eight
// lowercase hex digit lanes: lane 2k holds the high nibble of byte k and
// lane 2k+1 its low nibble, so a Store8 writes the digits in encoding/hex
// order.
func hexEncodeWord(x uint64) uint64 {
	const (
		evenPairs   = 0x0000FFFF0000FFFF // 16-bit groups 0 and 2
		evenBytes   = 0x00FF00FF00FF00FF // byte lanes 0, 2, 4, 6
		evenNibbles = 0x000F000F000F000F // low nibble of lanes 0, 2, 4, 6
		sixes       = 0x0606060606060606
		tens        = 0x1010101010101010
		zeros       = 0x3030303030303030 // '0' in every lane
		letterGap   = 0x27               // 'a' - '0' - 10
	)
	// Spread bytes 0..3 into lanes 0, 2, 4, 6; the two masked doublings
	// cannot collide because each step keeps only the copies that landed in
	// still-empty positions.
	v := (x | x<<16) & evenPairs
	v = (v | v<<8) & evenBytes
	// Split into nibble lanes: high nibbles stay in the even lanes, low
	// nibbles move to the odd lanes.
	nib := (v>>4)&evenNibbles | (v&evenNibbles)<<8
	// Branchless digit/letter selection: adding 6 carries into bit 4 exactly
	// for nibbles > 9, and each resulting 0x01 lane times letterGap bridges
	// the ASCII gap between '9'+1 and 'a'. The multiply cannot carry across
	// lanes because every lane factor is 0 or 1.
	return nib + zeros + ((nib+sixes)&tens)>>4*letterGap
}
