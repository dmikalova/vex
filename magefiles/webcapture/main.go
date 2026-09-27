// Command webcapture prunes the browser client's replay captures. It replays
// every entry the dev server has written and deletes the ones that are no longer
// open findings: the faults that have been fixed, and the entries recorded against
// a different command-log version or card pool. What is left is the list of
// findings still to fix.
//
// Run it via `mage capturePrune`. Sibling of `mage corpusPrune`.
package main

import (
	"fmt"
	"os"

	"github.com/dmikalova/vex/internal/web"
)

func main() {
	report, err := web.PruneCaptures(web.CaptureDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf(
		"scanned %d captures: kept %d open, dropped %d fixed and %d stale\n",
		report.Scanned, len(report.Open), report.Fixed, report.Stale)
	for _, f := range report.Open {
		fmt.Printf("  still failing: %s\n", f)
	}
}
