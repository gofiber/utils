// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// The ascii_amd64.s kernels are self-bounding: they clamp their own vector
// loops and finish short remainders, so they carry no length precondition;
// the Impl functions gate on MinLen to amortize vector setup.
//
//go:noescape
func isASCIIAVX2(data []byte) bool

//go:noescape
func firstNonASCIIAVX2(data []byte) int

//go:noescape
func countNonASCIIAVX2(data []byte) int

func isASCIIImpl(data []byte) bool {
	if hasAVX2 && len(data) >= MinLen {
		return isASCIIAVX2(data)
	}
	return isASCIIGeneric(data)
}

func firstNonASCIIImpl(data []byte) int {
	if hasAVX2 && len(data) >= MinLen {
		return firstNonASCIIAVX2(data)
	}
	return firstNonASCIIGeneric(data)
}

// countNonASCIIImpl also needs POPCNT (see hasPOPCNT). No prefix/suffix split:
// the kernel counts its own remainder, and a split costs a second call.
func countNonASCIIImpl(data []byte) int {
	if hasAVX2 && hasPOPCNT && len(data) >= MinLen {
		return countNonASCIIAVX2(data)
	}
	return countNonASCIIGeneric(data)
}
