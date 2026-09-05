//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

// A coverGate is one area held at 100%. test names the packages whose tests run
// and count names the packages whose statements are counted. They differ where a
// package is exercised through another's tests: the card definitions are covered
// by the per-card tests that sit beside them, so the gate counts the sets while
// running everything under internal/cards.
type coverGate struct {
	name  string
	test  string
	count string
}

// coverGates lists the areas the gate holds at 100%. Absent from it, and so
// deliberately ungated: internal/web, which is a view layer the tests exercise
// through only a few entry points.
var coverGates = []coverGate{
	{name: "engine", test: "./internal/engine/", count: "./internal/engine/"},
	{name: "cards", test: "./internal/cards/...", count: "./internal/cards/sets/..."},
	{name: "cardtest", test: "./internal/cards/...", count: "./internal/cards/cardtest/"},
	{name: "deckgen", test: "./internal/deckgen/", count: "./internal/deckgen/"},
}

// Cover runs the gated tests and reports coverage. It fails if any area drops
// below 100%, listing that area's functions that are short.
func Cover() error {
	var failed []string
	for _, g := range coverGates {
		if err := g.run(); err != nil {
			fmt.Println(err)
			failed = append(failed, g.name)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("coverage below 100%%: %s", strings.Join(failed, ", "))
	}
	return nil
}

// run measures one gate and reports whether it is short of 100%.
func (g coverGate) run() error {
	if err := os.MkdirAll("tmp", 0o755); err != nil {
		return err
	}
	// The profile is a throwaway artifact, so it is written under ./tmp (gitignored)
	// rather than the repo root.
	profile := filepath.Join("tmp", "coverage."+g.name+".out")
	// The per-package go test lines here report the coverpkg's coverage against a
	// different package's tests, which reads as confusing noise next to the real
	// figure below, so the output is captured and only surfaced when the run fails.
	if out, err := sh.Output(
		"go",
		"test",
		g.test,
		"-coverpkg="+g.count,
		"-coverprofile="+profile,
	); err != nil {
		fmt.Print(tidyTestOutput(out))
		return err
	}
	out, err := sh.Output("go", "tool", "cover", "-func="+profile)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	total := lines[len(lines)-1]
	fields := strings.Fields(total)
	pct, err := strconv.ParseFloat(strings.TrimSuffix(fields[len(fields)-1], "%"), 64)
	if err != nil {
		return fmt.Errorf("parsing coverage percent from %q: %w", total, err)
	}
	fmt.Printf("%-10s %.1f%%\n", g.name+":", pct)
	if pct < 100 {
		fmt.Printf("%s functions below 100%%:\n", g.name)
		for _, l := range lines[:len(lines)-1] {
			if !strings.HasSuffix(l, "100.0%") {
				fmt.Println("  " + l)
			}
		}
		return fmt.Errorf("%s coverage %.1f%% is below 100%%", g.name, pct)
	}
	return nil
}
