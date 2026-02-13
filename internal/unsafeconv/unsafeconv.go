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
