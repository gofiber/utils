// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loops replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, the per-range
// clamp-and-compare classifier replaced with a nibble table lookup, the
// 32-byte main loops replaced with 4x unrolled 128-byte blocks, and the
// block test reduced to one compare over the four raw class values.

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// Nibble classifier for \w = [A-Za-z0-9_]: b is a word character iff its class
// value lowTable[b&15] & highTable[b>>4] is nonzero. highTable gives rows
// 0x3_..0x7_ the bits 0x01..0x10, others 0; lowTable[n] is the OR of the
// bits of the rows where nibble n is a word character. VPSHUFB zeroes lanes
// whose index has bit 7 set, so bytes >= 0x80 classify as non-word.
//
// Skeleton and bounds as in memchr_amd64.s: nothing outside haystack[0:len] is
// read. The block loops reduce the four raw class values (per-lane min for
// NotWord, OR for Word) and compare once; the raw values stay live for the hit
// path.

// WORDBITS: class value of each lane of src into dst (t0 scratch; dst may
// alias src).
#define WORDBITS(src, dst, t0) \
	VPSRLW   $4, src, t0;  \
	VPAND    Y2, t0, t0;   \
	VPSHUFB  t0, Y1, t0;   \
	VPSHUFB  src, Y0, dst; \
	VPAND    t0, dst, dst

// NOTWORD: 0xFF in each non-word lane.
#define NOTWORD(src, dst, t0) \
	WORDBITS(src, dst, t0); \
	VPCMPEQB Y3, dst, dst

// SETUPTABLES: Y0 = lowTable, Y1 = highTable, Y2 = 0x0F mask, Y3 = zero.
#define SETUPTABLES \
	MOVQ        $0x1F1F1F1F1F1F1F15, AX; \
	VMOVQ       AX, X0;                  \
	MOVQ        $0x0E0A0A0A0A1E1F1F, AX; \
	VPINSRQ     $1, AX, X0, X0;          \
	VINSERTI128 $1, X0, Y0, Y0;          \
	MOVQ        $0x1008040201000000, AX; \
	VMOVQ       AX, X1;                  \
	VINSERTI128 $1, X1, Y1, Y1;          \
	MOVQ        $0x0F0F0F0F0F0F0F0F, AX; \
	VMOVQ       AX, X2;                  \
	VPBROADCASTQ X2, Y2;                 \
	VPXOR       Y3, Y3, Y3

// func memchrWordAVX2(haystack []byte) int
TEXT ·memchrWordAVX2(SB), NOSPLIT, $0-32
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX

	TESTQ   DX, DX
	JZ      word_not_found

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      word_tail_loop

	SETUPTABLES

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      word_loop32_entry

	BLOCKALIGN_LONG
word_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	WORDBITS(Y4, Y4, Y8)
	WORDBITS(Y5, Y5, Y8)
	WORDBITS(Y6, Y6, Y8)
	WORDBITS(Y7, Y7, Y8)

	VPOR    Y4, Y5, Y8
	VPOR    Y6, Y7, Y9
	VPOR    Y8, Y9, Y8
	VPCMPEQB Y3, Y8, Y8
	VPMOVMSKB Y8, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     word_loop128

word_loop32_entry:
	CMPQ    SI, R12
	JA      word_last32

word_loop32:
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     word_loop32

word_last32:
	// Overlapping final window: lanes before SI were ruled out, so BSF is exact.
	CMPQ    SI, R8
	JAE     word_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	JMP     word_not_found

word_tail_loop:
	MOVBLZX (SI), AX

	CMPB    AL, $65
	JB      word_check_lower
	CMPB    AL, $90
	JBE     word_found_scalar

word_check_lower:
	CMPB    AL, $97
	JB      word_check_digit
	CMPB    AL, $122
	JBE     word_found_scalar

word_check_digit:
	CMPB    AL, $48
	JB      word_check_underscore
	CMPB    AL, $57
	JBE     word_found_scalar

word_check_underscore:
	CMPB    AL, $95
	JE      word_found_scalar

	INCQ    SI
	CMPQ    SI, R8
	JB      word_tail_loop

word_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

word_found_in_block:
	VPCMPEQB Y3, Y4, Y4
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y5, Y5
	VPMOVMSKB Y5, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y6, Y6
	VPMOVMSKB Y6, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y7, Y7
	VPMOVMSKB Y7, CX
	NOTL    CX

word_found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI
	ADDQ    SI, BX
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET

word_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET

// func memchrNotWordAVX2(haystack []byte) int
TEXT ·memchrNotWordAVX2(SB), NOSPLIT, $0-32
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX

	TESTQ   DX, DX
	JZ      notword_not_found

	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      notword_tail_loop

	SETUPTABLES

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      notword_loop32_entry

	BLOCKALIGN_LONG
notword_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	WORDBITS(Y4, Y4, Y8)
	WORDBITS(Y5, Y5, Y8)
	WORDBITS(Y6, Y6, Y8)
	WORDBITS(Y7, Y7, Y8)

	VPMINUB Y4, Y5, Y8
	VPMINUB Y6, Y7, Y9
	VPMINUB Y8, Y9, Y8
	VPCMPEQB Y3, Y8, Y8
	VPMOVMSKB Y8, CX
	TESTL   CX, CX
	JNZ     notword_found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     notword_loop128

notword_loop32_entry:
	CMPQ    SI, R12
	JA      notword_last32

notword_loop32:
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     notword_loop32

notword_last32:
	CMPQ    SI, R8
	JAE     notword_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	JMP     notword_not_found

notword_tail_loop:
	MOVBLZX (SI), AX

	CMPB    AL, $65
	JB      notword_check_lower_gap
	CMPB    AL, $90
	JBE     notword_is_word

notword_check_lower_gap:
	CMPB    AL, $97
	JB      notword_check_digit_gap
	CMPB    AL, $122
	JBE     notword_is_word

notword_check_digit_gap:
	CMPB    AL, $48
	JB      notword_check_underscore_gap
	CMPB    AL, $57
	JBE     notword_is_word

notword_check_underscore_gap:
	CMPB    AL, $95
	JE      notword_is_word

	JMP     notword_found_scalar

notword_is_word:
	INCQ    SI
	CMPQ    SI, R8
	JB      notword_tail_loop

notword_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

notword_found_in_block:
	VPCMPEQB Y3, Y4, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y5, Y5
	VPMOVMSKB Y5, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y6, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y7, Y7
	VPMOVMSKB Y7, CX

notword_found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI
	ADDQ    SI, BX
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET

notword_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET
