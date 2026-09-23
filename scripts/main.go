package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

const readmePath = "README.md"

const usage = `usage: go run ./scripts [-update] [-block N]

  (no flags)  re-align the benchmark tables in README.md
  -update     re-run each table's own "// go test ..." command and refresh its numbers
  -block N    restrict -update to the Nth table (1-based); it still has to be a
              table this machine's architecture can measure
`

var errHelp = errors.New("help requested")

func main() {
	if err := run(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err) //nolint:errcheck // printing error only
		os.Exit(1)
	}
}

func run(args []string) error {
	update, block, err := parseFlags(args)
	switch {
	case errors.Is(err, errHelp):
		logf("%s", usage)
		return nil
	case err != nil:
		return fmt.Errorf("%w\n\n%s", err, usage)
	}

	if update {
		// A Ctrl-C during a 15 minute measurement should take the go test child with it.
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		if err := updateBenchmarks(ctx, readmePath, block); err != nil {
			return err
		}
	}

	return formatBenchmarks(readmePath)
}

// Parsed by hand: the repo's depguard rules deny the flag package.
func parseFlags(args []string) (bool, int, error) {
	update, block := false, 0
	for i := 0; i < len(args); i++ {
		name, value, hasValue := strings.Cut(args[i], "=")
		if hasValue && name != "-block" && name != "--block" {
			return false, 0, fmt.Errorf("%s takes no value", name)
		}

		switch name {
		case "-h", "--help":
			return false, 0, errHelp
		case "-update", "--update":
			update = true
		case "-block", "--block":
			if !hasValue {
				rest := args[i+1:]
				if len(rest) == 0 {
					return false, 0, errors.New("-block needs a table number")
				}
				value = rest[0]
				i++
			}
			n, err := strconv.Atoi(value)
			if err != nil || n < 1 {
				return false, 0, fmt.Errorf("-block wants a positive table number, got %q", value)
			}
			block = n
		default:
			return false, 0, fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if block > 0 && !update {
		return false, 0, errors.New("-block only applies together with -update")
	}
	return update, block, nil
}
