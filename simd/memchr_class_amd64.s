// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loops replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, the per-range
// clamp-and-compare classifier replaced with a nibble table lookup, and
// the 32-byte main loops replaced with 4x unrolled 128-byte blocks.

//go:build amd64

#include "textflag.h"

// Classifying \w = [A-Za-z0-9_] with unsigned clamp-and-compare costs three
// VPMINUB/VPMAXUB/VPCMPEQB triples plus a VPCMPEQB for '_' and three VPORs:
// thirteen vector ops per 32 bytes, ten of them on the two general vector
// ALU ports. A nibble table does the same job in six, and the two VPSHUFBs
// it adds run on the otherwise idle shuffle port.
//
// The table splits each byte into its high and low nibble. Every high
// nibble that contains word characters gets one bit:
//
//	0x3_ -> 0x01   0x4_ -> 0x02   0x5_ -> 0x04   0x6_ -> 0x08   0x7_ -> 0x10
//
// and each low-nibble entry carries the set of high-nibble rows in which
// that low nibble is a word character, so a byte is a word character
// exactly when lowTable[lo] & highTable[hi] != 0:
//
//	lo 0x0        -> 0x15   ('0', 'P', 'p')
//	lo 0x1..0x9   -> 0x1F   (digits and letters in every row)
//	lo 0xA        -> 0x1E   ('J', 'Z', 'j', 'z' but not ':')
//	lo 0xB..0xE   -> 0x0A   ('K'..'N', 'k'..'n' only)
//	lo 0xF        -> 0x0E   ('O', '_', 'o')
//
// Bytes >= 0x80 need no separate guard: VPSHUFB zeroes any lane whose
// index byte has bit 7 set, so the low-nibble lookup already returns 0 for
// them and they classify as non-word, matching the SWAR fallback.
//
// Both kernels compute the same "not a word character" mask (0xFF per
// non-word lane): the low and high lookups are AND-ed and compared against
// zero, which is the complement the class test needs anyway.
// memchrNotWord ORs those masks together across a block and scans the
// extracted bits directly; memchrWord ANDs them (a block holds a word
// character exactly when some lane's non-word mask is clear) and inverts
// the extracted bits with a single NOTL.

// NOTWORD leaves in dst a mask with 0xFF in every lane of src that is NOT
// a word character. t0 and t1 are scratch; dst may alias src.
#define NOTWORD(src, dst, t0, t1) \
	VPSRLW   $4, src, t0;  \
	VPAND    Y2, t0, t0;   \
	VPSHUFB  t0, Y1, t0;   \
	VPSHUFB  src, Y0, t1;  \
	VPAND    t0, t1, t0;   \
	VPCMPEQB Y3, t0, dst

// SETUPTABLES loads the nibble tables, the low-nibble mask and a zero
// vector into Y0-Y3. The tables are built from immediates rather than a
// RODATA section so the kernels stay position-independent and free of
// relocations.
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
//
// AVX2 implementation to find the first word character [A-Za-z0-9_].
//
// Parameters:
//   haystack_base+0(FP)  - pointer to haystack (8 bytes)
//   haystack_len+8(FP)   - length (8 bytes)
//   haystack_cap+16(FP)  - capacity (8 bytes, unused)
//   ret+24(FP)           - return value (8 bytes)
//
TEXT ·memchrWordAVX2(SB), NOSPLIT, $0-32
	// Load parameters
	MOVQ    haystack_base+0(FP), SI    // SI = haystack pointer
	MOVQ    haystack_len+8(FP), DX     // DX = length

	// Empty check
	TESTQ   DX, DX
	JZ      word_not_found

	// Save start for offset calculation and compute the end pointer
	MOVQ    SI, DI                     // DI = start pointer
	LEAQ    (SI)(DX*1), R8             // R8 = end pointer

	CMPQ    DX, $32
	JB      word_tail_loop             // sub-32 direct call: scalar

	SETUPTABLES

	LEAQ    -32(R8), R12               // last valid 32-byte window base
	LEAQ    -128(R8), R11              // last valid 128-byte block base
	CMPQ    SI, R11
	JA      word_loop32_entry

word_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	NOTWORD(Y4, Y4, Y8, Y9)
	NOTWORD(Y5, Y5, Y8, Y9)
	NOTWORD(Y6, Y6, Y8, Y9)
	NOTWORD(Y7, Y7, Y8, Y9)

	// A word character exists in the block iff some lane's non-word mask
	// is clear, so AND the four masks and look for a zero lane.
	VPAND   Y4, Y5, Y8
	VPAND   Y6, Y7, Y9
	VPAND   Y8, Y9, Y8
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
	NOTWORD(Y4, Y4, Y8, Y9)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     word_loop32

word_last32:
	// Fewer than 32 bytes remain past SI; rescan the final window. Lanes
	// before SI were already resolved as non-matching, so BSF stays exact.
	CMPQ    SI, R8
	JAE     word_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8, Y9)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	JMP     word_not_found

word_tail_loop:
	MOVBLZX (SI), AX                   // AX = byte

	// Check [A-Z]: 65 <= b <= 90
	CMPB    AL, $65
	JB      word_check_lower
	CMPB    AL, $90
	JBE     word_found_scalar

word_check_lower:
	// Check [a-z]: 97 <= b <= 122
	CMPB    AL, $97
	JB      word_check_digit
	CMPB    AL, $122
	JBE     word_found_scalar

word_check_digit:
	// Check [0-9]: 48 <= b <= 57
	CMPB    AL, $48
	JB      word_check_underscore
	CMPB    AL, $57
	JBE     word_found_scalar

word_check_underscore:
	// Check '_': 95
	CMPB    AL, $95
	JE      word_found_scalar

	// Not a word char, continue
	INCQ    SI
	CMPQ    SI, R8
	JB      word_tail_loop

word_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

word_found_in_block:
	// Locate the first hit by re-extracting the per-vector masks in
	// address order; only one of the four can hold the lowest match.
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y6, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y7, CX
	NOTL    CX

word_found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI                     // SI = offset to chunk start
	ADDQ    SI, BX                     // BX = absolute position
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET

word_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET

// func memchrNotWordAVX2(haystack []byte) int
//
// AVX2 implementation to find the first non-word character, the complement
// of memchrWordAVX2. The nibble classifier already produces the non-word
// mask directly, so this kernel skips the NOTL memchrWordAVX2 needs.
//
TEXT ·memchrNotWordAVX2(SB), NOSPLIT, $0-32
	// Load parameters
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX

	// Empty check
	TESTQ   DX, DX
	JZ      notword_not_found

	// Save start for offset calculation and compute the end pointer
	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      notword_tail_loop          // sub-32 direct call: scalar

	SETUPTABLES

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      notword_loop32_entry

notword_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	NOTWORD(Y4, Y4, Y8, Y9)
	NOTWORD(Y5, Y5, Y8, Y9)
	NOTWORD(Y6, Y6, Y8, Y9)
	NOTWORD(Y7, Y7, Y8, Y9)

	VPOR    Y4, Y5, Y8
	VPOR    Y6, Y7, Y9
	VPOR    Y8, Y9, Y8
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
	NOTWORD(Y4, Y4, Y8, Y9)
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
	NOTWORD(Y4, Y4, Y8, Y9)
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	JMP     notword_not_found

notword_tail_loop:
	MOVBLZX (SI), AX

	// Check if NOT a word char (inverse of word_tail_loop logic)
	// Word chars: [A-Z], [a-z], [0-9], '_'
	// If it's a word char, skip; if not, found

	// Check [A-Z]: 65 <= b <= 90
	CMPB    AL, $65
	JB      notword_check_lower_gap
	CMPB    AL, $90
	JBE     notword_is_word           // It's [A-Z], skip

notword_check_lower_gap:
	// Check [a-z]: 97 <= b <= 122
	CMPB    AL, $97
	JB      notword_check_digit_gap
	CMPB    AL, $122
	JBE     notword_is_word           // It's [a-z], skip

notword_check_digit_gap:
	// Check [0-9]: 48 <= b <= 57
	CMPB    AL, $48
	JB      notword_check_underscore_gap
	CMPB    AL, $57
	JBE     notword_is_word           // It's [0-9], skip

notword_check_underscore_gap:
	// Check '_': 95
	CMPB    AL, $95
	JE      notword_is_word           // It's '_', skip

	// Not a word char - found!
	JMP     notword_found_scalar

notword_is_word:
	// It's a word char, continue searching
	INCQ    SI
	CMPQ    SI, R8
	JB      notword_tail_loop

notword_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

notword_found_in_block:
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
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
