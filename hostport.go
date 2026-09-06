package utils

import (
	"bytes"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// SplitHostPort splits "host:port", "[host]:port", and their %zone forms
// with net.SplitHostPort's exact rules, reporting failure as ok == false
// instead of a *net.AddrError; it accepts byte slices as well as strings,
// and the parts alias hostport.
func SplitHostPort[S byteSeq](hostport S) (host, port S, ok bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	b := unsafeconv.Bytes(hostport)
	// The port starts after the last colon.
	i := bytes.LastIndexByte(b, ':')
	if i < 0 {
		return host, port, false
	}
	j, k := 0, 0
	if hostport[0] == '[' {
		// Expect the first ']' just before the last ':'.
		end := bytes.IndexByte(b, ']')
		if end < 0 || end+1 != i {
			return host, port, false
		}
		host = hostport[1:end]
		j, k = 1, end+1 // there can't be a '[' resp. ']' before these positions
	} else {
		host = hostport[:i]
		if bytes.IndexByte(b[:i], ':') >= 0 {
			var zero S
			return zero, port, false
		}
	}
	if bytes.IndexByte(b[j:], '[') >= 0 || bytes.IndexByte(b[k:], ']') >= 0 {
		var zero S
		return zero, port, false
	}
	return host, hostport[i+1:], true
}
