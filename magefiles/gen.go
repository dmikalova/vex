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

// Gen regenerates card comments.
func Gen() {
	mg.Deps(GenerateComments)
}
