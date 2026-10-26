package main

// Refreshes the numbers in the README benchmark tables: each ```text table starts with the
// `go test` command that produced it, so re-run that, replace the values of the rows it
// reports, and rewrite the Environment: header above the table with the machine that ran.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

var envKeys = []string{"goos", "goarch", "pkg", "cpu"}

// machineKeys omits pkg: pick() needs the table's own package, not the run's.
var machineKeys = []string{"goos", "goarch", "cpu"}

type benchBlock struct {
	env    map[string]int
	pkg    string
	goos   string
	goarch string
	cmd    []string
	open   int
	end    int
}

type benchResult struct {
	pkg    string
	fields []string
	ns     float64
}

func updateBenchmarks(ctx context.Context, path string, only int) error {
	lines, err := readLines(path)
	if err != nil {
		return err
	}

	blocks, err := parseBlocks(lines)
	if err != nil {
		return err
	}
	if only > len(blocks) {
		return fmt.Errorf("-block %d: %s has %d benchmark tables", only, path, len(blocks))
	}

	known := knownRows(lines, blocks)
	dir := filepath.Dir(path)
	ran := 0
	for i, b := range blocks {
		num := i + 1
		switch {
		case only > 0 && only != num:
			continue
		// Skip other architectures: unmeasured rows would keep that machine's numbers.
		case b.goarch != runtime.GOARCH:
			logf("table %d (%s/%s): skipped, this machine is %s/%s\n", num, b.goos, b.goarch, runtime.GOOS, runtime.GOARCH)
			continue
		case len(b.cmd) == 0:
			logf("table %d: skipped, no \"// go test ...\" line to re-run\n", num)
			continue
		}

		logf("table %d (%s/%s): %s\n", num, b.goos, b.goarch, strings.Join(b.cmd, " "))
		out, err := runBench(ctx, dir, b.cmd)
		if err != nil {
			return err
		}
		applyResults(lines, b, num, out, known)
		ran++
	}

	if ran == 0 {
		return fmt.Errorf("no benchmark table was refreshed on %s/%s, see the reasons above", runtime.GOOS, runtime.GOARCH)
	}
	return writeLines(path, lines)
}

// knownRows names every documented benchmark; one table's run can measure another's rows.
func knownRows(lines []string, blocks []benchBlock) map[string]bool {
	known := map[string]bool{}
	for _, b := range blocks {
		for i := b.open + 1; i < b.end; i++ {
			if name, ok := rowName(lines[i]); ok {
				known[stripProcs(name)] = true
			}
		}
	}
	return known
}

func rowName(line string) (string, bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 || !strings.HasPrefix(fields[0], "Benchmark") {
		return "", false
	}
	return fields[0], true
}

func applyResults(lines []string, b benchBlock, num int, out string, known map[string]bool) {
	env, results, odd := parseBenchOutput(out)

	updated := 0
	var stale, mixed, unusable, narrower []string
	for i := b.open + 1; i < b.end; i++ {
		fields, ok := benchFields(lines[i])
		if !ok {
			if name, isRow := rowName(lines[i]); isRow {
				unusable = append(unusable, name)
			}
			continue
		}
		name := stripProcs(fields[0])

		runs, ok := pick(results[name], b.pkg)
		if !ok {
			mixed = append(mixed, name)
			continue
		}
		best := median(runs)
		if best == nil {
			stale = append(stale, name)
			continue
		}
		row := canonical(best.fields)
		if len(row) < len(fields) {
			// This run measured fewer columns; replacing the row would drop the rest.
			narrower = append(narrower, name)
			continue
		}
		lines[i] = strings.Join(row, "  ")
		updated++
	}

	var missing []string
	for _, key := range machineKeys {
		at, ok := b.env[key]
		if !ok {
			continue
		}
		if env[key] == "" {
			// Better "unknown" than a stale machine; arm64 Linux reports no cpu line.
			lines[at] = key + ": unknown"
			missing = append(missing, key)
			continue
		}
		lines[at] = key + ": " + env[key]
	}

	var added []string
	for name := range results {
		if !known[name] {
			added = append(added, name)
		}
	}

	logf("table %d: %d rows refreshed\n", num, updated)
	report(num, "not in the new output, left as they were", stale)
	report(num, "measured but in none of the tables", added)
	report(num, "in the table, but not in a shape this tool can refresh", unusable)
	report(num, "measured with fewer columns than the table has, left as they were", narrower)
	report(num, "measured by several packages and none of them is the table's, left as they were", mixed)
	report(num, "measured in a shape the table cannot hold", odd)
	report(num, "not reported by this run, the header says unknown", missing)
}

// A name can come from two packages; the table's own package decides.
func pick(runs []benchResult, pkg string) ([]benchResult, bool) {
	packages := map[string]bool{}
	for _, r := range runs {
		packages[r.pkg] = true
	}
	if len(packages) < 2 {
		return runs, true
	}
	owned := slices.DeleteFunc(slices.Clone(runs), func(r benchResult) bool { return r.pkg != pkg })
	return owned, len(owned) > 0
}

// canonical drops extra columns such as MB/s, which would misalign the table.
func canonical(fields []string) []string {
	row := append([]string{}, fields[:2]...)
	for i := 2; i+1 < len(fields); i += 2 {
		switch fields[i+1] {
		case "ns/op", "B/op", "allocs/op":
			row = append(row, fields[i], fields[i+1])
		}
	}
	return row
}

// ns/op must be the first metric: tables are aligned by column index.
func benchFields(line string) ([]string, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 || !strings.HasPrefix(fields[0], "Benchmark") || fields[3] != "ns/op" {
		return nil, false
	}
	return fields, true
}

func report(num int, title string, names []string) {
	if len(names) == 0 {
		return
	}
	slices.Sort(names)
	shown := names
	if len(shown) > 10 {
		shown = shown[:10]
	}
	logf("  table %d: %s (%d): %s\n", num, title, len(names), strings.Join(shown, ", "))
	if len(names) > len(shown) {
		logf("  ... and %d more\n", len(names)-len(shown))
	}
}

func median(runs []benchResult) *benchResult {
	if len(runs) == 0 {
		return nil
	}
	sorted := slices.Clone(runs)
	slices.SortFunc(sorted, func(a, b benchResult) int {
		switch {
		case a.ns < b.ns:
			return -1
		case a.ns > b.ns:
			return 1
		default:
			return 0
		}
	})
	return &sorted[len(sorted)/2]
}

func parseBenchOutput(out string) (map[string]string, map[string][]benchResult, []string) {
	env := map[string]string{}
	results := map[string][]benchResult{}
	var odd []string
	pkg := ""

	for line := range strings.Lines(out) {
		line = strings.TrimSpace(line)
		if key, value, ok := strings.Cut(line, ": "); ok && slices.Contains(envKeys, key) {
			if key == "pkg" {
				pkg = strings.TrimSpace(value)
			}
			if env[key] == "" {
				env[key] = strings.TrimSpace(value)
			}
			continue
		}
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}

		fields, ok := benchFields(line)
		if !ok {
			odd = append(odd, strings.Fields(line)[0])
			continue
		}
		ns, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			odd = append(odd, fields[0])
			continue
		}
		name := stripProcs(fields[0])
		results[name] = append(results[name], benchResult{pkg: pkg, fields: fields, ns: ns})
	}
	return env, results, odd
}

func stripProcs(name string) string {
	if i := strings.LastIndexByte(name, '-'); i > 0 {
		if _, err := strconv.Atoi(name[i+1:]); err == nil {
			return name[:i]
		}
	}
	return name
}

func parseBlocks(lines []string) ([]benchBlock, error) {
	var blocks []benchBlock
	for i, line := range lines {
		if strings.TrimSpace(line) != "```text" {
			continue
		}

		b := benchBlock{open: i, env: map[string]int{}}
		for j := i + 1; j < len(lines); j++ {
			text := strings.TrimSpace(lines[j])
			if text == "```text" {
				return nil, fmt.Errorf("line %d: ```text block opened at line %d is never closed", j+1, i+1)
			}
			if text == "```" {
				b.end = j
				break
			}
			if b.cmd == nil && strings.HasPrefix(text, "// go test") {
				cmd, err := splitArgs(strings.TrimPrefix(text, "//"))
				if err != nil {
					return nil, fmt.Errorf("line %d: %w", j+1, err)
				}
				b.cmd = cmd
			}
		}
		if b.end == 0 {
			return nil, fmt.Errorf("line %d: ```text block is never closed", i+1)
		}

		readEnvHeader(lines, &b)
		// A block without a full Environment: header is sample output, not a table.
		if len(b.env) == len(envKeys) {
			blocks = append(blocks, b)
		}
	}
	return blocks, nil
}

func readEnvHeader(lines []string, b *benchBlock) {
	for j := b.open - 1; j >= 0; j-- {
		text := strings.TrimSpace(lines[j])
		switch text {
		case "":
			continue
		case "Environment:":
			return
		}

		key, value, ok := strings.Cut(text, ": ")
		if !ok {
			return
		}
		value = strings.TrimSpace(value)
		switch key {
		case "goos":
			b.goos = value
		case "goarch":
			b.goarch = value
		case "pkg":
			b.pkg = value
		case "cpu":
		default:
			return
		}
		b.env[key] = j
	}
}

func splitArgs(s string) ([]string, error) {
	var args []string
	var cur strings.Builder
	quote := byte(0)
	started := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0 && c == quote:
			quote = 0
		case quote == 0 && (c == '\'' || c == '"'):
			quote = c
			started = true
		case quote == 0 && (c == ' ' || c == '\t'):
			if started {
				args = append(args, cur.String())
				cur.Reset()
				started = false
			}
		default:
			cur.WriteByte(c)
			started = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unbalanced %c in %q", quote, strings.TrimSpace(s))
	}
	if started {
		args = append(args, cur.String())
	}

	if len(args) < 2 || args[0] != "go" || args[1] != "test" {
		return nil, fmt.Errorf("expected a `go test` command, got %q", strings.TrimSpace(s))
	}
	return args, nil
}

func runBench(ctx context.Context, dir string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec G204 -- `go test` runs this repo's own code either way; no shell means the line cannot mean more than arguments
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = io.MultiWriter(&out, os.Stderr)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", strings.Join(args, " "), err)
	}
	return out.String(), nil
}

func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n"), nil
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600) // #nosec G306
}

func logf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...) //nolint:errcheck // progress reporting only
}
