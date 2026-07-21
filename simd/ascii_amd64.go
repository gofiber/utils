// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// isASCIIAVX2 is implemented in ascii_amd64.s.
//
//go:noescape
func isASCIIAVX2(data []byte) bool

// isASCII dispatches to the AVX2 kernel for inputs large enough to amortize
// the vector setup; smaller inputs and pre-AVX2 CPUs take the SWAR path.
func isASCII(data []byte) bool {
	if hasAVX2 && len(data) >= MinLen {
		return isASCIIAVX2(data)
	}
	return isASCIIGeneric(data)
}
