package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if err := formatBenchmarks("README.md"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func formatBenchmarks(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	type meta struct {
		kind  string
		index int
	}

	lineMeta := make([]meta, len(lines))
	benchEntries := [][]string{}
	memEntries := [][]string{}

	inBlock := false
	re := regexp.MustCompile(`\S+`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "```go":
			inBlock = true
		case trimmed == "```" && inBlock:
			inBlock = false
		default:
			if inBlock {
				expanded := strings.ReplaceAll(line, "\t", "    ")
				if strings.Contains(expanded, "ns/op") {
					toks := re.FindAllString(expanded, -1)
					benchEntries = append(benchEntries, toks)
					lineMeta[i] = meta{kind: "bench", index: len(benchEntries) - 1}
				} else if strings.Contains(expanded, "allocs/op") {
					toks := re.FindAllString(expanded, -1)
					memEntries = append(memEntries, toks)
					lineMeta[i] = meta{kind: "mem", index: len(memEntries) - 1}
				}
			}
		}
	}

	benchWidths := columnWidths(benchEntries)
	memWidths := columnWidths(memEntries)

	var output []string
	inBlock = false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "```go":
			inBlock = true
			output = append(output, line)
			continue
		case trimmed == "```" && inBlock:
			inBlock = false
			output = append(output, line)
			continue
		}

		m := lineMeta[i]
		if m.kind == "bench" {
			toks := benchEntries[m.index]
			formatted := []string{}
			for j, t := range toks {
				switch j {
				case 0:
					formatted = append(formatted, padRight(t, benchWidths[j]))
				case 1, 2, 4, 6:
					formatted = append(formatted, padLeft(t, benchWidths[j]))
				default:
					formatted = append(formatted, t)
				}
			}
			output = append(output, strings.Join(formatted, "  "))
		} else if m.kind == "mem" {
			toks := memEntries[m.index]
			formatted := []string{}
			for j, t := range toks {
				switch j {
				case 0, 2:
					formatted = append(formatted, padLeft(t, memWidths[j]))
				default:
					formatted = append(formatted, t)
				}
			}
			output = append(output, strings.Join(formatted, "  "))
		} else {
			output = append(output, line)
		}
	}

	out := strings.Join(output, "\n")
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0o644)
}

func columnWidths(entries [][]string) []int {
	maxCols := 0
	for _, e := range entries {
		if len(e) > maxCols {
			maxCols = len(e)
		}
	}

	widths := make([]int, maxCols)
	for _, e := range entries {
		for i, t := range e {
			if len(t) > widths[i] {
				widths[i] = len(t)
			}
		}
	}
	return widths
}

func padLeft(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat(" ", n-len(s)) + s
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
