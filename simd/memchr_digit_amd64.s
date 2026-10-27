// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loop replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, and the 32-byte main
// loop replaced with a 4x unrolled 128-byte block.

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// DIGITS leaves 0xFF in each lane of src holding an ASCII digit (t0 scratch;
// dst may alias src). Adding Y0 = [0x46 x 32] maps '0'..'9' onto 0x76..0x7F, so
// one signed compare with Y1 = [0x75 x 32] is exact: 0x3A..0xB9 become negative,
// 0x00..0x2F stay <= 0x75, and 0xBA..0xFF wrap to 0x00..0x45.
#define DIGITS(src, dst, t0) \
	VPADDB   Y0, src, t0; \
	VPCMPGTB Y1, t0, dst

// func memchrDigitAVX2(haystack []byte) int
//
// Skeleton and bounds as in memchr_amd64.s: nothing outside haystack[0:len] is
// read; under 32 bytes (not reachable via the Go dispatch) a scalar loop runs.
TEXT ·memchrDigitAVX2(SB), NOSPLIT, $0-32
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX

	TESTQ   DX, DX
	JZ      not_found

	MOVQ    $0x4646464646464646, AX
	VMOVQ   AX, X0
	VPBROADCASTQ X0, Y0

	MOVQ    $0x7575757575757575, AX
	VMOVQ   AX, X1
	VPBROADCASTQ X1, Y1

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      tail_loop

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      loop32_entry

	BLOCKALIGN
loop128:
	VMOVDQU (SI), Y2
	VMOVDQU 32(SI), Y3
	VMOVDQU 64(SI), Y4
	VMOVDQU 96(SI), Y5
	DIGITS(Y2, Y2, Y6)
	DIGITS(Y3, Y3, Y6)
	DIGITS(Y4, Y4, Y6)
	DIGITS(Y5, Y5, Y6)

	VPOR    Y2, Y3, Y6
	VPOR    Y4, Y5, Y7
	VPOR    Y6, Y7, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128

loop32_entry:
	CMPQ    SI, R12
	JA      last32

loop32:
	VMOVDQU (SI), Y2
	DIGITS(Y2, Y2, Y6)
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32

last32:
	// Overlapping final window; lanes before SI are non-digits, so BSF is exact.
	CMPQ    SI, R8
	JAE     not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y2
	DIGITS(Y2, Y2, Y6)
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	JMP     not_found

tail_loop:
	MOVBLZX (SI), BX

	CMPB    BL, $0x30
	JB      tail_next
	CMPB    BL, $0x39
	JA      tail_next

	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET

tail_next:
	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop

not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

found_in_block:
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX

found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI
	ADDQ    SI, BX
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET
