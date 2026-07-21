// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, and the
// scalar tail loops replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees.

//go:build amd64

#include "textflag.h"

// func memchr2AVX2(haystack []byte, needle1, needle2 byte) int
//
// AVX2 implementation of memchr2 that searches for either of two bytes.
// Uses parallel comparison with two broadcast vectors and combines results
// with VPOR. The final 1-31 bytes are handled with one overlapping 32-byte
// vector at the buffer end: lanes before it were already resolved as
// non-matching, so the first set lane of the overlap falls in the new bytes
// (the same argument the SWAR fallbacks use for their overlapping word at
// n-8). Inputs shorter than 32 bytes — unreachable through the Go dispatch,
// which routes them to the SWAR fallback — take a scalar loop instead.
//
// Parameters (FP offsets, ABI0; the byte arguments are packed into
// adjacent 1-byte slots after the 24-byte slice header, then padded to
// align the return value):
//   haystack_base+0(FP)  - pointer (8 bytes)
//   haystack_len+8(FP)   - length (8 bytes)
//   haystack_cap+16(FP)  - capacity (8 bytes, unused)
//   needle1+24(FP)       - first needle (1 byte)
//   needle2+25(FP)       - second needle (1 byte, then 6 bytes padding)
//   ret+32(FP)           - return value (8 bytes)
//
// Total argument frame: 40 bytes (24 slice + 1 + 1 + 6 padding + 8 ret)
TEXT ·memchr2AVX2(SB), NOSPLIT, $0-40
	// Load parameters
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX
	MOVBQZX needle1+24(FP), AX
	MOVBQZX needle2+25(FP), BX

	// Empty check
	TESTQ   DX, DX
	JZ      not_found2

	// Broadcast both needles
	VMOVD   AX, X0
	VPBROADCASTB X0, Y0                  // Y0 = [needle1 × 32]
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1                  // Y1 = [needle2 × 32]

	// Save start pointer
	MOVQ    SI, DI

	// Calculate end
	LEAQ    (SI)(DX*1), R8

loop32_2:
	LEAQ    32(SI), R9
	CMPQ    R9, R8
	JA      handle_tail2

	// Load chunk
	VMOVDQU (SI), Y2

	// Compare with both needles
	VPCMPEQB Y0, Y2, Y3                  // Y3 = matches with needle1
	VPCMPEQB Y1, Y2, Y4                  // Y4 = matches with needle2

	// Combine results with OR (match if either needle found)
	VPOR    Y3, Y4, Y3                   // Y3 = needle1_matches | needle2_matches

	// Extract mask
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector2

	ADDQ    $32, SI
	JMP     loop32_2

handle_tail2:
	// Fewer than 32 bytes remain.
	CMPQ    SI, R8
	JAE     not_found2
	CMPQ    DX, $32
	JB      tail_loop2                   // sub-32 direct call: scalar

	// Rescan the final 32-byte window ending at the buffer end; lanes
	// before SI are known non-matching, so BSF stays exact.
	LEAQ    -32(R8), SI
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
	CMPB    CL, AL                       // Compare with needle1
	JE      found_scalar2
	CMPB    CL, BL                       // Compare with needle2
	JE      found_scalar2

	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop2

not_found2:
	MOVQ    $-1, AX
	MOVQ    AX, ret+32(FP)
	VZEROUPPER
	RET

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
// AVX2 implementation of memchr3 that searches for any of three bytes.
// Same structure as memchr2AVX2 (including the overlapping vector tail),
// with a third broadcast/compare.
//
// Parameters (FP offsets, ABI0; byte arguments packed into adjacent 1-byte
// slots after the slice header, then padded):
//   haystack_base+0(FP)  - pointer (8 bytes)
//   haystack_len+8(FP)   - length (8 bytes)
//   haystack_cap+16(FP)  - capacity (8 bytes, unused)
//   needle1+24(FP)       - first needle (1 byte)
//   needle2+25(FP)       - second needle (1 byte)
//   needle3+26(FP)       - third needle (1 byte, then 5 bytes padding)
//   ret+32(FP)           - return value (8 bytes)
//
// Total argument frame: 40 bytes (24 slice + 3 + 5 padding + 8 ret)
TEXT ·memchr3AVX2(SB), NOSPLIT, $0-40
	// Load parameters
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX
	MOVBQZX needle1+24(FP), AX
	MOVBQZX needle2+25(FP), BX
	MOVBQZX needle3+26(FP), R10

	// Empty check
	TESTQ   DX, DX
	JZ      not_found3

	// Broadcast all three needles
	VMOVD   AX, X0
	VPBROADCASTB X0, Y0                  // Y0 = [needle1 × 32]
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1                  // Y1 = [needle2 × 32]
	VMOVD   R10, X2
	VPBROADCASTB X2, Y2                  // Y2 = [needle3 × 32]

	// Save start pointer
	MOVQ    SI, DI

	// Calculate end
	LEAQ    (SI)(DX*1), R8

loop32_3:
	LEAQ    32(SI), R9
	CMPQ    R9, R8
	JA      handle_tail3

	// Load chunk
	VMOVDQU (SI), Y3

	// Compare with all three needles
	VPCMPEQB Y0, Y3, Y4                  // Y4 = matches with needle1
	VPCMPEQB Y1, Y3, Y5                  // Y5 = matches with needle2
	VPCMPEQB Y2, Y3, Y6                  // Y6 = matches with needle3

	// Combine results: match if any needle found
	VPOR    Y4, Y5, Y4                   // Y4 = needle1 | needle2
	VPOR    Y4, Y6, Y4                   // Y4 = (needle1 | needle2) | needle3

	// Extract mask
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector3

	ADDQ    $32, SI
	JMP     loop32_3

handle_tail3:
	// Fewer than 32 bytes remain; overlapping vector tail, see memchr2AVX2.
	CMPQ    SI, R8
	JAE     not_found3
	CMPQ    DX, $32
	JB      tail_loop3                   // sub-32 direct call: scalar

	LEAQ    -32(R8), SI
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
	CMPB    CL, AL                       // Compare with needle1
	JE      found_scalar3
	CMPB    CL, BL                       // Compare with needle2
	JE      found_scalar3
	CMPB    CL, R10B                     // Compare with needle3
	JE      found_scalar3

	INCQ    SI
	CMPQ    SI, R8
	JB      tail_loop3

not_found3:
	MOVQ    $-1, AX
	MOVQ    AX, ret+32(FP)
	VZEROUPPER
	RET

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
// AVX2 implementation of paired-byte search. Finds positions where byte1
// appears at position i and byte2 appears at position i+offset.
//
// Algorithm:
//  1. Broadcast byte1 to YMM0, byte2 to YMM1
//  2. For each 32-byte chunk at position p:
//     - Load haystack[p:p+32], compare with byte1 → mask1
//     - Load haystack[p+offset:p+offset+32], compare with byte2 → mask2
//     - AND mask1, mask2 → combined mask (both bytes at correct distance)
//  3. First set bit in combined mask is the answer
//  4. The final 1-31 pair positions are rescanned with one overlapping
//     vector pair whose second load ends exactly at the buffer end; pair
//     positions before it are known non-matching, so BSF stays exact.
//     Inputs with fewer than offset+32 bytes — unreachable through the Go
//     dispatch — take a scalar loop instead.
//
// Parameters (FP offsets, ABI0; the byte arguments are packed into
// adjacent 1-byte slots, then padded to align the int argument):
//   haystack_base+0(FP)  - pointer (8 bytes)
//   haystack_len+8(FP)   - length (8 bytes)
//   haystack_cap+16(FP)  - capacity (8 bytes, unused)
//   byte1+24(FP)         - first byte to find (1 byte)
//   byte2+25(FP)         - second byte to find (1 byte, then 6 bytes padding)
//   offset+32(FP)        - distance between byte1 and byte2 (8 bytes)
//   ret+40(FP)           - return value (8 bytes)
//
// Total argument frame: 48 bytes (24 slice + 1 + 1 + 6 padding + 8 + 8 ret)
TEXT ·memchrPairAVX2(SB), NOSPLIT, $0-48
	// Load parameters
	MOVQ    haystack_base+0(FP), SI      // SI = haystack pointer
	MOVQ    haystack_len+8(FP), DX       // DX = haystack length
	MOVBQZX byte1+24(FP), AX             // AX = byte1
	MOVBQZX byte2+25(FP), BX             // BX = byte2
	MOVQ    offset+32(FP), R10           // R10 = offset

	// Empty check and bounds check
	TESTQ   DX, DX
	JZ      not_found_pair
	CMPQ    R10, DX                      // offset >= length?
	JAE     not_found_pair

	// Broadcast both bytes
	VMOVD   AX, X0
	VPBROADCASTB X0, Y0                  // Y0 = [byte1 × 32]
	VMOVD   BX, X1
	VPBROADCASTB X1, Y1                  // Y1 = [byte2 × 32]

	// Save start pointer
	MOVQ    SI, DI                       // DI = haystack start

	// Calculate limits:
	// - We need 32 bytes at position p (for byte1 check)
	// - We need 32 bytes at position p+offset (for byte2 check)
	// - So we need p + offset + 32 <= length
	// - Which means p <= length - offset - 32
	MOVQ    DX, R8                       // R8 = length
	SUBQ    R10, R8                      // R8 = length - offset
	SUBQ    $32, R8                      // R8 = length - offset - 32 (last valid p)
	CMPQ    R8, $0
	JL      scalar_entry_pair            // sub-(offset+32) direct call: scalar

	// R8 now contains the last valid position for the vector loop
	ADDQ    SI, R8                       // R8 = absolute last valid base

loop32_pair:
	// Check if SI <= R8 (still have room for a full vector)
	CMPQ    SI, R8
	JA      vector_tail_pair

	// Load 32 bytes at position p (for byte1)
	VMOVDQU (SI), Y2                     // Y2 = haystack[p:p+32]

	// Load 32 bytes at position p+offset (for byte2)
	VMOVDQU (SI)(R10*1), Y3              // Y3 = haystack[p+offset:p+offset+32]

	// Compare byte1
	VPCMPEQB Y0, Y2, Y4                  // Y4 = positions where haystack[p+k] == byte1

	// Compare byte2
	VPCMPEQB Y1, Y3, Y5                  // Y5 = positions where haystack[p+offset+k] == byte2

	// AND the results: bit k is set only if:
	// - haystack[p+k] == byte1
	// - haystack[p+offset+k] == byte2
	// This means the pair (byte1, byte2) appears at (p+k, p+k+offset)
	VPAND   Y4, Y5, Y4                   // Y4 = combined mask

	// Extract mask
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector_pair

	// No match in this chunk, advance by 32 bytes
	ADDQ    $32, SI
	JMP     loop32_pair

vector_tail_pair:
	// Fewer than 32 pair positions remain past SI. If any do, rescan the
	// final 32-position window based at R8 = base+len-offset-32: its second
	// load ends exactly at the buffer end, and pair positions before SI are
	// known non-matching.
	LEAQ    (DI)(DX*1), R9               // R9 = buffer end
	SUBQ    R10, R9                      // R9 = end of valid pair positions
	CMPQ    SI, R9
	JAE     not_found_pair
	MOVQ    R8, SI                       // SI = base of the final window
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
	// Scalar path for direct calls with fewer than offset+32 bytes.
	LEAQ    (DI)(DX*1), R8               // R8 = end pointer

tail_loop_pair:
	// Check if SI + offset < end
	LEAQ    (SI)(R10*1), R9              // R9 = SI + offset
	CMPQ    R9, R8
	JAE     not_found_pair               // If SI + offset >= end, we're done

	// Load and compare both bytes
	MOVBLZX (SI), CX                     // CX = haystack[SI]
	CMPB    CL, AL                       // Compare with byte1
	JNE     tail_next_pair               // Not byte1, skip

	// byte1 matches, check byte2
	MOVBLZX (SI)(R10*1), CX              // CX = haystack[SI + offset]
	CMPB    CL, BL                       // Compare with byte2
	JE      found_scalar_pair            // Both match!

tail_next_pair:
	INCQ    SI
	JMP     tail_loop_pair

not_found_pair:
	MOVQ    $-1, AX
	MOVQ    AX, ret+40(FP)
	VZEROUPPER
	RET

found_in_vector_pair:
	// Match found in vector! CX contains mask
	BSFL    CX, CX                       // CX = index of first set bit (0-31)
	SUBQ    DI, SI                       // SI = offset from start to current chunk
	ADDQ    SI, CX                       // CX = absolute index
	MOVQ    CX, ret+40(FP)
	VZEROUPPER
	RET

found_scalar_pair:
	SUBQ    DI, SI
	MOVQ    SI, ret+40(FP)
	VZEROUPPER
	RET
