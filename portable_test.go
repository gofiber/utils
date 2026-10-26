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

// IsASCII, IndexAny2 and IndexAny3 hand inputs of simd.MinLen bytes or more to
// the AVX2 kernels where there are any, and scan them with their own SWAR
// loops where there are not. On an amd64 machine with AVX2, which is every
// machine the amd64 CI jobs and most developers use, the SWAR loops for long
// inputs would therefore never run, and the one tier that is not exercised is
// the one every other CPU runs.
//
// Test_PortableTier closes that gap from inside an ordinary test run: it starts
// this test binary again with the AVX2 feature switched off through the
// GODEBUG setting golang.org/x/sys/cpu reads, and has the child run the tests
// that depend on the tier. When the run measures coverage the child writes its
// counters to the same directory as the parent, so the report is one that
// includes the portable tier.

const (
	portableTierEnv = "UTILS_TEST_PORTABLE_TIER_CHILD"
	noAVX2          = "cpu.avx2=off"
)

// portableTierTests are the tests the child runs. They are checked by name in
// its output, so renaming one without updating this list fails loudly instead
// of quietly dropping it from the portable run.
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

// coverDirArg returns the directory go test asked this binary to write its
// coverage counters to, or "" when it is not measuring coverage.
func coverDirArg() string {
	for _, arg := range os.Args[1:] {
		if dir, ok := strings.CutPrefix(arg, "-test.gocoverdir="); ok {
			return dir
		}
	}
	return ""
}

// Test_PortableTier re-runs the tier-dependent tests with AVX2 disabled.
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
		// The timeout of this run ends this process, not the child, which would
		// go on to its own default of ten minutes. It gets what is left of ours.
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

// Test_PortableTier_IsPortable fails the child run if it somehow still has the
// AVX2 kernels, which would make the whole exercise a repeat of the parent's.
func Test_PortableTier_IsPortable(t *testing.T) {
	t.Parallel()
	if !isPortableTierChild() {
		t.Skip("only meaningful in the portable-tier run")
	}
	require.False(t, simd.Accelerated(), "%s did not switch the AVX2 kernels off", noAVX2)
}

// Test_PortableTier_PositionSweep walks inputs of every length through the
// scan loops' word, block and tail paths with the byte of interest at every
// position. It runs in both tiers, and in the portable one it is what pins the
// SWAR loops for lengths the AVX2 kernels otherwise take.
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
				// A second hit after the first never moves the answer.
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
