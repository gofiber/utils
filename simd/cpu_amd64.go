//go:build amd64

package simd

import (
	"golang.org/x/sys/cpu"
)

// hasAVX2 gates every AVX2 kernel; x/sys/cpu also checks OS YMM support.
var hasAVX2 = cpu.X86.HasAVX2

// hasPOPCNT gates countNonASCIIAVX2. POPCNT has its own CPUID bit that a
// hypervisor can mask separately from AVX2; ungated, it would SIGILL.
var hasPOPCNT = cpu.X86.HasPOPCNT
