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

// Bytes returns a []byte view of s without allocation, relying on string and
// slice headers both starting with the data pointer. For string-backed S, the
// result aliases immutable memory and must only be read.
func Bytes[S ~string | ~[]byte](s S) []byte {
	// #nosec G103
	return unsafe.Slice(*(**byte)(unsafe.Pointer(&s)), len(s))
}

// Seq is Bytes' inverse: it returns b viewed as S without allocation. For
// string-backed S, b must never be written again.
func Seq[S ~string | ~[]byte](b []byte) S {
	// #nosec G103
	return *(*S)(unsafe.Pointer(&b))
}
