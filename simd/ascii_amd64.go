// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// vectorLen is the AVX2 register width in bytes. MinLen happens to equal
// it today, but the two mean different things — MinLen is a dispatch
// tuning knob, vectorLen is a hardware fact — so the alignment arithmetic
// below uses this constant rather than MinLen.
const vectorLen = 32

// The ASCII kernels are implemented in ascii_amd64.s. All three read the
// same signal (a byte is non-ASCII exactly when its high bit is set, which
// VPMOVMSKB gathers 32 at a time) and process four vectors per iteration.
//
//go:noescape
func isASCIIAVX2(data []byte) bool

//go:noescape
func firstNonASCIIAVX2(data []byte) int

// countNonASCIIAVX2 requires len(data) to be a multiple of vectorLen; see
// countNonASCIIImpl.
//
//go:noescape
func countNonASCIIAVX2(data []byte) int

// isASCIIImpl dispatches to the AVX2 kernel for inputs large enough to amortize
// the vector setup; smaller inputs and pre-AVX2 CPUs take the SWAR path.
func isASCIIImpl(data []byte) bool {
	if hasAVX2 && len(data) >= MinLen {
		return isASCIIAVX2(data)
	}
	return isASCIIGeneric(data)
}

// firstNonASCIIImpl dispatches like isASCIIImpl.
func firstNonASCIIImpl(data []byte) int {
	if hasAVX2 && len(data) >= MinLen {
		return firstNonASCIIAVX2(data)
	}
	return firstNonASCIIGeneric(data)
}

// countNonASCIIImpl splits the input at the last whole vector. Counting is
// the one scan here that cannot finish with the overlapping final window
// the first-match kernels use — rescanning bytes would double-count them —
// so instead of paying for a masked epilogue in assembly, the kernel takes
// the vector-aligned prefix and the SWAR loop counts the at most 31
// remaining bytes.
func countNonASCIIImpl(data []byte) int {
	if hasAVX2 && len(data) >= MinLen {
		head := len(data) &^ (vectorLen - 1)
		return countNonASCIIAVX2(data[:head]) + countNonASCIIGeneric(data[head:])
	}
	return countNonASCIIGeneric(data)
}
