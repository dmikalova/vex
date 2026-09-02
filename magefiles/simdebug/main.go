// Command simdebug replays one simulated game with the game log turned on and
// prints the log tail alongside the invariant violation that ended it. It exists
// so a soak, fuzz, or property-test failure can be read as the sequence of plays
// that produced it instead of a bare invariant message.
//
// Run it via `mage debug`, which passes SCRIPT and TAIL through as its arguments:
//
//	mage debug                     # first failing fixed-seed property game
//	SCRIPT=00ff1a... mage debug    # a specific script (the hex a failure prints)
//	TAIL=200 mage debug            # widen the log tail from its default
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/dmikalova/vactrol/internal/sim"
)

// seedBatch matches the batch size TestSimulateSeeds plays, so a bare `mage debug`
// searches exactly the games `mage test` runs.
const seedBatch = 1000

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var scriptHex string
	if len(args) > 0 {
		scriptHex = args[0]
	}
	tail := 60
	if len(args) > 1 && args[1] != "" {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("tail %q is not a number: %w", args[1], err)
		}
		tail = n
	}

	var f sim.Failure
	if scriptHex == "" {
		fmt.Printf("searching the %d fixed-seed property games…\n", seedBatch)
		f = sim.FirstSeedFailure(seedBatch)
	} else {
		script, err := hex.DecodeString(scriptHex)
		if err != nil {
			return fmt.Errorf("script %q is not hex: %w", scriptHex, err)
		}
		f = sim.Debug(script)
	}

	if f.Err == nil {
		fmt.Println("no failure: every game played out cleanly")
		return nil
	}

	fmt.Printf("script: %s\n\n", hex.EncodeToString(f.Script))
	lines := f.Log
	if tail > 0 && len(lines) > tail {
		fmt.Printf("… %d earlier log lines elided\n", len(lines)-tail)
		lines = lines[len(lines)-tail:]
	}
	for _, l := range lines {
		fmt.Println(l)
	}
	fmt.Printf("\nVIOLATION: %v\n", f.Err)
	return nil
}
