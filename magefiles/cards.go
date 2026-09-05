//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Tool groups the card-management commands used while researching and
// implementing cards — lookup, missing, coverage, and stub. They are invoked
// under the namespace, e.g. `mage tool:stub callofthearchons`.
type Tool mg.Namespace

// Lookup prints the source cards whose name contains a query. Each hit carries
// its set code, collector number, house/type/rarity, printed text, and a
// ready-made card.Provenance(...) call — the details a card.New(...) definition
// needs. Run `mage tool:lookup "ether spider"`.
func (Tool) Lookup(query string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "lookup", query)
}

// Missing lists a set's cards still to implement. Those are the source cards no
// implemented card tags with a provenance Ref yet. With no set chosen it opens an
// interactive ↑/↓ picker; set SET=<slug> to name one directly, e.g.
// `SET=callofthearchons mage tool:missing` (slugs match the files in
// internal/cards/provenance, minus the .json).
func (Tool) Missing() error {
	args := []string{"run", "./magefiles/cardlookup", "missing"}
	if set := os.Getenv("SET"); set != "" {
		args = append(args, set)
	}
	return sh.RunV("go", args...)
}

// Coverage reports how many cards of each set are implemented. A card
// counts as covered once an implemented card tags it with a provenance Ref.
func (Tool) Coverage() error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "coverage")
}

// Stub scaffolds a stub for each unimplemented card in a set. Each stub is
// build-excluded (`//go:build todo`) and carries the card's printed text and a
// TODO marker. Excluded stubs do not compile or register, so the card database and
// coverage numbers stay honest until a card is actually implemented; to implement
// one, remove the build tag and write the real ability. It also (re)generates the
// set package's `0set.go`, cataloging the cards this set reprints from earlier
// sets so they join its deck-generation pool as full members (ADR 0021). Pass a
// set slug, e.g. `mage tool:stub callofthearchons`.
func (Tool) Stub(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "stub", setSlug)
}
