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

// GenerateProvenance rebuilds the provenance card catalogs. Their source is the
// master-vault pack data. It is intentionally NOT part of Gen: it needs the external
// master-vault-data checkout, so it is run by hand when that source data changes.
func GenerateProvenance() error {
	return sh.RunV("go", "run", "./magefiles/genprovenance")
}

// Gen regenerates card comments.
func Gen() {
	mg.Deps(GenerateComments)
}
