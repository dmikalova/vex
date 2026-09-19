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

// goldmarkLintVersion pins the markdown linter, run via `go run` so — unlike the
// old brew-installed quickmark — it is always available and the gate can never
// silently pass without checking Markdown.
const goldmarkLintVersion = "v0.5.3"

// Lint runs golangci-lint using the pinned version.
func Lint() error {
	return sh.RunV("go", "run",
		"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@"+golangciLintVersion,
		"run", "--fix")
}

// Markdownlint auto-fixes what it can and fails only on what it cannot. It runs
// goldmark-lint (a Go port of markdownlint) with --fix, so fixable issues are
// corrected in place — the same "format first, then check" contract as Fmt — and
// only genuinely unfixable violations exit non-zero. Scope, disabled rules, and
// the docs/todo.md exclusion live in .markdownlint-cli2.yaml.
func Markdownlint() error {
	return sh.RunV("go", "run",
		"github.com/mrueg/goldmark-lint/cmd/goldmark-lint@"+goldmarkLintVersion,
		"--fix", "**/*.md")
}

// Semgrep checks the tree with semgrep. It is deliberately NOT part of
// `mage check`: it is an external, non-Go binary that pulls remote rulesets, so
// it cannot be pinned via `go run` and is not always installed. Install it
// (`brew install semgrep` or `pipx install semgrep`) before running; without it
// this skips with a warning rather than failing.
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
