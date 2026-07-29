// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loop replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, the 32-byte main loop
// replaced with a 4x unrolled 128-byte block, and firstNonASCIIAVX2 /
// countNonASCIIAVX2 added so the locating and counting scans no longer
// have to fall back to SWAR words on amd64.

//go:build amd64

#include "textflag.h"

// All three kernels share one observation: a byte is non-ASCII exactly when
// its high bit is set, and VPMOVMSKB gathers those 32 high bits into a
// scalar register in one instruction. OR-ing four vectors before the
// extraction therefore preserves "some byte in these 128 was non-ASCII"
// exactly, which is what lets the block loop test 128 bytes with a single
// port-0 VPMOVMSKB. Both loops test their bound at the bottom against a
// precomputed limit pointer.

// func isASCIIAVX2(data []byte) bool
//
// AVX2 implementation of ASCII detection that checks if all bytes are < 0x80.
//
// Algorithm:
//  1. Block loop: OR four 32-byte vectors, extract the high bits with
//     VPMOVMSKB, and return false if any is set. The loads fold into the
//     VPOR operands, so a block costs four memory ops and one extraction.
//  2. 32-byte loop for the remainder, then one overlapping 32-byte window
//     ending at the buffer end (rescanning bytes already proven ASCII
//     cannot change the outcome). A scalar loop serves only sub-32 direct
//     calls, unreachable through the Go dispatch.
//  3. Always call VZEROUPPER before returning (critical for performance)
//
// Parameters (FP offsets):
//   data_base+0(FP)  - pointer to data (8 bytes)
//   data_len+8(FP)   - data length (8 bytes)
//   data_cap+16(FP)  - data capacity (8 bytes, unused but part of slice)
//   ret+24(FP)       - return value: bool (1 byte)
//
// Total argument frame size: 25 bytes (8+8+8+1)
TEXT ·isASCIIAVX2(SB), NOSPLIT, $0-25
	// Load parameters
	MOVQ    data_base+0(FP), SI     // SI = data pointer
	MOVQ    data_len+8(FP), DX      // DX = data length

	// Empty data check
	TESTQ   DX, DX
	JZ      all_ascii

	// Calculate end pointer
	LEAQ    (SI)(DX*1), R8          // R8 = SI + length (end pointer)

	CMPQ    DX, $32
	JB      tail_loop               // sub-32 direct call: scalar

	LEAQ    -32(R8), R12            // last valid 32-byte window base
	LEAQ    -128(R8), R11           // last valid 128-byte block base
	CMPQ    SI, R11
	JA      loop32_entry

loop128:
	VMOVDQU (SI), Y0
	VPOR    32(SI), Y0, Y0
	VPOR    64(SI), Y0, Y0
	VPOR    96(SI), Y0, Y0
	VPMOVMSKB Y0, AX
	TESTL   AX, AX
	JNZ     found_non_ascii

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128

loop32_entry:
	CMPQ    SI, R12
	JA      last32

loop32:
	// Load 32 bytes (unaligned load is safe and fast with AVX2) and
	// extract the high bit of each byte into a 32-bit mask.
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, AX
	TESTL   AX, AX
	JNZ     found_non_ascii

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32

last32:
	// Fewer than 32 bytes remain past SI; retest the final window.
	CMPQ    SI, R8
	JAE     all_ascii
	VMOVDQU (R12), Y0
	VPMOVMSKB Y0, AX
	TESTL   AX, AX
	JNZ     found_non_ascii
	JMP     all_ascii

tail_loop:
	// Load one byte and check high bit
	MOVBLZX (SI), AX                // AX = data[SI] (zero-extended)
	TESTB   $0x80, AL               // Check if high bit is set
	JNZ     found_non_ascii         // If set, non-ASCII found

	// Advance to next byte
	INCQ    SI
	CMPQ    SI, R8                  // Check if reached end
	JB      tail_loop               // Continue if SI < end

all_ascii:
	// All bytes are ASCII
	MOVB    $1, ret+24(FP)          // Return true
	VZEROUPPER                      // Clear upper YMM bits (CRITICAL!)
	RET

found_non_ascii:
	MOVB    $0, ret+24(FP)          // Return false
	VZEROUPPER                      // Clear upper YMM bits
	RET

// func firstNonASCIIAVX2(data []byte) int
//
// AVX2 kernel behind FirstNonASCII: the same scan as isASCIIAVX2, but the
// block loop keeps the four source vectors live so the hit path can
// re-extract them in address order and BSF the first non-ASCII byte.
//
// Parameters (FP offsets):
//   data_base+0(FP)  - pointer to data (8 bytes)
//   data_len+8(FP)   - data length (8 bytes)
//   data_cap+16(FP)  - data capacity (8 bytes, unused)
//   ret+24(FP)       - return value: index or -1 (8 bytes)
TEXT ·firstNonASCIIAVX2(SB), NOSPLIT, $0-32
	MOVQ    data_base+0(FP), SI
	MOVQ    data_len+8(FP), DX

	TESTQ   DX, DX
	JZ      fna_not_found

	MOVQ    SI, DI                  // DI = data start (preserved)
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      fna_tail_loop           // sub-32 direct call: scalar

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      fna_loop32_entry

fna_loop128:
	VMOVDQU (SI), Y0
	VMOVDQU 32(SI), Y1
	VMOVDQU 64(SI), Y2
	VMOVDQU 96(SI), Y3
	VPOR    Y0, Y1, Y4
	VPOR    Y2, Y3, Y5
	VPOR    Y4, Y5, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     fna_found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     fna_loop128

fna_loop32_entry:
	CMPQ    SI, R12
	JA      fna_last32

fna_loop32:
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, CX
	TESTL   CX, CX
	JNZ     fna_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     fna_loop32

fna_last32:
	// Fewer than 32 bytes remain past SI; rescan the final window. Lanes
	// before SI are known ASCII, so the first set lane falls in the new
	// bytes.
	CMPQ    SI, R8
	JAE     fna_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, CX
	TESTL   CX, CX
	JNZ     fna_found_in_vector
	JMP     fna_not_found

fna_tail_loop:
	MOVBLZX (SI), AX
	TESTB   $0x80, AL
	JNZ     fna_found_scalar
	INCQ    SI
	CMPQ    SI, R8
	JB      fna_tail_loop

fna_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

fna_found_in_block:
	VPMOVMSKB Y0, CX
	TESTL   CX, CX
	JNZ     fna_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y1, CX
	TESTL   CX, CX
	JNZ     fna_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     fna_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y3, CX

fna_found_in_vector:
	BSFL    CX, CX
	SUBQ    DI, SI
	ADDQ    SI, CX
	MOVQ    CX, ret+24(FP)
	VZEROUPPER
	RET

fna_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET

// func countNonASCIIAVX2(data []byte) int
//
// AVX2 kernel behind CountNonASCII. Counting is the one scan here that
// cannot finish with the overlapping final window the first-match kernels
// use — rescanning bytes would count them twice — so instead the vector
// loops stop at the last whole vector and a branch-free byte loop takes
// the remaining 0-31 bytes. That keeps the kernel self-bounding for any
// length, like every other kernel in this package: it never reads past
// data[len(data)-1], so no caller-side length preconditioning stands
// between it and an out-of-bounds read.
//
// Each vector's high bits become a 32-bit mask and POPCNT turns it into a
// lane count. Four masks are extracted before the four POPCNTs so the
// port-0 extractions and the port-1 population counts overlap; each POPCNT
// writes the register its own VPMOVMSKB just produced, so the destination
// false dependency Intel CPUs carry on POPCNT is already satisfied.
//
// POPCNT is a separate CPUID feature from AVX2, so countNonASCIIImpl gates
// on both; see cpu_amd64.go.
//
// Parameters (FP offsets):
//   data_base+0(FP)  - pointer to data (8 bytes)
//   data_len+8(FP)   - data length (8 bytes)
//   data_cap+16(FP)  - data capacity (8 bytes, unused)
//   ret+24(FP)       - number of bytes >= 0x80 (8 bytes)
TEXT ·countNonASCIIAVX2(SB), NOSPLIT, $0-32
	MOVQ    data_base+0(FP), SI
	MOVQ    data_len+8(FP), DX

	XORQ    R9, R9                  // R9 = running count
	XORQ    R13, R13                // R13 = second accumulator

	TESTQ   DX, DX
	JZ      cna_done

	LEAQ    (SI)(DX*1), R8          // R8 = end pointer
	CMPQ    DX, $32
	JB      cna_tail                // no whole vector: byte loop only

	LEAQ    -32(R8), R12            // last valid 32-byte window base
	LEAQ    -128(R8), R11           // last valid 128-byte block base
	CMPQ    SI, R11
	JA      cna_loop32_entry

cna_loop128:
	VMOVDQU (SI), Y0
	VMOVDQU 32(SI), Y1
	VMOVDQU 64(SI), Y2
	VMOVDQU 96(SI), Y3
	VPMOVMSKB Y0, AX
	VPMOVMSKB Y1, BX
	VPMOVMSKB Y2, CX
	VPMOVMSKB Y3, R10
	POPCNTL AX, AX
	POPCNTL BX, BX
	POPCNTL CX, CX
	POPCNTL R10, R10
	ADDQ    AX, R9
	ADDQ    BX, R13
	ADDQ    CX, R9
	ADDQ    R10, R13

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     cna_loop128

cna_loop32_entry:
	CMPQ    SI, R12
	JA      cna_tail

cna_loop32:
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, AX
	POPCNTL AX, AX
	ADDQ    AX, R9
	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     cna_loop32

cna_tail:
	// 0-31 bytes left. Shifting the byte down to its high bit turns the
	// test into an add, so the loop carries no data-dependent branch.
	CMPQ    SI, R8
	JAE     cna_done
	MOVBLZX (SI), AX
	SHRL    $7, AX
	ADDQ    AX, R9
	INCQ    SI
	JMP     cna_tail

cna_done:
	ADDQ    R13, R9
	MOVQ    R9, ret+24(FP)
	VZEROUPPER
	RET
