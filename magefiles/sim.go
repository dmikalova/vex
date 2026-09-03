//go:build mage

package main

import (
	"os"
	"time"

	"github.com/magefile/mage/sh"
)

// Fuzz runs the whole-game fuzzer with assertions on. It drives
// internal/sim.FuzzPlay under the assert build tag, explores the game tree with
// coverage-guided mutation, and saves any minimized failing script under
// internal/sim/testdata/fuzz/FuzzPlay.
//
// The budget defaults to 60s. Pass a Go duration to change it, either as a mage
// flag or through FUZZTIME, which also takes an "Nx" execution count that a
// duration cannot express:
//
//	mage fuzz
//	mage fuzz -duration=5m
//	mage fuzz -d=5m
//	FUZZTIME=10000x mage fuzz
//
// The generated corpus lives in the Go build cache, not the repo — run
// `mage fuzzClean` to reset it.
func Fuzz(duration, d *time.Duration) error {
	if budget := firstSet(duration, d); budget != "" {
		os.Setenv("FUZZTIME", budget)
	}
	fuzztime := os.Getenv("FUZZTIME")
	if fuzztime == "" {
		fuzztime = "60s"
	}
	return sh.RunV("go", "test", "-tags", "assert",
		"-run", "^$", "-fuzz", "^FuzzPlay$", "-fuzztime", fuzztime,
		"./internal/sim")
}

// firstSet returns the first duration flag the caller actually passed, so a
// target can offer both a spelled-out name and a one-letter alias. mage derives
// each flag from a parameter name and has no alias mechanism, so the alias has to
// be a second parameter.
func firstSet(flags ...*time.Duration) string {
	for _, f := range flags {
		if f != nil {
			return f.String()
		}
	}
	return ""
}

// FuzzClean discards the generated fuzz corpus. It lives in the Go build cache.
// Reach for this when a stale corpus keeps steering the fuzzer down the same
// paths. The checked-in seed corpus under internal/sim/testdata is untouched.
func FuzzClean() error {
	return sh.RunV("go", "clean", "-fuzzcache")
}

// CorpusPrune prunes FuzzPlay's seed corpus of fixed bugs. Each bug that still
// reproduces is left as one minimized entry. The soak saves every failing
// script verbatim, so a single bug can leave hundreds of near-identical
// multi-kilobyte entries; this replays them all, drops the ones whose bug is
// fixed, and shrinks what remains.
func CorpusPrune() error {
	return sh.RunV("go", "run", "./magefiles/simcorpus")
}

// Debug replays a failing simulated game with the game log on. It prints the log
// tail next to the invariant violation that ended the game, turning a soak, fuzz,
// or property-test find into a readable sequence of plays.
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

// Trace writes a full game log to a file, end to end. It plays the fixed-seed
// property games once with the log turned on, so a whole game can be read at once
// instead of one prompt at a time. Unlike `mage debug`, which shows the tail of a
// game that broke, a trace is the full log of games that pass.
//
// COUNT sets how many of the property batch's games to play (default 1) and OUT
// the destination; it defaults to tmp/sim/trace.log, under the repo's gitignored
// scratch directory.
//
//	mage trace
//	COUNT=25 mage trace
//	OUT=tmp/sim/mine.log mage trace
func Trace() error {
	return sh.RunV("go", "run", "./magefiles/simtrace",
		os.Getenv("COUNT"), os.Getenv("OUT"))
}

// Soak runs the long-running game soak with assertions on. It
// drives internal/sim.TestSoak under the assert build tag, churning fresh random
// games across GOMAXPROCS workers until the time budget elapses. It does not stop
// at the first failure: every failing script is saved into
// internal/sim/testdata/fuzz/FuzzPlay as a permanent FuzzPlay regression.
//
// The budget defaults to 30s. Pass a Go duration to change it, either as a mage
// flag or through SOAK_DURATION:
//
//	mage soak
//	mage soak -duration=5m
//	mage soak -d=5m
//	SOAK_DURATION=5m mage soak
func Soak(duration, d *time.Duration) error {
	if budget := firstSet(duration, d); budget != "" {
		os.Setenv("SOAK_DURATION", budget)
	}
	if os.Getenv("SOAK_DURATION") == "" {
		os.Setenv("SOAK_DURATION", "30s")
	}
	return sh.RunV("go", "test", "-tags", "assert",
		"-timeout", soakTimeout(os.Getenv("SOAK_DURATION")),
		"-run", "^TestSoak$", "-count", "1", "-v", "./internal/sim")
}

// soakTimeout gives go test a panic timeout past the soak's own budget. The soak
// stops itself at the budget, so the timeout is only there to catch a genuine
// hang — but go test defaults it to 10m, which kills a soak of 10m or longer
// mid-run and reports the budget elapsing as a timeout panic. An unparsable
// budget gets no timeout at all; the soak itself fails on it with a better
// message than a goroutine dump.
func soakTimeout(budget string) string {
	d, err := time.ParseDuration(budget)
	if err != nil {
		return "0"
	}
	return (d + 5*time.Minute).String()
}
