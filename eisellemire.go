// The Eisel-Lemire algorithm in this file follows the implementation in Go's
// strconv package, which carries this license:
//
// Copyright 2020 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google LLC nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package utils

import (
	"math"
	"math/bits"
)

//go:generate go run pow10gen.go

// eiselLemire returns man * 10^exp10 rounded to the nearest float64 (float32 when
// bitSize is 32, held in a float64), or false when it cannot tell how to round: the
// exponent is outside the table, the result is subnormal or overflows, or the product
// is too close to a float or a halfway point for the 128-bit approximation.
// man must not be zero. See https://nigeltao.github.io/blog/2020/eisel-lemire.html.
func eiselLemire(man uint64, exp10 int64, bitSize int) (float64, bool) {
	if exp10 < pow10Min || exp10 > pow10Max {
		return 0, false
	}
	mantBits, bias := uint(52), uint64(1023)
	if bitSize == 32 {
		mantBits, bias = 23, 127
	}
	// Keep mantBits+2 bits of the product (one for rounding); the lowBits below
	// them tell whether the approximation is decisive.
	lowBits := 61 - mantBits
	lowMask := uint64(1)<<lowBits - 1

	// Normalize man and estimate the binary exponent: 217706/2^16 is log2(10).
	clz := uint(bits.LeadingZeros64(man))
	man <<= clz
	exp2 := uint64(217706*exp10>>16+64) + bias - uint64(clz)

	// Multiply by the top 64 bits of the power; if the low bits of the product are
	// all ones, the dropped bits could carry into the result, so use them too.
	pow := &pow10Tab[exp10-pow10Min]
	hi, lo := bits.Mul64(man, pow.hi)
	if hi&lowMask == lowMask && lo+man < man {
		yHi, yLo := bits.Mul64(man, pow.lo)
		mergedHi, mergedLo := hi, lo+yHi
		if mergedLo < lo {
			mergedHi++
		}
		if mergedHi&lowMask == lowMask && mergedLo+1 == 0 && yLo+man < man {
			return 0, false
		}
		hi, lo = mergedHi, mergedLo
	}

	// The product has 127 or 128 bits, so the shift and the exponent depend on the top one.
	msb := hi >> 63
	m := hi >> (msb + uint64(lowBits))
	exp2 -= 1 ^ msb

	// With the rounding bit set and an even mantissa, an exact tie rounds down and
	// anything above rounds up; the approximation cannot tell which this is.
	if lo == 0 && hi&lowMask == 0 && m&3 == 1 {
		return 0, false
	}

	// Drop the rounding bit, rounding up if it is set.
	m = (m + 1) >> 1
	if m>>(mantBits+1) > 0 {
		m >>= 1
		exp2++
	}

	// exp2 is unsigned, so zero or below wrapped; 2*bias+1 or more is Inf or NaN.
	if exp2-1 >= 2*bias {
		return 0, false
	}
	b := exp2<<mantBits | m&(1<<mantBits-1)
	if bitSize == 32 {
		return float64(math.Float32frombits(uint32(b))), true
	}
	return math.Float64frombits(b), true
}
