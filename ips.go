package utils

import (
	"net"
)

// IsIPv4 works the same way as net.ParseIP,
// but without check for IPv6 case and without returning net.IP slice, whereby IsIPv4 makes no allocations.
func IsIPv4(s string) bool {
	l := len(s)
	if l < 7 || l > 15 {
		return false
	}
	for i := 0; i < net.IPv4len; i++ {
		if l == 0 {
			return false
		}

		if i > 0 {
			if s[0] != '.' {
				return false
			}
			s = s[1:]
			l--
		}

		n, ci := 0, 0

		for ci = 0; ci < l && '0' <= s[ci] && s[ci] <= '9'; ci++ {
			n = n*10 + int(s[ci]-'0')
			if n > 0xFF {
				return false
			}
		}

		if ci == 0 || (ci > 1 && s[0] == '0') {
			return false
		}

		s = s[ci:]
		l -= ci
	}

	return l == 0
}

// IsIPv6 works the same way as net.ParseIP,
// but without check for IPv4 case and without returning net.IP slice, whereby IsIPv6 makes no allocations.
func IsIPv6(s string) bool {
	l := len(s)
	if l < 2 || l > 39 {
		return false
	}
	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if l >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		l -= 2
		// Might be only ellipsis
		if l == 0 {
			return true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i := 0
	for i < net.IPv6len {
		// Hex number.
		n, ci := 0, 0

		for ci = 0; ci < l; ci++ {
			if '0' <= s[ci] && s[ci] <= '9' {
				n *= 16
				n += int(s[ci] - '0')
			} else if 'a' <= s[ci] && s[ci] <= 'f' {
				n *= 16
				n += int(s[ci]-'a') + 10
			} else if 'A' <= s[ci] && s[ci] <= 'F' {
				n *= 16
				n += int(s[ci]-'A') + 10
			} else {
				break
			}
			if n > 0xFFFF {
				return false
			}
		}
		if ci == 0 || n > 0xFFFF {
			return false
		}

		if ci < l && s[ci] == '.' {
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
			l = 0
			i += net.IPv4len
			break
		}

		// Save this 16-bit chunk.
		i += 2

		// Stop at end of string.
		s = s[ci:]
		l -= ci
		if l == 0 {
			break
		}

		// Otherwise must be followed by colon and more.
		if s[0] != ':' || l == 1 {
			return false
		}
		s = s[1:]
		l--

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				return false
			}
			ellipsis = i
			s = s[1:]
			l--
			if l == 0 { // can be at end
				break
			}
		}
	}

	// Must have used entire string.
	if l != 0 {
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
