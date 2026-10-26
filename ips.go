package utils

import (
	"net"
)

const hexFieldMaxLen = 4

// IsIPv4 reports whether s is a dotted-decimal IPv4 address, as net.ParseIP accepts it.
// It does not allocate.
func IsIPv4(s string) bool {
	_, ok := ParseIPv4(s)
	return ok
}

// IsIPv6 reports whether s is an IPv6 address, as net.ParseIP accepts it, without allocating.
// Like net.ParseIP it rejects a "%zone" suffix, which ParseIPv6 accepts: screen '%' yourself
// before replacing IsIPv6 with ParseIPv6.
func IsIPv6(s string) bool {
	ellipsis := -1 // position of ellipsis in ip

	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		if len(s) == 0 {
			return true
		}
	}

	i := 0
	for i < net.IPv6len {
		// Bound the digit count, not the value: leading zeros would let "00001" through.
		ci := 0
		for ; ci < len(s); ci++ {
			c := s[ci]
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				break
			}
			if ci == hexFieldMaxLen {
				return false
			}
		}
		if ci == 0 {
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

		i += 2

		s = s[ci:]
		if len(s) == 0 {
			break
		}

		if s[0] != ':' || len(s) == 1 {
			return false
		}
		s = s[1:]

		if s[0] == ':' {
			if ellipsis >= 0 {
				return false
			}
			ellipsis = i
			s = s[1:]
			if len(s) == 0 {
				break
			}
		}
	}

	if len(s) != 0 {
		return false
	}

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
