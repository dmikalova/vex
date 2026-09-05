//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// GenerateComments rewrites each card's doc comment from its definition.
func GenerateComments() error {
	return sh.RunV("go", "run", "./magefiles/gencomments")
}

// GenerateRules generates docs/rulebook.md from engine doc comments.
func GenerateRules() error {
	return sh.RunV("go", "run", "./magefiles/genrules")
}

// RulebookFresh fails if docs/rulebook.md is stale (out of sync with the engine
// doc comments it is generated from). It regenerates in memory and compares,
// leaving the working tree untouched. See ADR 0018.
func RulebookFresh() error {
	return sh.RunV("go", "run", "./magefiles/genrules", "-check")
}

// GenerateProvenance rebuilds the provenance card catalogs. Their source is the
// master-vault pack data. It is intentionally NOT part of Gen: it needs the external
// master-vault-data checkout, so it is run by hand when that source data changes.
func GenerateProvenance() error {
	return sh.RunV("go", "run", "./magefiles/genprovenance")
}

// Gen regenerates card comments and the rulebook.
func Gen() {
	mg.Deps(GenerateComments, GenerateRules)
}
