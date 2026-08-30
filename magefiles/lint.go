//go:build mage

package main

import "github.com/magefile/mage/sh"

// golangciLintVersion pins the linter so local, hook, and CI runs all agree.
// Run via `go run` so it never enters the module's own dependency graph.
const golangciLintVersion = "v2.13.2"

// Lint runs golangci-lint using the pinned version.
func Lint() error {
	return sh.RunV("go", "run",
		"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@"+golangciLintVersion,
		"run")
}
