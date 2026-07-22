//go:build amd64

package simd

import (
	"golang.org/x/sys/cpu"
)

// hasAVX2 gates every AVX2 kernel in this package. golang.org/x/sys/cpu
// verifies both the CPU capability and OS support for saving the YMM
// register state, and the value is read once at init so the per-call
// dispatch is a single predictable branch.
var hasAVX2 = cpu.X86.HasAVX2
