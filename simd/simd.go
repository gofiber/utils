// Package simd provides SIMD-accelerated byte searching and validation
// primitives for the Fiber ecosystem: multi-needle scans, character-class
// scans, paired-byte scans, substring search, and ASCII validation.
//
// Acceleration comes in two tiers. The scan kernels (Memchr2, Memchr3,
// MemchrPair, MemchrDigit, MemchrWord, MemchrNotWord, IsASCII, FirstNonASCII,
// CountNonASCII) run AVX2 assembly on amd64 for inputs of MinLen+ bytes
// (MinLen+offset for MemchrPair) and portable SWAR loops built on the swar
// package everywhere else. Memmem prefilters on amd64 (with MemchrPair or a
// rare-byte bytes.IndexByte scan) and delegates to bytes.Index for short
// haystacks or without AVX2. MemchrInTable and MemchrNotInTable are scalar on
// every architecture, and single-byte search is absent because bytes.IndexByte
// is already vectorized.
//
// Unlike the top-level utils helpers, this package operates on []byte only.
//
// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
package simd

// MinLen is the input length from which the AVX2 kernels engage on amd64;
// shorter inputs take the SWAR fallbacks. Callers that route between their own
// word loops and this package should gate on it rather than hard-code 32.
const MinLen = 32

// Accelerated reports whether the AVX2 assembly kernels are in use on this CPU.
func Accelerated() bool {
	return hasAVX2
}
