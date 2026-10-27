package utils

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/utils/v2/simd"
	"github.com/stretchr/testify/require"
)

// IsASCII, IndexAny2 and IndexAny3 run their SWAR loops on inputs of simd.MinLen bytes or more only without
// the AVX2 kernels, so an ordinary run on AVX2 hardware never reaches them; Test_PortableTier does.

const (
	portableTierEnv = "UTILS_TEST_PORTABLE_TIER_CHILD"
	noAVX2          = "cpu.avx2=off"
)

// portableTierTests are checked by name, so renaming one fails loudly instead of dropping it.
var portableTierTests = []string{
	"Test_PortableTier_IsPortable",
	"Test_PortableTier_PositionSweep",
	"Test_IsASCII_Fixed",
	"Test_IndexAny_Fixed",
	"Test_Search_Randomized",
	"Test_Search_NamedTypes_SIMDDispatch",
	"Test_AppendUnescape_LongAndBoundaryInputs",
	"Test_AppendQueryUnescape",
}

func isPortableTierChild() bool {
	return os.Getenv(portableTierEnv) != ""
}

func coverDirArg() string {
	for _, arg := range os.Args[1:] {
		if dir, ok := strings.CutPrefix(arg, "-test.gocoverdir="); ok {
			return dir
		}
	}
	return ""
}

// Test_PortableTier re-runs the tier-dependent tests in a child with AVX2 off, merging its coverage.
func Test_PortableTier(t *testing.T) {
	t.Parallel()
	if isPortableTierChild() {
		t.Skip("this is the portable-tier run")
	}
	if !simd.Accelerated() {
		t.Skip("the AVX2 tier is not in use here, so every test already runs the portable tier")
	}

	exe, err := os.Executable()
	require.NoError(t, err)
	args := []string{"-test.run=^(" + strings.Join(portableTierTests, "|") + ")$", "-test.v", "-test.count=1", "-test.shuffle=off"}
	if dir := coverDirArg(); dir != "" {
		args = append(args, "-test.gocoverdir="+dir)
	}
	if deadline, ok := t.Deadline(); ok {
		// Our timeout ends this process, not the child; pass on what is left.
		if left := time.Until(deadline); left > 0 {
			args = append(args, "-test.timeout="+left.String())
		}
	}
	cmd := exec.CommandContext(t.Context(), exe, args...) // #nosec G204 -- the child is this test binary
	godebug := noAVX2
	if existing := os.Getenv("GODEBUG"); existing != "" {
		godebug = existing + "," + noAVX2
	}
	cmd.Env = append(os.Environ(), "GODEBUG="+godebug, portableTierEnv+"=1")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "the portable-tier run failed:\n%s", out)
	for _, name := range portableTierTests {
		require.Contains(t, string(out), "--- PASS: "+name+" (", "%s did not run in the portable-tier run:\n%s", name, out)
	}
}

// Test_PortableTier_IsPortable fails the child if AVX2 is still on.
func Test_PortableTier_IsPortable(t *testing.T) {
	t.Parallel()
	if !isPortableTierChild() {
		t.Skip("only meaningful in the portable-tier run")
	}
	require.False(t, simd.Accelerated(), "%s did not switch the AVX2 kernels off", noAVX2)
}

// Test_PortableTier_PositionSweep moves the byte of interest through every position of every length.
func Test_PortableTier_PositionSweep(t *testing.T) {
	t.Parallel()
	const maxLen = 130
	clean := bytes.Repeat([]byte("abcdefg."), maxLen/8+1)[:maxLen]
	for n := 0; n <= maxLen; n++ {
		base := clean[:n]
		if !IsASCII(base) || !IsASCII(string(base)) {
			t.Fatalf("IsASCII(%d clean bytes) = false", n)
		}
		if got := IndexAny2(base, ',', '"'); got != -1 {
			t.Fatalf("IndexAny2(%d clean bytes) = %d", n, got)
		}
		if got := IndexAny3(base, ',', '"', '\\'); got != -1 {
			t.Fatalf("IndexAny3(%d clean bytes) = %d", n, got)
		}
		for pos := range n {
			for _, hi := range []byte{0x80, 0xc3, 0xff} {
				b := bytes.Clone(base)
				b[pos] = hi
				if IsASCII(b) || IsASCII(string(b)) {
					t.Fatalf("IsASCII(%d bytes with %#x at %d) = true", n, hi, pos)
				}
			}
			for _, planted := range []byte{',', '"', '\\', 0x00, 0x80, 0xff} {
				b := bytes.Clone(base)
				b[pos] = planted
				s := string(b)
				if got, want := IndexAny2(b, planted, '"'), refIndexAny2(s, planted, '"'); got != want {
					t.Fatalf("IndexAny2(%d bytes with %#x at %d) = %d, want %d", n, planted, pos, got, want)
				}
				if got, want := IndexAny3(s, 'x', planted, 'y'), refIndexAny3(s, 'x', planted, 'y'); got != want {
					t.Fatalf("IndexAny3(%d bytes with %#x at %d) = %d, want %d", n, planted, pos, got, want)
				}
				// A later hit never moves the answer.
				if pos+1 < n {
					b[n-1] = planted
					if got := IndexAny2(b, planted, planted); got != pos {
						t.Fatalf("IndexAny2(%d bytes with %#x at %d and at the end) = %d", n, planted, pos, got)
					}
				}
			}
		}
	}
}
