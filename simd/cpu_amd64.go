//go:build amd64

package simd

// hasAVX2 gates every AVX2 kernel in this package. It is set once at init
// from CPUID so the per-call dispatch is a single predictable branch.
var hasAVX2 = detectAVX2()

// cpuid executes the CPUID instruction with the given leaf (EAX) and
// sub-leaf (ECX), yielding the four result registers. Implemented in
// cpu_amd64.s.
func cpuid(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32)

// cpuidRegs packs cpuid's four registers as {EAX, EBX, ECX, EDX} so callers
// can pick single registers without blank-heavy assignments.
func cpuidRegs(leaf, subleaf uint32) [4]uint32 {
	eax, ebx, ecx, edx := cpuid(leaf, subleaf)
	return [4]uint32{eax, ebx, ecx, edx}
}

// xgetbv executes XGETBV with ECX=0, returning the low 32 bits of XCR0.
// It must only be called when CPUID reports OSXSAVE support. Implemented in
// cpu_amd64.s.
func xgetbv() uint32

// detectAVX2 reports whether both the CPU and the OS support AVX2: the CPU
// must advertise AVX and AVX2, and the OS must save the XMM and YMM register
// state on context switches (XCR0 bits 1 and 2). This mirrors the checks in
// golang.org/x/sys/cpu without taking on the dependency.
func detectAVX2() bool {
	maxLeaf := cpuidRegs(0, 0)[0]
	if maxLeaf < 7 {
		return false
	}
	ecx1 := cpuidRegs(1, 0)[2]
	const osxsaveBit = 1 << 27
	const avxBit = 1 << 28
	if ecx1&osxsaveBit == 0 || ecx1&avxBit == 0 {
		return false
	}
	const xmmAndYMMState = 0x6 // XCR0 bit 1 (XMM) and bit 2 (YMM)
	if xgetbv()&xmmAndYMMState != xmmAndYMMState {
		return false
	}
	ebx7 := cpuidRegs(7, 0)[1]
	const avx2Bit = 1 << 5
	return ebx7&avx2Bit != 0
}
