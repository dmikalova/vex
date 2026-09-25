//go:build mage

// Vex developer tasks. Run `mage <target>` to run one, `mage -l` to list them.
// The generic gate (ci:fix, ci:check, ci:build, ci:test, ci:cover, ci:lint and
// the rest) comes from project-standards' ci package; the targets here are the
// ones only vex has. The local validator is `mage ci:fix && mage ci:check`.
package main

import (
	"context"
	"os"
	"path/filepath"

	// mage:import ci
	"github.com/dmikalova/project-standards/ci"
	"github.com/magefile/mage/sh"
)

func init() {
	// `go build ./...` compiles the web client for the host, which misses
	// anything that only breaks under js/wasm, so ci:build compiles it again for
	// the target it ships to. It does not precompress the static assets: that
	// step (WebAssets) is slow max-level brotli, only the shipped server needs it,
	// and the server falls back to the raw files when the .br/.gz siblings are
	// absent. The Dockerfile runs `mage webAssets` explicitly.
	ci.ExtraBuilds = []func(context.Context) error{buildWasmCheck}

	// The areas held at 100% statement coverage. test names the packages whose
	// tests run and count the packages whose statements are counted. They differ
	// where a package is exercised through another's tests: the card definitions
	// are covered by the per-card tests that sit beside them, so the cards gate
	// counts the sets while running everything under internal/cards. Absent, and
	// so deliberately ungated: internal/web, a view layer the tests exercise
	// through only a few entry points.
	ci.CoverGates = []ci.CoverGate{
		{Name: "engine", Test: "./internal/engine/", Count: "./internal/engine/"},
		{Name: "cards", Test: "./internal/cards/...", Count: "./internal/cards/sets/..."},
		{Name: "cardtest", Test: "./internal/cards/...", Count: "./internal/cards/cardtest/"},
		{Name: "deckgen", Test: "./internal/deckgen/", Count: "./internal/deckgen/"},
	}
}

// buildWasmCheck compiles the web client for js/wasm into a temporary
// directory, so ci:check stays read-only. WebWasm is the target that writes the
// bundle the server ships.
func buildWasmCheck(context.Context) error {
	dir, err := os.MkdirTemp("", "vex-wasm-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(dir) }()
	return wasmBuild(filepath.Join(dir, "app.wasm"))
}

// TestRun runs only the tests matching a name pattern. It covers every package,
// and the pattern is a Go regexp — `mage testRun TestHeal`, or
// `mage testRun 'TestBumpsy|TestUrchin'`.
func TestRun(pattern string) error {
	return goTest("./...", "-run", pattern)
}

// Fmtmklv checks or rewrites multiline keyed composite literals. It defaults to
// checking the whole repo and accepts a -fix flag to rewrite in place, with an
// optional path argument for focused runs.
func Fmtmklv(path *string, fix *bool) error {
	args := []string{"run", "./magefiles/mklvfmt"}
	if fix != nil && *fix {
		args = append(args, "-fix")
	}
	if path != nil && *path != "" {
		args = append(args, *path)
	} else {
		args = append(args, "./...")
	}
	return sh.RunV("go", args...)
}
