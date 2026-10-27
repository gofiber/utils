// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loops replaced with overlapping vector rescans of the final
// window (a vector pair in memchrPairAVX2) for the MinLen+ inputs the Go
// dispatch guarantees, and the 32-byte main loops replaced with 4x
// unrolled 128-byte blocks that OR the per-vector masks together and
// extract a single VPMOVMSKB per block.

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// The kernels share one skeleton:
//
//  1. A 128-byte block loop while base <= end-128 (R11). Four vectors are
//     compared, their match masks OR-ed and tested with one VPMOVMSKB; the
//     per-vector masks stay live so the (rare) hit path re-extracts them in
//     address order to find the first match.
//  2. A 32-byte loop while base <= end-32 (R12).
//  3. One overlapping 32-byte window ending at the buffer end for the last
//     1-31 bytes. Lanes before the current base were already ruled out, so
//     the first set lane lies in the new bytes and BSF stays exact.
//
// Loads are unaligned and never reach outside haystack[0:len]. Inputs under 32
// bytes, unreachable through the Go dispatch, take the scalar tail loops. All
// exits run VZEROUPPER.

// func memchr2AVX2(haystack []byte, needle1, needle2 byte) int
//
// Either of two bytes: two broadcast needles, match masks OR-ed.
TEXT ·memchr2AVX2(SB), NOSPLIT, $0-40
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX
	MOVBQZX needle1+24(FP), AX
	MOVBQZX needle2+25(FP), BX

	TESTQ   DX, DX
	JZ      not_found2

	VMOVD   AX, X0
	VPBROADCASTB X0, Y0
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      tail_loop2

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      loop32_2_entry

	BLOCKALIGN
loop128_2:
	VMOVDQU (SI), Y2
	VMOVDQU 32(SI), Y3
	VMOVDQU 64(SI), Y4
	VMOVDQU 96(SI), Y5

	VPCMPEQB Y0, Y2, Y6
	VPCMPEQB Y1, Y2, Y7
	VPOR    Y6, Y7, Y2
	VPCMPEQB Y0, Y3, Y6
	VPCMPEQB Y1, Y3, Y7
	VPOR    Y6, Y7, Y3
	VPCMPEQB Y0, Y4, Y6
	VPCMPEQB Y1, Y4, Y7
	VPOR    Y6, Y7, Y4
	VPCMPEQB Y0, Y5, Y6
	VPCMPEQB Y1, Y5, Y7
	VPOR    Y6, Y7, Y5

	VPOR    Y2, Y3, Y6
	VPOR    Y4, Y5, Y7
	VPOR    Y6, Y7, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     found_in_block2

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128_2

loop32_2_entry:
	CMPQ    SI, R12
	JA      last32_2

loop32_2:
	VMOVDQU (SI), Y2
	VPCMPEQB Y0, Y2, Y3
	VPCMPEQB Y1, Y2, Y4
	VPOR    Y3, Y4, Y3
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector2

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32_2

last32_2:
	CMPQ    SI, R8
	JAE     not_found2
	MOVQ    R12, SI
	VMOVDQU (SI), Y2
	VPCMPEQB Y0, Y2, Y3
	VPCMPEQB Y1, Y2, Y4
	VPOR    Y3, Y4, Y3
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector2
	JMP     not_found2

tail_loop2:
	MOVBLZX (SI), CX
	CMPB    CL, AL
	JE      found_scalar2
	CMPB    CL, BL
	JE      found_scalar2

	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop2

not_found2:
	MOVQ    $-1, AX
	MOVQ    AX, ret+32(FP)
	VZEROUPPER
	RET

found_in_block2:
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector2
	ADDQ    $32, SI
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector2
	ADDQ    $32, SI
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector2
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX

found_in_vector2:
	BSFL    CX, CX
	SUBQ    DI, SI
	ADDQ    SI, CX
	MOVQ    CX, ret+32(FP)
	VZEROUPPER
	RET

found_scalar2:
	SUBQ    DI, SI
	MOVQ    SI, ret+32(FP)
	VZEROUPPER
	RET

// func memchr3AVX2(haystack []byte, needle1, needle2, needle3 byte) int
//
// Same as memchr2AVX2 with a third needle.
TEXT ·memchr3AVX2(SB), NOSPLIT, $0-40
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX
	MOVBQZX needle1+24(FP), AX
	MOVBQZX needle2+25(FP), BX
	MOVBQZX needle3+26(FP), R10

	TESTQ   DX, DX
	JZ      not_found3

	VMOVD   AX, X0
	VPBROADCASTB X0, Y0
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1
	VMOVD   R10, X2
	VPBROADCASTB X2, Y2

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      tail_loop3

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      loop32_3_entry

	BLOCKALIGN
loop128_3:
	VMOVDQU (SI), Y3
	VMOVDQU 32(SI), Y4
	VMOVDQU 64(SI), Y5
	VMOVDQU 96(SI), Y6

	VPCMPEQB Y0, Y3, Y7
	VPCMPEQB Y1, Y3, Y8
	VPCMPEQB Y2, Y3, Y9
	VPOR    Y7, Y8, Y7
	VPOR    Y7, Y9, Y3
	VPCMPEQB Y0, Y4, Y7
	VPCMPEQB Y1, Y4, Y8
	VPCMPEQB Y2, Y4, Y9
	VPOR    Y7, Y8, Y7
	VPOR    Y7, Y9, Y4
	VPCMPEQB Y0, Y5, Y7
	VPCMPEQB Y1, Y5, Y8
	VPCMPEQB Y2, Y5, Y9
	VPOR    Y7, Y8, Y7
	VPOR    Y7, Y9, Y5
	VPCMPEQB Y0, Y6, Y7
	VPCMPEQB Y1, Y6, Y8
	VPCMPEQB Y2, Y6, Y9
	VPOR    Y7, Y8, Y7
	VPOR    Y7, Y9, Y6

	VPOR    Y3, Y4, Y7
	VPOR    Y5, Y6, Y8
	VPOR    Y7, Y8, Y7
	VPMOVMSKB Y7, CX
	TESTL   CX, CX
	JNZ     found_in_block3

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128_3

loop32_3_entry:
	CMPQ    SI, R12
	JA      last32_3

loop32_3:
	VMOVDQU (SI), Y3
	VPCMPEQB Y0, Y3, Y4
	VPCMPEQB Y1, Y3, Y5
	VPCMPEQB Y2, Y3, Y6
	VPOR    Y4, Y5, Y4
	VPOR    Y4, Y6, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector3

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32_3

last32_3:
	CMPQ    SI, R8
	JAE     not_found3
	MOVQ    R12, SI
	VMOVDQU (SI), Y3
	VPCMPEQB Y0, Y3, Y4
	VPCMPEQB Y1, Y3, Y5
	VPCMPEQB Y2, Y3, Y6
	VPOR    Y4, Y5, Y4
	VPOR    Y4, Y6, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector3
	JMP     not_found3

tail_loop3:
	MOVBLZX (SI), CX
	CMPB    CL, AL
	JE      found_scalar3
	CMPB    CL, BL
	JE      found_scalar3
	CMPB    CL, R10B
	JE      found_scalar3

	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop3

not_found3:
	MOVQ    $-1, AX
	MOVQ    AX, ret+32(FP)
	VZEROUPPER
	RET

found_in_block3:
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector3
	ADDQ    $32, SI
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector3
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX
	TESTL   CX, CX
	JNZ     found_in_vector3
	ADDQ    $32, SI
	VPMOVMSKB Y6, CX

found_in_vector3:
	BSFL    CX, CX
	SUBQ    DI, SI
	ADDQ    SI, CX
	MOVQ    CX, ret+32(FP)
	VZEROUPPER
	RET

found_scalar3:
	SUBQ    DI, SI
	MOVQ    SI, ret+32(FP)
	VZEROUPPER
	RET

// func memchrPairAVX2(haystack []byte, byte1, byte2 byte, offset int) int
//
// Finds the first p with haystack[p] == byte1 and haystack[p+offset] == byte2:
// each step ANDs the compare of the 32 bytes at p with byte1 and of the 32 at
// p+offset with byte2; the block loop does four such pairs at once. Both loads
// must stay in bounds, so the last valid base is R12 = start+len-offset-32
// (R11 = R12-96 for blocks), and offset >= len (unsigned) finds nothing.
TEXT ·memchrPairAVX2(SB), NOSPLIT, $0-48
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX
	MOVBQZX byte1+24(FP), AX
	MOVBQZX byte2+25(FP), BX
	MOVQ    offset+32(FP), R10

	TESTQ   DX, DX
	JZ      not_found_pair
	CMPQ    R10, DX
	JAE     not_found_pair

	VMOVD   AX, X0
	VPBROADCASTB X0, Y0
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1

	MOVQ    SI, DI

	MOVQ    DX, R12
	SUBQ    R10, R12
	SUBQ    $32, R12
	JL      scalar_entry_pair
	ADDQ    SI, R12

	LEAQ    -96(R12), R11
	CMPQ    SI, R11
	JA      loop32_pair_entry

	BLOCKALIGN
loop128_pair:
	VMOVDQU (SI), Y2
	VMOVDQU 32(SI), Y3
	VMOVDQU 64(SI), Y4
	VMOVDQU 96(SI), Y5
	VPCMPEQB Y0, Y2, Y2
	VPCMPEQB Y0, Y3, Y3
	VPCMPEQB Y0, Y4, Y4
	VPCMPEQB Y0, Y5, Y5

	VMOVDQU (SI)(R10*1), Y6
	VMOVDQU 32(SI)(R10*1), Y7
	VMOVDQU 64(SI)(R10*1), Y8
	VMOVDQU 96(SI)(R10*1), Y9
	VPCMPEQB Y1, Y6, Y6
	VPCMPEQB Y1, Y7, Y7
	VPCMPEQB Y1, Y8, Y8
	VPCMPEQB Y1, Y9, Y9

	VPAND   Y2, Y6, Y2
	VPAND   Y3, Y7, Y3
	VPAND   Y4, Y8, Y4
	VPAND   Y5, Y9, Y5

	VPOR    Y2, Y3, Y6
	VPOR    Y4, Y5, Y7
	VPOR    Y6, Y7, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     found_in_block_pair

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128_pair

loop32_pair_entry:
	CMPQ    SI, R12
	JA      vector_tail_pair

loop32_pair:
	VMOVDQU (SI), Y2
	VMOVDQU (SI)(R10*1), Y3
	VPCMPEQB Y0, Y2, Y4
	VPCMPEQB Y1, Y3, Y5

	VPAND   Y4, Y5, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32_pair

vector_tail_pair:
	// Fewer than 32 pair positions remain; if any, rescan the window at R12,
	// whose second load ends exactly at the buffer end. Earlier positions were
	// ruled out.
	LEAQ    (DI)(DX*1), R9
	SUBQ    R10, R9
	CMPQ    SI, R9
	JAE     not_found_pair
	MOVQ    R12, SI
	VMOVDQU (SI), Y2
	VMOVDQU (SI)(R10*1), Y3
	VPCMPEQB Y0, Y2, Y4
	VPCMPEQB Y1, Y3, Y5
	VPAND   Y4, Y5, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair
	JMP     not_found_pair

scalar_entry_pair:
	// Scalar path: direct calls with fewer than offset+32 bytes only.
	LEAQ    (DI)(DX*1), R8

tail_loop_pair:
	LEAQ    (SI)(R10*1), R9
	CMPQ    R9, R8
	JAE     not_found_pair

	MOVBLZX (SI), CX
	CMPB    CL, AL
	JNE     tail_next_pair

	MOVBLZX (SI)(R10*1), CX
	CMPB    CL, BL
	JE      found_scalar_pair

tail_next_pair:
	INCQ    SI
	JMP     tail_loop_pair

not_found_pair:
	MOVQ    $-1, AX
	MOVQ    AX, ret+40(FP)
	VZEROUPPER
	RET

found_in_block_pair:
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair
	ADDQ    $32, SI
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair
	ADDQ    $32, SI
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX

found_in_vector_pair:
	BSFL    CX, CX
	SUBQ    DI, SI
	ADDQ    SI, CX
	MOVQ    CX, ret+40(FP)
	VZEROUPPER
	RET

found_scalar_pair:
	SUBQ    DI, SI
	MOVQ    SI, ret+40(FP)
	VZEROUPPER
	RET
