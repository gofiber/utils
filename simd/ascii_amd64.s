// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loop replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, the 32-byte main loop
// replaced with a 4x unrolled 128-byte block, and firstNonASCIIAVX2 /
// countNonASCIIAVX2 added so the locating and counting scans no longer
// have to fall back to SWAR words on amd64 (the counting kernel's block
// loop accumulating in vector byte lanes rather than extracting and
// population-counting every vector's mask).

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// A byte is non-ASCII iff its high bit is set, so OR-ing four vectors before one
// VPMOVMSKB tests 128 bytes. The kernels are self-bounding and never read
// outside data[0:len] (loads are unaligned): the block loop runs while SI <= R11
// (end-128), the 32-byte loop while SI <= R12 (end-32); sub-32 direct calls use a
// scalar loop. All exits run VZEROUPPER.

// func isASCIIAVX2(data []byte) bool
TEXT ·isASCIIAVX2(SB), NOSPLIT, $0-25
	MOVQ    data_base+0(FP), SI
	MOVQ    data_len+8(FP), DX

	TESTQ   DX, DX
	JZ      all_ascii

	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      tail_loop

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      loop32_entry

	BLOCKALIGN
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
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, AX
	TESTL   AX, AX
	JNZ     found_non_ascii

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32

last32:
	// Overlapping final window; rescanning proven-ASCII bytes is harmless.
	CMPQ    SI, R8
	JAE     all_ascii
	VMOVDQU (R12), Y0
	VPMOVMSKB Y0, AX
	TESTL   AX, AX
	JNZ     found_non_ascii
	JMP     all_ascii

tail_loop:
	MOVBLZX (SI), AX
	TESTB   $0x80, AL
	JNZ     found_non_ascii

	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop

all_ascii:
	MOVB    $1, ret+24(FP)
	VZEROUPPER
	RET

found_non_ascii:
	MOVB    $0, ret+24(FP)
	VZEROUPPER
	RET

// func firstNonASCIIAVX2(data []byte) int
//
// The hit path re-extracts the block's four live vectors in address order.
TEXT ·firstNonASCIIAVX2(SB), NOSPLIT, $0-32
	MOVQ    data_base+0(FP), SI
	MOVQ    data_len+8(FP), DX

	TESTQ   DX, DX
	JZ      fna_not_found

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      fna_tail_loop

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      fna_loop32_entry

	BLOCKALIGN
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
	// Overlapping final window; lanes before SI are ASCII, so BSF stays exact.
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
// No overlapping final window (rescanned bytes would count twice): vector
// loops stop at the last whole vector and a byte loop takes the last 0-31 bytes.
// The block loop subtracts 0 > byte (signed VPCMPGTB: 0xFF exactly for high-bit
// bytes) from byte accumulators Y4/Y5; a lane holds 255 and gains at most 2 per
// block, so FLUSH folds them into the 64-bit totals in Y8 every 127 blocks.
// Leftover vectors use VPMOVMSKB+POPCNT; callers must also check for POPCNT.

// FLUSH needs Y7 = 0 and clobbers Y0/Y1.
#define FLUSH \
	VPSADBW Y7, Y4, Y0; \
	VPSADBW Y7, Y5, Y1; \
	VPADDQ  Y0, Y8, Y8; \
	VPADDQ  Y1, Y8, Y8; \
	VPXOR   Y4, Y4, Y4; \
	VPXOR   Y5, Y5, Y5

TEXT ·countNonASCIIAVX2(SB), NOSPLIT, $0-32
	MOVQ    data_base+0(FP), SI
	MOVQ    data_len+8(FP), DX

	XORQ    R9, R9
	XORQ    R10, R10

	TESTQ   DX, DX
	JZ      cna_done

	LEAQ    (SI)(DX*1), R8
	CMPQ    DX, $32
	JB      cna_tail

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      cna_loop32_entry

	VPXOR   Y7, Y7, Y7
	VPXOR   Y4, Y4, Y4
	VPXOR   Y5, Y5, Y5
	VPXOR   Y8, Y8, Y8
	MOVQ    $127, R13

	// The DECQ/JZ countdown is a third fused pair: hence BLOCKALIGN_LONG.
	BLOCKALIGN_LONG
cna_loop128:
	VPCMPGTB (SI), Y7, Y0
	VPCMPGTB 32(SI), Y7, Y1
	VPCMPGTB 64(SI), Y7, Y2
	VPCMPGTB 96(SI), Y7, Y3
	VPSUBB  Y0, Y4, Y4
	VPSUBB  Y1, Y5, Y5
	VPSUBB  Y2, Y4, Y4
	VPSUBB  Y3, Y5, Y5

	ADDQ    $128, SI
	DECQ    R13
	JZ      cna_flush

cna_loop128_next:
	CMPQ    SI, R11
	JBE     cna_loop128
	JMP     cna_block_done

cna_flush:
	FLUSH
	MOVQ    $127, R13
	JMP     cna_loop128_next

cna_block_done:
	VPSADBW Y7, Y4, Y0
	VPSADBW Y7, Y5, Y1
	VPADDQ  Y0, Y8, Y8
	VPADDQ  Y1, Y8, Y8
	VEXTRACTI128 $1, Y8, X0
	VPADDQ  X0, X8, X8
	VPSRLDQ $8, X8, X0
	VPADDQ  X0, X8, X8
	VMOVQ   X8, R9

cna_loop32_entry:
	CMPQ    SI, R12
	JA      cna_tail

cna_loop32:
	VMOVDQU (SI), Y0
	VPMOVMSKB Y0, AX
	// POPCNT reuses VPMOVMSKB's destination: no false dependency on Intel.
	POPCNTL AX, AX
	ADDQ    AX, R10
	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     cna_loop32

cna_tail:
	CMPQ    SI, R8
	JAE     cna_done
	MOVBLZX (SI), AX
	SHRL    $7, AX
	ADDQ    AX, R10
	INCQ    SI
	JMP     cna_tail

cna_done:
	ADDQ    R10, R9
	MOVQ    R9, ret+24(FP)
	VZEROUPPER
	RET
