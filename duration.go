package utils

import (
	"time"
)

const durationBufLen = 32 // the longest rendering is "-2562047h47m16.854775808s"

// AppendDuration appends d formatted exactly like time.Duration.String to dst
// and returns the extended slice, without the string allocation.
func AppendDuration(dst []byte, d time.Duration) []byte {
	var buf [durationBufLen]byte
	w := durationToBuf(&buf, d)
	return append(dst, buf[w:]...)
}

// durationToBuf renders d right-aligned into buf and returns the first offset.
func durationToBuf(buf *[durationBufLen]byte, d time.Duration) int {
	w := len(buf)
	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Sub-second: ns, µs (two UTF-8 bytes), or ms.
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			buf[w] = '0'
			return w
		case u < uint64(time.Microsecond):
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			prec = 3
			w--
			buf[w] = 0xC2 // "µ" is 0xC2 0xB5
			buf[w+1] = 0xB5
		default:
			prec = 6
			buf[w] = 'm'
		}
		w, u = durationFrac(buf, w, u, prec)
		w = durationInt(buf, w, u)
	} else {
		w--
		buf[w] = 's'
		w, u = durationFrac(buf, w, u, 9)
		// u is now whole seconds.
		w = durationInt(buf, w, u%60)
		u /= 60
		if u > 0 {
			w--
			buf[w] = 'm'
			w = durationInt(buf, w, u%60)
			u /= 60
			if u > 0 {
				w--
				buf[w] = 'h'
				w = durationInt(buf, w, u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}
	return w
}

// durationFrac writes v/10**prec's fraction ending at buf[w] without trailing
// zeros and returns the new offset with v/10**prec.
func durationFrac(buf *[durationBufLen]byte, w int, v uint64, prec int) (int, uint64) {
	started := false
	for range prec {
		digit := v % 10
		v /= 10
		started = started || digit != 0
		if started {
			w--
			buf[w] = byte(digit) + '0'
		}
	}
	if started {
		w--
		buf[w] = '.'
	}
	return w, v
}

// durationInt writes the digits of v ending at buf[w] and returns the new offset.
func durationInt(buf *[durationBufLen]byte, w int, v uint64) int {
	for v >= 100 {
		q := v / 100
		r := v - q*100
		w -= 2
		buf[w] = decimalPairs[2*r]
		buf[w+1] = decimalPairs[2*r+1]
		v = q
	}
	if v >= 10 {
		w -= 2
		buf[w] = decimalPairs[2*v]
		buf[w+1] = decimalPairs[2*v+1]
		return w
	}
	w--
	buf[w] = byte(v) + '0'
	return w
}
