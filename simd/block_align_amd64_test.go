//go:build amd64

package simd

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

// BLOCKALIGN / BLOCKALIGN_LONG keep every macro-fused compare-and-branch pair of the block
// loops from crossing or ending on a 32-byte boundary (JCC erratum on Skylake-derived cores:
// a uop-cache miss per iteration). This test reads each kernel's code from the running binary,
// so a branch moved onto a boundary fails here; the log gives each pair's offset from the loop
// head. Branches are matched on exact encodings (fusedPairs), not decoded: a kernel with a new
// branch shape must add its encoding.

var blockLoopKernels = map[string]any{
	"memchr2AVX2":       memchr2AVX2,
	"memchr3AVX2":       memchr3AVX2,
	"memchrPairAVX2":    memchrPairAVX2,
	"memchrDigitAVX2":   memchrDigitAVX2,
	"memchrWordAVX2":    memchrWordAVX2,
	"memchrNotWordAVX2": memchrNotWordAVX2,
	"isASCIIAVX2":       isASCIIAVX2,
	"firstNonASCIIAVX2": firstNonASCIIAVX2,
	"countNonASCIIAVX2": countNonASCIIAVX2,
}

var fusedPairs = []struct {
	name     string
	prefix   []byte
	backEdge bool
}{
	{name: "TESTL CX, CX", prefix: []byte{0x85, 0xC9}},
	{name: "TESTL AX, AX", prefix: []byte{0x85, 0xC0}},
	{name: "DECQ R13", prefix: []byte{0x49, 0xFF, 0xCD}},
	{name: "CMPQ SI, R11", prefix: []byte{0x4C, 0x39, 0xDE}, backEdge: true},
}

var nop9 = []byte{0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00}

var nopForms = [][]byte{
	{0x66, 0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	nop9,
	{0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x0F, 0x1F, 0x80, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x40, 0x00},
	{0x0F, 0x1F, 0x00},
	{0x66, 0x90},
	{0x90},
}

func nopLen(b []byte) int {
	for _, form := range nopForms {
		if bytes.HasPrefix(b, form) {
			return len(form)
		}
	}
	return 0
}

// kernelBody returns the kernel's code: an ABI0 kernel's Go func value is a wrapper
// whose first CALL leads there.
func kernelBody(t *testing.T, name string, fn any) ([]byte, uintptr) {
	t.Helper()
	wrapper := reflect.ValueOf(fn).UnsafePointer()
	w := unsafe.Slice((*byte)(wrapper), 128) //nolint:gosec // reading this binary's own text segment
	call := bytes.IndexByte(w, 0xE8)
	require.GreaterOrEqual(t, call, 0, "%s: no CALL in the ABI wrapper", name)
	rel := int32(binary.LittleEndian.Uint32(w[call+1:]))
	body := unsafe.Add(wrapper, call+5+int(rel))
	code := unsafe.Slice((*byte)(body), 2048) //nolint:gosec // reading this binary's own text segment
	require.True(t, bytes.HasPrefix(code, []byte{0x48, 0x8B, 0x74, 0x24, 0x08}),
		"%s: wrapper CALL does not lead to the kernel body", name)
	return code, uintptr(body)
}

// Test_BlockLoopBranchesStayInWindow checks that no fused pair crosses or ends on a 32-byte boundary.
func Test_BlockLoopBranchesStayInWindow(t *testing.T) {
	signature := slices.Concat(nop9, nop9)
	for name, fn := range blockLoopKernels {
		code, base := kernelBody(t, name, fn)

		// The signature can also match in PCALIGN padding, so the head is the first non-NOP after it.
		head := bytes.Index(code, signature)
		require.GreaterOrEqual(t, head, 0, "%s: BLOCKALIGN padding not found", name)
		head += len(signature)
		for n := nopLen(code[head:]); n > 0; n = nopLen(code[head:]) {
			head += n
		}

		pairs := 0
		closed := false
		for i := head; i < len(code)-6 && !closed; i++ {
			for _, fp := range fusedPairs {
				if !bytes.HasPrefix(code[i:], fp.prefix) {
					continue
				}
				j := i + len(fp.prefix)
				var size int
				switch {
				case code[j] >= 0x70 && code[j] <= 0x7F:
					size = 2
				case code[j] == 0x0F && code[j+1] >= 0x80 && code[j+1] <= 0x8F:
					size = 6
				default:
					continue
				}
				start := base + uintptr(i)
				end := base + uintptr(j+size-1)
				t.Logf("%s: head mod 32 = %d, fused %s/Jcc at head+%d..head+%d", name, (base+uintptr(head))%32, fp.name, i-head, j+size-1-head)
				if start/32 != end/32 {
					t.Errorf("%s: fused %s/Jcc at %#x-%#x crosses a 32-byte boundary", name, fp.name, start, end)
				}
				if end%32 == 31 {
					t.Errorf("%s: fused %s/Jcc at %#x-%#x ends on a 32-byte boundary", name, fp.name, start, end)
				}
				pairs++
				if fp.backEdge {
					closed = true
				}
				i = j + size - 1
				break
			}
		}
		require.True(t, closed, "%s: block loop back-edge not found", name)
		require.GreaterOrEqual(t, pairs, 1, "%s: no fused branch found in the block loop", name)
	}
}
