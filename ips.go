package utils

import (
	"net"
)

var hexTable = [256]byte{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

// IsIPv4 works the same way as net.ParseIP,
// but without check for IPv6 case and without returning net.IP slice, whereby IsIPv4 makes no allocations.
func IsIPv4(s string) bool {
	for i := range net.IPv4len {
		if len(s) == 0 {
			return false
		}

		if i > 0 {
			if s[0] != '.' {
				return false
			}
			s = s[1:]
		}

		n, ci := 0, 0

		for ci = 0; ci < len(s) && '0' <= s[ci] && s[ci] <= '9'; ci++ {
			n = n*10 + int(s[ci]-'0')
			if n > 0xFF {
				return false
			}
		}

		if ci == 0 || (ci > 1 && s[0] == '0') {
			return false
		}

		s = s[ci:]
	}

	return len(s) == 0
}

// IsIPv6 works the same way as net.ParseIP,
// but without check for IPv4 case and without returning net.IP slice, whereby IsIPv6 makes no allocations.
func IsIPv6(s string) bool {
	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		// Might be only ellipsis
		if len(s) == 0 {
			return true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i := 0
	for i < net.IPv6len {
		// Hex number.
		n, ci := 0, 0

		for ci = 0; ci < len(s); ci++ {
			if (s[ci] < '0' || s[ci] > '9') && (s[ci] < 'a' || s[ci] > 'f') && (s[ci] < 'A' || s[ci] > 'F') {
				break
			}
			n *= 16
			n += int(hexTable[s[ci]])

			if n > 0xFFFF {
				return false
			}
		}
		if ci == 0 || n > 0xFFFF {
			return false
		}

		if ci < len(s) && s[ci] == '.' {
			if ellipsis < 0 && i != net.IPv6len-net.IPv4len {
				return false
			}
			if i+net.IPv4len > net.IPv6len {
				return false
			}

			if !IsIPv4(s) {
				return false
			}

			s = ""
			i += net.IPv4len
			break
		}

		// Save this 16-bit chunk.
		i += 2

		// Stop at end of string.
		s = s[ci:]
		if len(s) == 0 {
			break
		}

		// Otherwise must be followed by colon and more.
		if s[0] != ':' || len(s) == 1 {
			return false
		}
		s = s[1:]

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				return false
			}
			ellipsis = i
			s = s[1:]
			if len(s) == 0 { // can be at end
				break
			}
		}
	}

	// Must have used entire string.
	if len(s) != 0 {
		return false
	}

	// If didn't parse enough, expand ellipsis.
	if i < net.IPv6len {
		if ellipsis < 0 {
			return false
		}
	} else if ellipsis >= 0 {
		// Ellipsis must represent at least one 0 group.
		return false
	}
	return true
}
