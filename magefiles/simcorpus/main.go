// Command simcorpus rewrites FuzzPlay's seed corpus as one minimized entry per
// bug that still reproduces. A soak saves every failing script verbatim, so one
// bug can leave hundreds of near-identical multi-kilobyte entries behind; this
// replays them all, drops the ones whose bug is fixed, and keeps the rest keyed
// by the violation they trip.
//
// Run it via `mage corpusPrune`.
package main

import (
	"fmt"
	"os"

	"github.com/dmikalova/vactrol/internal/sim"
)

// corpusDir is FuzzPlay's seed corpus from the repository root, where mage runs.
const corpusDir = "internal/sim/" + sim.CorpusDir

func main() {
	report, err := sim.PruneCorpus(corpusDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("scanned %d entries: kept %d, dropped %d whose bug is fixed\n",
		report.Scanned, report.Kept, report.Fixed)
	for _, f := range report.Failures {
		fmt.Printf("  still failing: %s\n", f)
	}
}
