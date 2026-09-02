//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

// Fuzz runs the whole-game fuzzer (internal/sim.FuzzPlay) with the engine's
// in-game invariant checks turned on via the assert build tag. It explores the
// game tree with coverage-guided mutation and saves any minimized failing script
// under internal/sim/testdata/fuzz/FuzzPlay.
//
// Set FUZZTIME to bound the run (a Go duration or an "Nx" count); it defaults to
// 60s. The generated corpus lives in the Go build cache, not the repo — run
// `mage fuzzClean` to reset it.
func Fuzz() error {
	fuzztime := os.Getenv("FUZZTIME")
	if fuzztime == "" {
		fuzztime = "60s"
	}
	return sh.RunV("go", "test", "-tags", "assert",
		"-run", "^$", "-fuzz", "^FuzzPlay$", "-fuzztime", fuzztime,
		"./internal/sim")
}

// FuzzClean discards the generated fuzz corpus in the Go build cache, for when a
// stale corpus keeps steering the fuzzer down the same paths. The checked-in seed
// corpus under internal/sim/testdata is untouched.
func FuzzClean() error {
	return sh.RunV("go", "clean", "-fuzzcache")
}

// CorpusPrune rewrites FuzzPlay's seed corpus as one minimized entry per bug that
// still reproduces. The soak saves every failing script verbatim, so a single bug
// can leave hundreds of near-identical multi-kilobyte entries; this replays them
// all, drops the ones whose bug is fixed, and shrinks what remains.
func CorpusPrune() error {
	return sh.RunV("go", "run", "./magefiles/simcorpus")
}

// Debug replays one simulated game with the game log on and prints the log tail
// next to the invariant violation that ended it, turning a soak, fuzz, or
// property-test find into a readable sequence of plays.
//
// With no SCRIPT it searches the fixed-seed property batch that `mage test` plays
// and replays the first failure; set SCRIPT to the hex script a failure printed to
// replay that one instead. TAIL widens the log tail (default 60 lines).
//
//	mage debug
//	SCRIPT=00ff1a... mage debug
//	TAIL=200 mage debug
func Debug() error {
	return sh.RunV("go", "run", "./magefiles/simdebug",
		os.Getenv("SCRIPT"), os.Getenv("TAIL"))
}

// Soak runs the long-running game soak (internal/sim.TestSoak) with the assert
// build tag, churning fresh random games across GOMAXPROCS workers until the time
// budget elapses. It does not stop at the first failure: every failing script is
// saved into internal/sim/testdata/fuzz/FuzzPlay as a permanent FuzzPlay regression.
// Set SOAK_DURATION to a Go duration (e.g. SOAK_DURATION=5m); it defaults to 30s.
func Soak() error {
	if os.Getenv("SOAK_DURATION") == "" {
		os.Setenv("SOAK_DURATION", "30s")
	}
	return sh.RunV("go", "test", "-tags", "assert",
		"-run", "^TestSoak$", "-count", "1", "-v", "./internal/sim")
}
