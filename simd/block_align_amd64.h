// BLOCKALIGN / BLOCKALIGN_LONG start the block loop at byte 18 / 27 of a 32-byte
// window: PCALIGN, then two / three 9-byte NOPs before the label (once per call).
// Skylake-derived cores' JCC-erratum mitigation keeps any jump or fused cmp+jcc
// crossing or ending on a 32-byte boundary out of the uop cache; these
// layout-dependent offsets keep every pair clear of it, which
// block_align_amd64_test.go verifies on the linked binary.
#define BLOCKALIGN \
	PCALIGN $32; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00

#define BLOCKALIGN_LONG \
	BLOCKALIGN; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00
