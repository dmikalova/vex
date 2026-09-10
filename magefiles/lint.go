//go:build mage

package main

import (
	"fmt"
	"os/exec"
	"strings"

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

// Markdownlint runs quickmark over every non-ignored markdown file except
// docs/todo.md. That file is the human's personal scratch list — freeform, often
// mid-edit, and never kept lint-clean — and agents never write into it (see the
// repo AGENTS.md), so holding it to the doc lint standard would only ever fail the
// gate on the human's placeholders. Install qmark (`brew install quickmark-cli` or
// see https://github.com/ekropotin/quickmark) before running — unlike
// golangci-lint it is not a Go module mage can pin via `go run`, so environments
// without it (e.g. CI images that haven't installed it yet) skip with a warning
// instead of failing the whole gate.
func Markdownlint() error {
	if _, err := exec.LookPath("qmark"); err != nil {
		fmt.Println("qmark not found on PATH, skipping markdown lint")
		return nil
	}
	// Match qmark's own default scope — tracked plus untracked files, minus
	// anything .gitignore excludes — so a new uncommitted doc is still linted.
	out, err := sh.Output(
		"git", "ls-files", "--cached", "--others", "--exclude-standard", "*.md",
	)
	if err != nil {
		return err
	}
	var files []string
	for _, f := range strings.Split(strings.TrimSpace(out), "\n") {
		if f == "" || f == "docs/todo.md" {
			continue
		}
		files = append(files, f)
	}
	if len(files) == 0 {
		return nil
	}
	return sh.RunV("qmark", files...)
}
