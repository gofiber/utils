// Package simd provides SIMD-accelerated byte searching and validation
// primitives for the Fiber ecosystem: multi-needle scans, character-class
// scans, paired-byte scans, substring search, and ASCII validation.
//
// On amd64 CPUs with AVX2, the functions dispatch to assembly kernels that
// process 32 bytes per iteration; everywhere else they fall back to portable
// SWAR (SIMD within a register) implementations, so the package builds and
// behaves identically on every platform. Accelerated reports which mode is
// active. Single-byte search is deliberately absent: bytes.IndexByte is
// already vector-accelerated by the Go runtime on all major architectures.
//
// Unlike the top-level utils helpers, this package operates on []byte only —
// it is the raw engine underneath the generic helpers. Callers holding a
// string can use the top-level generic functions instead of converting.
//
// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
package simd

// SWAR lane constants shared by the portable fallbacks: lo8 carries 0x01 and
// hi8 carries 0x80 in every byte lane of a 64-bit word.
const (
	lo8 = uint64(0x0101010101010101)
	hi8 = uint64(0x8080808080808080)
)

// Accelerated reports whether the AVX2 assembly kernels are in use on this
// CPU. When false, all functions use portable fallback implementations.
func Accelerated() bool {
	return hasAVX2
}
