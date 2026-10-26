package main

import (
	"strings"
)

// formatBenchmarks uses the updater's parser, so prose and sample output stay untouched.
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
	at := map[int]int{}
	for _, b := range blocks {
		for i := b.open + 1; i < b.end; i++ {
			if fields, ok := benchFields(lines[i]); ok {
				at[i] = len(rows)
				rows = append(rows, fields)
			}
		}
	}

	// One grid for every table; per-table widths would re-pad existing rows once.
	widths := columnWidths(rows)
	for i, r := range at {
		lines[i] = padRow(rows[r], widths)
	}
	return writeLines(path, lines)
}

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
