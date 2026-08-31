//go:build mage

// Vactrol developer tasks. Run `mage` to list targets, `mage <target>` to run
// one, or `mage -l` for the same list with descriptions.
package main

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build builds all packages.
func Build() error {
	return sh.RunV("go", "build", "./...")
}

// Test runs all tests.
func Test() error {
	return sh.RunV("go", "test", "./...")
}

// Vet runs go vet.
func Vet() error {
	return sh.RunV("go", "vet", "./...")
}

// golinesVersion pins the formatter so local, hook, and CI runs all agree. Run
// via `go run` (as with the linter) so it never enters the module's own
// dependency graph. golines applies gofmt's formatting and additionally shortens
// long lines.
const golinesVersion = "v0.15.0"

// golinesCmd builds the `go run golines` invocation. gofmt is pinned as the base
// formatter so the result is identical whether or not goimports is on PATH.
func golinesCmd(extra ...string) []string {
	return append([]string{
		"run", "github.com/golangci/golines@" + golinesVersion,
		"--base-formatter=gofmt",
	}, extra...)
}

// Fmt formats all Go files with golines (gofmt plus long-line shortening).
func Fmt() error {
	return sh.RunV("go", golinesCmd("-w", ".")...)
}

// FmtCheck fails if any Go file is not golines-clean (does not modify files).
func FmtCheck() error {
	out, err := sh.Output("go", golinesCmd("--dry-run", ".")...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("golines needed:\n%s", out)
	}
	return nil
}

// Tidy tidies module dependencies.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// Check is the full green gate: fmt-check, build, vet, lint, test, coverage.
func Check() error {
	mg.Deps(FmtCheck, Build, Vet, Lint, Test, Cover)
	fmt.Println("ALL GREEN")
	return nil
}
