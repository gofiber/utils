package main

import (
	"strings"
)

// formatBenchmarks finds the rows through the same parser the updater uses, so
// prose, sample output and anything that is not a table stay exactly as they are.
func formatBenchmarks(path string) error {
	lines, err := readLines(path)
	if err != nil {
		return err
	}
	blocks, err := parseBlocks(lines)
	if err != nil {
		return err
	}

	var rows [][]string
	at := map[int]int{} // line index -> index in rows
	for _, b := range blocks {
		for i := b.open + 1; i < b.end; i++ {
			if fields, ok := benchFields(lines[i]); ok {
				at[i] = len(rows)
				rows = append(rows, fields)
			}
		}
	}

	// One grid for every table. Per table would re-pad every existing row once and
	// only narrows which table a refresh reflows, which is not worth that diff.
	widths := columnWidths(rows)
	for i, r := range at {
		lines[i] = padRow(rows[r], widths)
	}
	return writeLines(path, lines)
}

// padRow puts the name against the left edge and each measured value against the
// right edge of its column; the units follow whatever they belong to.
func padRow(fields []string, widths []int) string {
	out := make([]string, len(fields))
	for i, f := range fields {
		switch i {
		case 0:
			out[i] = padRight(f, widths[i])
		case 1, 2, 4, 6:
			out[i] = padLeft(f, widths[i])
		default:
			out[i] = f
		}
	}
	return strings.Join(out, "  ")
}

func columnWidths(rows [][]string) []int {
	maxCols := 0
	for _, r := range rows {
		if len(r) > maxCols {
			maxCols = len(r)
		}
	}

	widths := make([]int, maxCols)
	for _, r := range rows {
		for i, f := range r {
			if len(f) > widths[i] {
				widths[i] = len(f)
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
