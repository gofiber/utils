package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const fixture = `# Title

` + "```text" + `
sample output that is not a benchmark table
` + "```" + `

Environment:
goos: darwin
goarch: arm64
pkg: github.com/gofiber/utils/v2
cpu: Apple M2 Pro

` + "```text" + `
// go test ./... -benchmem -run=^$ -bench=Benchmark_ -count=1

# Group
Benchmark_Kept/a-12      100    1.000  ns/op     0  B/op   0  allocs/op
Benchmark_Gone/b-12      200    2.000  ns/op     8  B/op   1  allocs/op
Benchmark_Both/c-12      300    3.000  ns/op     0  B/op   0  allocs/op
Benchmark_Wide/d-12      400    4.000  ns/op     0  B/op   0  allocs/op
` + "```" + `
`

func TestParseBlocks(t *testing.T) {
	blocks, err := parseBlocks(strings.Split(fixture, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	// The plain ```text block above carries no Environment header and must not
	// shift the numbering that -block uses.
	if len(blocks) != 1 {
		t.Fatalf("got %d tables, want 1", len(blocks))
	}
	b := blocks[0]
	if b.goos != "darwin" || b.goarch != "arm64" {
		t.Errorf("got %s/%s, want darwin/arm64", b.goos, b.goarch)
	}
	if b.pkg != "github.com/gofiber/utils/v2" {
		t.Errorf("got pkg %q", b.pkg)
	}
	if want := "go test ./... -benchmem -run=^$ -bench=Benchmark_ -count=1"; strings.Join(b.cmd, " ") != want {
		t.Errorf("got command %q, want %q", strings.Join(b.cmd, " "), want)
	}
}

func TestParseBlocksRejectsBrokenFences(t *testing.T) {
	for name, md := range map[string]string{
		"unclosed at end of file": "```text\nrow\n",
		"never closed before the next one": "```text\nrow\n" +
			"```text\nrow\n```\n",
	} {
		if _, err := parseBlocks(strings.Split(md, "\n")); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func TestApplyResults(t *testing.T) {
	lines := strings.Split(fixture, "\n")
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}

	applyResults(lines, blocks[0], 1, strings.Join([]string{
		"goos: linux",
		"goarch: amd64",
		"pkg: github.com/gofiber/utils/v2",
		"cpu: Intel(R) Xeon(R) Processor @ 2.80GHz",
		"Benchmark_Kept/a-4   	     300	         3.500 ns/op	      16 B/op	       2 allocs/op",
		// b.SetBytes adds an MB/s column that the table has no room for.
		"Benchmark_Wide/d-4   	     500	         5.500 ns/op	  181.82 MB/s	       0 B/op	       0 allocs/op",
		"Benchmark_New/e-4    	     400	         4.500 ns/op	       0 B/op	       0 allocs/op",
		// Benchmark_Both is measured by two packages, and only one of them is the
		// one this table's header names.
		"Benchmark_Both/c-4   	     100	        10.000 ns/op	       0 B/op	       0 allocs/op",
		"pkg: github.com/gofiber/utils/v2/simd",
		"Benchmark_Both/c-4   	     600	       600.000 ns/op	       0 B/op	       0 allocs/op",
		"PASS",
	}, "\n"), knownRows(lines, blocks))

	out := strings.Join(lines, "\n")
	for _, want := range []string{
		"Benchmark_Kept/a-4  300  3.500  ns/op  16  B/op  2  allocs/op",
		"Benchmark_Wide/d-4  500  5.500  ns/op  0  B/op  0  allocs/op",
		"goos: linux",
		"goarch: amd64",
		"cpu: Intel(R) Xeon(R) Processor @ 2.80GHz",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "MB/s") {
		t.Errorf("the MB/s column must not reach the table:\n%s", out)
	}
	if !strings.Contains(out, "Benchmark_Gone/b-12      200    2.000  ns/op     8  B/op   1  allocs/op") {
		t.Error("a row missing from the output must stay untouched")
	}
	if !strings.Contains(out, "Benchmark_Both/c-4  100  10.000  ns/op") {
		t.Errorf("a name in two packages must take the one the table names:\n%s", out)
	}
	if strings.Contains(out, "Benchmark_New") {
		t.Error("a new benchmark must be reported, not inserted")
	}
}

func TestApplyResultsNamesNoMachineItCannotRead(t *testing.T) {
	lines := strings.Split(fixture, "\n")
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}

	// arm64 Linux carries no model name, so `go test` prints no cpu line at all.
	applyResults(lines, blocks[0], 1, strings.Join([]string{
		"goos: linux",
		"goarch: arm64",
		"pkg: github.com/gofiber/utils/v2",
		"Benchmark_Kept/a-4   1   3.500 ns/op   0 B/op   0 allocs/op",
	}, "\n"), knownRows(lines, blocks))

	out := strings.Join(lines, "\n")
	if strings.Contains(out, "Apple M2 Pro") {
		t.Errorf("the header must not keep a machine that did not run:\n%s", out)
	}
	if !strings.Contains(out, "cpu: unknown") {
		t.Errorf("a cpu the run did not report must be named unknown:\n%s", out)
	}
}

func TestPickResolvesOnThePackageOfTheTable(t *testing.T) {
	root := benchResult{pkg: "m", ns: 1}
	simd := benchResult{pkg: "m/simd", ns: 2}

	if runs, ok := pick([]benchResult{root, root}, "m/other"); !ok || len(runs) != 2 {
		t.Error("one package means no ambiguity, whatever the table says")
	}
	if runs, ok := pick([]benchResult{root, simd}, "m/simd"); !ok || len(runs) != 1 || runs[0].pkg != "m/simd" {
		t.Errorf("got %v, want only the table's own package", runs)
	}
	if _, ok := pick([]benchResult{root, simd}, "m/third"); ok {
		t.Error("an unresolvable name must be reported, not guessed")
	}
}

func TestMedianPicksMiddleRun(t *testing.T) {
	_, results, odd := parseBenchOutput(strings.Join([]string{
		"Benchmark_X-4   1   9.000 ns/op   0 B/op   0 allocs/op",
		"Benchmark_X-4   1   1.000 ns/op   0 B/op   0 allocs/op",
		"Benchmark_X-4   1   5.000 ns/op   0 B/op   0 allocs/op",
		"Benchmark_Odd-4   1   1 x   0 B/op",
	}, "\n"))
	if got := median(results["Benchmark_X"]); got == nil || got.ns != 5 {
		t.Errorf("got %v, want the 5.000 ns/op run", got)
	}
	if median(results["Benchmark_missing"]) != nil {
		t.Error("an unknown benchmark must not resolve to a run")
	}
	// An even -count has no middle; the doc comment promises the slower of the two.
	two := []benchResult{{ns: 1}, {ns: 9}}
	if got := median(two); got == nil || got.ns != 9 {
		t.Errorf("got %v, want the slower of the two middles", got)
	}
	if len(odd) != 1 || odd[0] != "Benchmark_Odd-4" {
		t.Errorf("got %v, want the unusable row reported", odd)
	}
}

func TestSplitArgs(t *testing.T) {
	got, err := splitArgs(` go test -run=^$ . -bench='^Benchmark_(A|B)$' `)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"go", "test", "-run=^$", ".", "-bench=^Benchmark_(A|B)$"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("got %q, want %q", got, want)
	}

	for _, bad := range []string{"go build ./...", "make bench", `go test -bench='^X`, ""} {
		if _, err := splitArgs(bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}
}

func TestStripProcs(t *testing.T) {
	for in, want := range map[string]string{
		"Benchmark_X-12":                 "Benchmark_X",
		"Benchmark_X":                    "Benchmark_X",
		"Benchmark_IndexControl/a-16B/f": "Benchmark_IndexControl/a-16B/f",
		"Benchmark_IndexControl/a-16B-4": "Benchmark_IndexControl/a-16B",
	} {
		if got := stripProcs(in); got != want {
			t.Errorf("stripProcs(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseFlags(t *testing.T) {
	for _, args := range [][]string{{"-update"}, {"--update"}, {"-update", "-block", "2"}, {"-update", "-block=2"}} {
		update, block, err := parseFlags(args)
		if err != nil || !update || (len(args) > 1 && block != 2) {
			t.Errorf("parseFlags(%q) = %v, %d, %v", args, update, block, err)
		}
	}
	for _, args := range [][]string{
		{"-bogus"},
		{"-block", "2"},
		{"-update", "-block"},
		{"-update", "-block", "0"},
		{"-update", "-block", "x"},
		{"-update=false"},
		{"-h=x"},
	} {
		if _, _, err := parseFlags(args); err == nil {
			t.Errorf("parseFlags(%q) was accepted", args)
		}
	}
	if _, _, err := parseFlags([]string{"-h"}); err != errHelp { //nolint:errorlint,err113 // the sentinel is returned directly
		t.Errorf("-h returned %v, want errHelp", err)
	}
}

func TestReadWriteLinesRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "f.md")
	for in, want := range map[string]string{
		"a\nb\n": "a\nb\n",
		"a\nb":   "a\nb\n", // a missing final newline is added
		"a\n\n":  "a\n\n",
	} {
		if err := os.WriteFile(path, []byte(in), 0o600); err != nil {
			t.Fatal(err)
		}
		lines, err := readLines(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := writeLines(path, lines); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(path) // #nosec G304 -- path is the test's own temp file
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%q round-tripped to %q, want %q", in, got, want)
		}
	}
}

// The tool finds its tables by shape, so a README that drifts out of that shape
// must fail here rather than silently refreshing nothing.
func TestParseBlocksOnReadme(t *testing.T) {
	lines, err := readLines(filepath.Join("..", readmePath))
	if err != nil {
		t.Fatal(err)
	}
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 3 {
		t.Fatalf("got %d benchmark tables in README.md, want 3", len(blocks))
	}
	for i, b := range blocks {
		if len(b.cmd) == 0 {
			t.Errorf(`table %d has no "// go test ..." line`, i+1)
		}
	}
}

func TestUpdateLeavesForeignArchitecturesAlone(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, readmePath)
	readme := strings.Replace(fixture, "goarch: arm64", "goarch: "+runtime.GOARCH+"-other", 1)
	if err := os.WriteFile(path, []byte(readme), 0o600); err != nil {
		t.Fatal(err)
	}

	err := updateBenchmarks(t.Context(), path, 0)
	if err == nil || !strings.Contains(err.Error(), "no benchmark table was refreshed") {
		t.Fatalf("got %v, want the no-table error", err)
	}
	after, err := os.ReadFile(path) // #nosec G304 -- path is the test's own temp file
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != readme {
		t.Error("a skipped table must not be rewritten")
	}

	if err := updateBenchmarks(t.Context(), path, 1); err == nil || !strings.Contains(err.Error(), "no benchmark table was refreshed") {
		t.Errorf("-block must not override the architecture guard, got %v", err)
	}
	if err := updateBenchmarks(t.Context(), path, 2); err == nil {
		t.Error("-block past the last table was accepted")
	}
}

func TestFormatBenchmarksTouchesOnlyTables(t *testing.T) {
	path := filepath.Join(t.TempDir(), readmePath)
	sample := "Benchmark_Sample-8   1   9.000  ns/op   0  B/op   0  allocs/op"
	note := "// swar wins below 32B: 3.2 ns/op against 11.6 ns/op"
	md := strings.Replace(fixture, "sample output that is not a benchmark table", sample, 1)
	md = strings.Replace(md, "# Group", note+"\n\n# Group", 1)
	if err := os.WriteFile(path, []byte(md), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := formatBenchmarks(path); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path) // #nosec G304 -- path is the test's own temp file
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), sample) {
		t.Errorf("a fence without an Environment header must stay as it is:\n%s", after)
	}
	if !strings.Contains(string(after), note) {
		t.Errorf("prose inside a table must stay as it is:\n%s", after)
	}
	if !strings.Contains(string(after), "Benchmark_Kept/a-12  100  1.000  ns/op  0  B/op  0  allocs/op") {
		t.Errorf("table rows must be re-aligned:\n%s", after)
	}
}

func TestFormatBenchmarksRefusesABrokenFence(t *testing.T) {
	path := filepath.Join(t.TempDir(), readmePath)
	md := "```text\nBenchmark_X-8   1   1.000  ns/op   0  B/op   0  allocs/op\nprose: 5 ns/op matters\n"
	if err := os.WriteFile(path, []byte(md), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := formatBenchmarks(path); err == nil {
		t.Fatal("an unclosed fence was accepted")
	}
	after, err := os.ReadFile(path) // #nosec G304 -- path is the test's own temp file
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != md {
		t.Errorf("a refused file must not be rewritten:\n%s", after)
	}
}

func TestKnownRowsCoversRowsItCannotRefresh(t *testing.T) {
	// No ns/op column, which is what `go test` prints for a sub-nanosecond result.
	md := strings.Replace(fixture,
		"Benchmark_Gone/b-12      200    2.000  ns/op     8  B/op   1  allocs/op",
		"Benchmark_Gone/b-12      200    8  B/op   1  allocs/op", 1)
	lines := strings.Split(md, "\n")
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}
	if !knownRows(lines, blocks)["Benchmark_Gone/b"] {
		t.Error("a row this tool cannot refresh would be reported as documented nowhere")
	}
}

func TestApplyResultsKeepsThePackageTheTableNames(t *testing.T) {
	lines := strings.Split(fixture, "\n")
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}

	// `go test ./...` spans packages, and whichever prints first is not the one the
	// table documents. Taking it from the run would decide the next run's rows.
	applyResults(lines, blocks[0], 1, strings.Join([]string{
		"goos: darwin",
		"goarch: arm64",
		"pkg: github.com/gofiber/utils/v2/simd",
		"cpu: Apple M2 Pro",
		"Benchmark_Both/c-12   1   600.000 ns/op   0 B/op   0 allocs/op",
	}, "\n"), knownRows(lines, blocks))

	out := strings.Join(lines, "\n")
	if !strings.Contains(out, "pkg: github.com/gofiber/utils/v2\n") {
		t.Errorf("the table's own package must survive the run:\n%s", out)
	}
}

func TestApplyResultsKeepsRowsAMeasurementCannotFill(t *testing.T) {
	lines := strings.Split(fixture, "\n")
	blocks, err := parseBlocks(lines)
	if err != nil {
		t.Fatal(err)
	}

	// What `go test` prints without -benchmem.
	applyResults(lines, blocks[0], 1, strings.Join([]string{
		"goos: darwin",
		"goarch: arm64",
		"pkg: github.com/gofiber/utils/v2",
		"cpu: Apple M2 Pro",
		"Benchmark_Kept/a-12   1   9.000 ns/op",
	}, "\n"), knownRows(lines, blocks))

	out := strings.Join(lines, "\n")
	if !strings.Contains(out, "Benchmark_Kept/a-12      100    1.000  ns/op     0  B/op   0  allocs/op") {
		t.Errorf("a row must not lose the columns the table documents:\n%s", out)
	}
}
