package unsafeconv

import "unsafe"

// UnsafeString returns a string view of b without allocation.
func UnsafeString(b []byte) string {
	// #nosec G103
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes returns a byte slice view of s without allocation.
func UnsafeBytes(s string) []byte {
	// #nosec G103
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Bytes returns a []byte view of s without allocation for any string- or
// byte-slice-backed type. It relies on string and slice headers both
// starting with the data pointer (guaranteed by the gc runtime layout, and
// the same assumption UnsafeString makes). When S's underlying type is
// string, the result aliases immutable memory and must only be read.
func Bytes[S ~string | ~[]byte](s S) []byte {
	// #nosec G103
	return unsafe.Slice(*(**byte)(unsafe.Pointer(&s)), len(s))
}
