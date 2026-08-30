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

// Fmt formats all Go files.
func Fmt() error {
	return sh.RunV("gofmt", "-w", ".")
}

// FmtCheck fails if any Go file is not gofmt-clean (does not modify files).
func FmtCheck() error {
	out, err := sh.Output("gofmt", "-l", ".")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("gofmt needed:\n%s", out)
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
