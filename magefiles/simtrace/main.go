// Command simtrace plays the fixed-seed property games once with the game log
// turned on and writes every log line to a file, so a whole game can be read as
// the sequence of plays it was rather than watched one prompt at a time.
//
// Run it via `mage trace`, which passes COUNT and OUT through as its arguments:
//
//	mage trace              # the first property game, into tmp/sim/trace.log
//	COUNT=25 mage trace     # the first 25 games
//	OUT=tmp/sim/x.log mage trace
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dmikalova/vactrol/internal/sim"
)

// defaultOut keeps traces in ./tmp, the repo's gitignored scratch directory, so a
// log dump is never a candidate for a commit.
const defaultOut = "tmp/sim/trace.log"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	count := 1
	if len(args) > 0 && args[0] != "" {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("count %q is not a number: %w", args[0], err)
		}
		if n < 1 {
			return fmt.Errorf("count %d must be at least 1", n)
		}
		count = n
	}
	out := defaultOut
	if len(args) > 1 && args[1] != "" {
		out = args[1]
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	if err := sim.Trace(count, w); err != nil {
		_ = f.Close()
		return err
	}
	if err := w.Flush(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	fmt.Printf("wrote the log of %d property game(s) to %s\n", count, out)
	return nil
}
