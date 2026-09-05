//go:build mage

// Vactrol developer tasks. Run `mage <target>` to run one, `mage -l` to list them.
package main

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build builds all packages, for the host and the browser. `go build
// ./...` compiles the web client for the host, which misses anything that only
// breaks under js/wasm, so WebWasm compiles it again for the target it ships to.
func Build() error {
	if err := sh.RunV("go", "build", "./..."); err != nil {
		return err
	}
	return WebWasm()
}

// Test runs all tests.
func Test() error {
	return sh.RunV("go", "test", "./...")
}

// TestRun runs only the tests matching a name pattern. It covers every package,
// and the pattern is a Go regexp — `mage testRun TestHeal`, or
// `mage testRun 'TestBumpsy|TestUrchin'`.
func TestRun(pattern string) error {
	return sh.RunV("go", "test", "./...", "-run", pattern)
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

// gciVersion pins the import organizer (also run via `go run`). gci groups imports
// into standard, third-party, then local-module sections — something gofmt/golines
// do not do.
const gciVersion = "v0.14.0"

// golinesCmd builds the `go run golines` invocation. gofmt is pinned as the base
// formatter so the result is identical whether or not goimports is on PATH.
func golinesCmd(extra ...string) []string {
	return append([]string{
		"run", "github.com/golangci/golines@" + golinesVersion,
		"--base-formatter=gofmt",
	}, extra...)
}

// gciCmd builds a `go run gci <sub>` invocation grouping imports as standard →
// third-party → local module. localmodule auto-detects the path from go.mod, so
// there is nothing to keep in sync.
func gciCmd(sub string) []string {
	return []string{
		"run", "github.com/daixiang0/gci@" + gciVersion, sub,
		"--skip-generated",
		"-s", "standard", "-s", "default", "-s", "localmodule", "--custom-order",
		".",
	}
}

// Fmt formats all Go files and organizes their imports. It runs golines (gofmt
// plus long-line shortening), then gci.
func Fmt() error {
	if err := sh.RunV("go", golinesCmd("-w", ".")...); err != nil {
		return err
	}
	return sh.RunV("go", gciCmd("write")...)
}

// FmtCheck checks formatting and import order, writing nothing. It fails
// if any Go file is not golines-clean or has unorganized imports, and each tool
// lists the files it would change.
func FmtCheck() error {
	golOut, err := sh.Output("go", golinesCmd("-l", ".")...)
	if err != nil {
		return err
	}
	gciOut, err := sh.Output("go", gciCmd("list")...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(golOut) != "" || strings.TrimSpace(gciOut) != "" {
		return fmt.Errorf("formatting needed:\n%s", strings.TrimSpace(golOut+"\n"+gciOut))
	}
	return nil
}

// Tidy tidies module dependencies.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// Check is the full green gate before calling work done. It runs fmt-check,
// build, vet, lint, markdown lint, test, coverage, and the rulebook freshness
// check.
func Check() error {
	mg.Deps(FmtCheck, Build, Vet, Lint, Markdownlint, Test, Cover, RulebookFresh)
	fmt.Println("ALL GREEN")
	return nil
}
