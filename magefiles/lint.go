//go:build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/sh"
)

// golangciLintVersion pins the linter so local, hook, and CI runs all agree.
// Run via `go run` so it never enters the module's own dependency graph.
const golangciLintVersion = "v2.13.2"

// Lint runs golangci-lint using the pinned version.
func Lint() error {
	return sh.RunV("go", "run",
		"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@"+golangciLintVersion,
		"run")
}

// Markdownlint runs quickmark to lint all markdown files. Install it (`brew install quickmark-cli` or
// see https://github.com/ekropotin/quickmark) before running this target — unlike
// golangci-lint, it is not a Go module mage can pin via `go run`, so environments
// without it (e.g. CI images that haven't installed it yet) skip with a warning
// instead of failing the whole gate.
func Markdownlint() error {
	if _, err := exec.LookPath("qmark"); err != nil {
		fmt.Println("qmark not found on PATH, skipping markdown lint")
		return nil
	}
	return sh.RunV("qmark")
}
