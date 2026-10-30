package main

import (
	"context"
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
		"Benchmark_Wide/d-4   	     500	         5.500 ns/op	  181.82 MB/s	       0 B/op	       0 allocs/op",
		"Benchmark_New/e-4    	     400	         4.500 ns/op	       0 B/op	       0 allocs/op",
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

	// arm64 Linux prints no cpu line.
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
	// An even -count has no middle: the slower of the two is promised.
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
		"a\nb":   "a\nb\n",
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

// The tool parses tables by shape, so README drift must fail here.
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
	// No ns/op column, as for a sub-nanosecond result.
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

	// go test ./... spans packages; the first to print need not be the table's.
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

func TestPadding(t *testing.T) {
	for _, tc := range []struct {
		name, got, want string
	}{
		{"left pads", padLeft("7", 3), "  7"},
		{"left keeps a wide value", padLeft("777", 2), "777"},
		{"left keeps an exact fit", padLeft("77", 2), "77"},
		{"right pads", padRight("a", 3), "a  "},
		{"right keeps a wide value", padRight("abc", 2), "abc"},
		{"right keeps an exact fit", padRight("ab", 2), "ab"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, tc.got, tc.want)
		}
	}

	row := padRow([]string{"Benchmark_A-4", "7", "1.5", "ns/op", "0", "B/op", "0", "allocs/op"}, []int{16, 4, 6, 5, 3, 4, 3, 9})
	if want := "Benchmark_A-4     " + "   7" + "     1.5" + "  ns/op" + "    0" + "  B/op" + "  " + "  0" + "  allocs/op"; row != want {
		t.Errorf("padRow = %q, want %q", row, want)
	}
}

func TestReadingAMissingReadmeIsAnError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.md")
	if _, err := readLines(missing); err == nil {
		t.Error("readLines of a missing file succeeded")
	}
	if err := formatBenchmarks(missing); err == nil {
		t.Error("formatBenchmarks of a missing file succeeded")
	}
	if err := updateBenchmarks(context.Background(), missing, 0); err == nil {
		t.Error("updateBenchmarks of a missing file succeeded")
	}
}

func TestUpdateBenchmarksRefusesABrokenFence(t *testing.T) {
	broken := writeReadme(t, t.TempDir(), "```text\nrow\n")
	if err := updateBenchmarks(context.Background(), broken, 0); err == nil {
		t.Error("a README with an unclosed fence was accepted")
	}
}

func TestUpdateBenchmarksSelectsOneTable(t *testing.T) {
	dir := writeMiniModule(t)
	two := "# Two tables\n\n" +
		miniTable("// this one has no command to run") + "\n" +
		miniTable(miniCommand)
	readme := writeReadme(t, dir, two)

	var err error
	logged := captureStderr(t, func() { err = updateBenchmarks(context.Background(), readme, 1) })
	if err == nil || !strings.Contains(err.Error(), "no benchmark table was refreshed") {
		t.Fatalf("got %v, want the refusal to write an unmeasured README", err)
	}
	if !strings.Contains(logged, "table 1: skipped, no") {
		t.Errorf("the reason was not logged:\n%s", logged)
	}
	if readFile(t, readme) != two {
		t.Error("nothing was measured, so nothing may change")
	}

	logged = captureStderr(t, func() { err = updateBenchmarks(context.Background(), readme, 2) })
	if err != nil {
		t.Fatalf("%v\nlog:\n%s", err, logged)
	}
	out := readFile(t, readme)
	first, second, _ := strings.Cut(out, miniCommand)
	if strings.Count(first, "777") != 1 || strings.Contains(second, "777") {
		t.Errorf("want table 1's placeholder row kept and table 2's replaced:\n%s", out)
	}
}

func TestApplyResultsReportsWhatItCannotRefresh(t *testing.T) {
	lines := strings.Split(strings.Join([]string{
		"Environment:", "goos: darwin", "goarch: arm64", "pkg: m", "cpu: c", "",
		"```text",
		"// go test . -bench=.",
		"Benchmark_Odd/x-12   not a row",
		"Benchmark_Both/c-12  300  3.000  ns/op  0  B/op  0  allocs/op",
		"```",
	}, "\n"), "\n")
	blocks, err := parseBlocks(lines)
	if err != nil || len(blocks) != 1 {
		t.Fatalf("got %d tables, %v", len(blocks), err)
	}

	logged := captureStderr(t, func() {
		applyResults(lines, blocks[0], 1, strings.Join([]string{
			"goos: linux", "goarch: arm64", "cpu: c2",
			"pkg: m/one", "Benchmark_Both/c-4  1  1.000 ns/op  0 B/op  0 allocs/op",
			"pkg: m/two", "Benchmark_Both/c-4  1  2.000 ns/op  0 B/op  0 allocs/op",
		}, "\n"), knownRows(lines, blocks))
	})
	for _, want := range []string{
		"not in a shape this tool can refresh (1): Benchmark_Odd/x-12",
		"measured by several packages and none of them is the table's, left as they were (1): Benchmark_Both/c",
		"0 rows refreshed",
	} {
		if !strings.Contains(logged, want) {
			t.Errorf("the log lacks %q:\n%s", want, logged)
		}
	}
	if out := strings.Join(lines, "\n"); !strings.Contains(out, "Benchmark_Both/c-12  300  3.000  ns/op") {
		t.Errorf("an ambiguous row must keep its numbers:\n%s", out)
	}
}

func TestApplyResultsTruncatesLongReports(t *testing.T) {
	lines := strings.Split(strings.Join([]string{
		"Environment:", "goos: darwin", "goarch: arm64", "pkg: m", "cpu: c", "",
		"```text", "// go test . -bench=.", "```",
	}, "\n"), "\n")
	blocks, err := parseBlocks(lines)
	if err != nil || len(blocks) != 1 {
		t.Fatalf("got %d tables, %v", len(blocks), err)
	}
	var out []string
	for i := range 12 {
		out = append(out, "Benchmark_New"+string(rune('A'+i))+"-4  1  1.000 ns/op  0 B/op  0 allocs/op")
	}
	logged := captureStderr(t, func() {
		applyResults(lines, blocks[0], 1, strings.Join(out, "\n"), knownRows(lines, blocks))
	})
	if !strings.Contains(logged, "measured but in none of the tables (12): Benchmark_NewA") ||
		!strings.Contains(logged, "... and 2 more") || strings.Contains(logged, "Benchmark_NewL") {
		t.Errorf("a report of twelve names must show ten and count the rest:\n%s", logged)
	}
}

func TestApplyResultsSkipsHeaderLinesTheBlockLacks(t *testing.T) {
	// applyResults must not trust that parseBlocks always finds the Environment lines.
	lines := []string{"```text", "// go test .", "```"}
	block := benchBlock{open: 0, end: 2, env: map[string]int{}}
	captureStderr(t, func() { applyResults(lines, block, 1, "", map[string]bool{}) })
	if strings.Join(lines, "\n") != "```text\n// go test .\n```" {
		t.Errorf("lines changed: %q", lines)
	}
}

func TestMedianOfEqualRuns(t *testing.T) {
	got := median([]benchResult{{ns: 2, fields: []string{"a"}}, {ns: 2, fields: []string{"b"}}, {ns: 1, fields: []string{"c"}}})
	if got == nil || got.ns != 2 {
		t.Errorf("got %v, want a run of 2ns", got)
	}
}

func TestParseBenchOutputReportsAnUnreadableTime(t *testing.T) {
	_, results, odd := parseBenchOutput("Benchmark_X-4   1   abc ns/op   0 B/op   0 allocs/op\nBenchmark_Y-4   1   2.5 ns/op   0 B/op   0 allocs/op")
	if len(odd) != 1 || odd[0] != "Benchmark_X-4" {
		t.Errorf("got %v, want the unreadable row reported", odd)
	}
	if len(results["Benchmark_X"]) != 0 || len(results["Benchmark_Y"]) != 1 {
		t.Errorf("got %v", results)
	}
}

func TestParseBlocksRejectsAnUnbalancedCommand(t *testing.T) {
	md := "Environment:\ngoos: linux\ngoarch: amd64\npkg: m\ncpu: c\n\n```text\n// go test -bench='^X\nrow\n```\n"
	_, err := parseBlocks(strings.Split(md, "\n"))
	if err == nil || !strings.Contains(err.Error(), "line 8") {
		t.Errorf("got %v, want an error at the command line", err)
	}
}

func TestReadEnvHeaderStopsAtAForeignLine(t *testing.T) {
	md := "Environment:\ngoos: linux\ngoarch: amd64\npkg: m\ncpu: c\nNote: measured on a laptop\n```text\n// go test .\n```\n"
	blocks, err := parseBlocks(strings.Split(md, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 0 {
		t.Errorf("got %d tables, want none", len(blocks))
	}
}
