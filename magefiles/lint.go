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

// Semgrep runs the semgrep static analyzer over the tree with its curated Go
// ruleset (p/golang). It is deliberately NOT part of `mage check`: it is an
// external, non-Go binary that pulls remote rulesets, so it cannot be pinned via
// `go run` and is not always installed. Install it (`brew install semgrep` or
// `pipx install semgrep`) before running; without it this skips with a warning
// rather than failing.
//
// Two of p/golang's rules are excluded because they are categorically wrong for
// this codebase, not one-off suppressions:
//   - math-random-used ("use crypto/rand"): every rand draw here is deterministic
//     seeded randomness for reproducible games and deck generation (ADR 0005,
//     0039). crypto/rand is unseedable, so it would break the reproducibility the
//     test/sim/debug/fuzz infrastructure relies on. This randomness is never
//     security-sensitive (no tokens, nonces, or secrets).
//   - use-tls ("use ListenAndServeTLS"): the server runs behind Cloud Run, which
//     terminates TLS at its front end and forwards plain HTTP, and locally over
//     plain HTTP. Serving TLS from the container would break both.
//
// The one real finding — filepath.Clean on a request path — is confined by the
// web/ prefix guard at its call site and suppressed there with an inline
// nosemgrep comment, so that rule stays active for any future code.
func Semgrep() error {
	if _, err := exec.LookPath("semgrep"); err != nil {
		fmt.Println("semgrep not found on PATH, skipping (brew install semgrep)")
		return nil
	}
	return sh.RunV("semgrep", "scan", "--error", "--config", "p/golang",
		"--exclude-rule", "go.lang.security.audit.crypto.math_random.math-random-used",
		"--exclude-rule", "go.lang.security.audit.net.use-tls.use-tls",
		".")
}
