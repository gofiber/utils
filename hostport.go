package utils

import (
	"bytes"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// SplitHostPort splits hostport with net.SplitHostPort's rules, reporting
// failure as ok == false; the parts alias hostport.
func SplitHostPort[S byteSeq](hostport S) (host, port S, ok bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	b := unsafeconv.Bytes(hostport)
	i := bytes.LastIndexByte(b, ':')
	if i < 0 {
		return host, port, false
	}
	j, k, hostEnd := 0, 0, i
	if hostport[0] == '[' {
		end := bytes.IndexByte(b, ']')
		if end < 0 || end+1 != i {
			return host, port, false
		}
		j, k, hostEnd = 1, end+1, end // no '[' resp. ']' before these positions
	} else if bytes.IndexByte(b[:i], ':') >= 0 {
		return host, port, false
	}
	if bytes.IndexByte(b[j:], '[') >= 0 || bytes.IndexByte(b[k:], ']') >= 0 {
		return host, port, false
	}
	return hostport[j:hostEnd], hostport[i+1:], true
}
