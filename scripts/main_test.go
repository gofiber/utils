package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// captureStderr runs fn with os.Stderr pointing at a pipe and returns what fn
// wrote to it, which is where the tool logs and where the child go test writes.
// The tests of this package do not run in parallel, so swapping it is safe.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stderr
	os.Stderr = w
	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) //nolint:errcheck // the pipe only ends when the writer closes
		done <- buf.String()
	}()
	func() {
		defer func() {
			os.Stderr = original
			_ = w.Close() //nolint:errcheck // nothing to do about a failed close of a pipe
		}()
		fn()
	}()
	return <-done
}

// writeMiniModule creates a module with one benchmark in a directory of its own
// and returns the directory. It has no dependencies and no go directive, so
// `go test` in it needs neither the network nor another toolchain, only the go
// command, without which the tests that run it are skipped.
//
// The benchmark sleeps for a millisecond: one iteration of an empty body can
// measure as 0 ns on a coarse clock, and `go test` prints no ns/op column then,
// which is a row this tool does not take.
func writeMiniModule(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("the go command is needed to run a benchmark: %v", err)
	}
	dir := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":       "module example.com/mini\n",
		"mini.go":      "package mini\n",
		"mini_test.go": "package mini\n\nimport (\n\t\"testing\"\n\t\"time\"\n)\n\nfunc Benchmark_Mini(b *testing.B) {\n\tfor i := 0; i < b.N; i++ {\n\t\ttime.Sleep(time.Millisecond)\n\t}\n}\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// Neither a parent workspace nor flags meant for the real module (-mod=vendor
	// say) may reach the temporary one. An empty GOFLAGS counts as unset, so the
	// settings of `go env -w` would still apply without GOENV=off.
	t.Setenv("GOWORK", "off")
	t.Setenv("GOFLAGS", "")
	t.Setenv("GOENV", "off")
	return dir
}

// miniTable is a README table for the mini module, for this machine's architecture: a placeholder row whose
// iteration count, 777, no real run of -benchtime=1x can produce.
func miniTable(command string) string {
	return "Environment:\ngoos: plan9\ngoarch: " + runtime.GOARCH + "\npkg: example.com/mini\ncpu: an old cpu\n\n" +
		"```text\n" + command + "\n\nBenchmark_Mini-8   777   9.000  ns/op   8  B/op   1  allocs/op\n```\n"
}

const miniCommand = "// go test -run=^$ -bench=Benchmark_Mini -benchmem -benchtime=1x -count=3 ."

func writeReadme(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestRunHelp(t *testing.T) {
	var err error
	out := captureStderr(t, func() { err = run([]string{"-h"}, "does-not-exist.md") })
	if err != nil {
		t.Fatalf("help is not an error, got %v", err)
	}
	if out != usage {
		t.Errorf("help printed %q, want the usage text", out)
	}
}

func TestRunRejectsBadArguments(t *testing.T) {
	for _, args := range [][]string{{"-bogus"}, {"-block"}, {"-block", "0"}, {"-block", "2"}, {"-update=1"}} {
		err := run(args, "does-not-exist.md")
		if err == nil {
			t.Errorf("%v was accepted", args)
			continue
		}
		if !strings.Contains(err.Error(), usage) {
			t.Errorf("%v: the error must carry the usage text, got %q", args, err)
		}
	}
}

func TestRunFormatsTheReadme(t *testing.T) {
	readme := writeReadme(t, t.TempDir(), fixture)
	if err := run(nil, readme); err != nil {
		t.Fatal(err)
	}
	out := readFile(t, readme)
	if !strings.Contains(out, "Benchmark_Kept/a-12  100  1.000  ns/op  0  B/op  0  allocs/op") {
		t.Errorf("the rows were not aligned:\n%s", out)
	}
	// Formatting again changes nothing.
	if err := run(nil, readme); err != nil {
		t.Fatal(err)
	}
	if again := readFile(t, readme); again != out {
		t.Errorf("formatting twice differs:\n%s\n---\n%s", out, again)
	}
}

func TestRunReportsAMissingReadme(t *testing.T) {
	if err := run(nil, filepath.Join(t.TempDir(), "missing.md")); err == nil {
		t.Error("a missing README was accepted")
	}
}

func TestRunUpdateMeasuresTableWithItsOwnCommand(t *testing.T) {
	dir := writeMiniModule(t)
	readme := writeReadme(t, dir, "# Mini\n\n"+miniTable(miniCommand))

	var err error
	logged := captureStderr(t, func() { err = run([]string{"-update"}, readme) })
	if err != nil {
		t.Fatalf("%v\nlog:\n%s", err, logged)
	}
	out := readFile(t, readme)
	for _, want := range []string{"goos: " + runtime.GOOS, "goarch: " + runtime.GOARCH} {
		if !strings.Contains(out, want) {
			t.Errorf("the header lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "plan9") || strings.Contains(out, "777") {
		t.Errorf("the placeholder survived the refresh:\n%s", out)
	}
	var row []string
	for line := range strings.Lines(out) {
		if strings.HasPrefix(line, "Benchmark_Mini") {
			row = strings.Fields(line)
		}
	}
	if len(row) < 8 || row[1] != "1" || row[3] != "ns/op" {
		t.Errorf("the row was not replaced by a measured one: %q", row)
	}
	if !strings.Contains(logged, "table 1 (") || !strings.Contains(logged, "1 rows refreshed") {
		t.Errorf("the progress log lacks the table summary:\n%s", logged)
	}
}

func TestRunUpdateReportsAFailingCommand(t *testing.T) {
	dir := writeMiniModule(t)
	readme := writeReadme(t, dir, miniTable("// go test ./no-such-package"))

	var err error
	captureStderr(t, func() { err = run([]string{"-update"}, readme) })
	if err == nil || !strings.Contains(err.Error(), "go test ./no-such-package") {
		t.Fatalf("got %v, want the failing command in the error", err)
	}
	if out := readFile(t, readme); !strings.Contains(out, "plan9") || !strings.Contains(out, "777") {
		t.Errorf("a failed run must leave the README as it was:\n%s", out)
	}
}

func TestRunBenchStopsWithItsContext(t *testing.T) {
	// A context that is canceled already keeps the child from starting, and the
	// error names the command.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := runBench(ctx, t.TempDir(), []string{"go", "version"})
	if err == nil || !strings.Contains(err.Error(), "go version") {
		t.Errorf("got %v, want the canceled command in the error", err)
	}
}

// TestRunBenchKillsARunningChild is the Ctrl-C case: the context is canceled while
// the child runs, and the child goes with it. The child is this test binary, run
// again to do nothing but announce itself and sleep.
func TestRunBenchKillsARunningChild(t *testing.T) {
	if marker := os.Getenv("SCRIPTS_TEST_CHILD_MARKER"); marker != "" {
		err := os.WriteFile(marker, []byte("up"), 0o600) // #nosec G703 -- the marker is the temporary file the parent test named
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Minute)
		return
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "started")
	t.Setenv("SCRIPTS_TEST_CHILD_MARKER", marker)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := runBench(ctx, t.TempDir(), []string{exe, "-test.run=^TestRunBenchKillsARunningChild$"})
		done <- err
	}()

	// Cancel only once the child is up, or the kill is not what is tested.
	for deadline := time.Now().Add(time.Minute); ; time.Sleep(10 * time.Millisecond) {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the child never started")
		}
	}
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Error("a child that was killed ended without an error")
		}
	case <-time.After(time.Minute):
		t.Fatal("runBench did not return after its context was canceled")
	}
}

func TestMainEndsWithStatusOneOnFailure(t *testing.T) {
	originalArgs, originalExit := os.Args, exit
	t.Cleanup(func() { os.Args, exit = originalArgs, originalExit })
	var statuses []int
	exit = func(code int) { statuses = append(statuses, code) }

	os.Args = []string{"scripts", "-no-such-flag"}
	out := captureStderr(t, main)
	if len(statuses) != 1 || statuses[0] != 1 {
		t.Errorf("exit statuses %v, want [1]", statuses)
	}
	if !strings.Contains(out, `unknown argument "-no-such-flag"`) {
		t.Errorf("the error was not printed:\n%s", out)
	}

	statuses = nil
	os.Args = []string{"scripts", "-h"}
	out = captureStderr(t, main)
	if len(statuses) != 0 {
		t.Errorf("help exited with %v", statuses)
	}
	if out != usage {
		t.Errorf("help printed %q", out)
	}
}
