// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build !amd64

package simd

func isASCIIImpl(data []byte) bool {
	return isASCIIGeneric(data)
}

func firstNonASCIIImpl(data []byte) int {
	return firstNonASCIIGeneric(data)
}

func countNonASCIIImpl(data []byte) int {
	return countNonASCIIGeneric(data)
}
