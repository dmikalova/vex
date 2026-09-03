//go:build mage

package main

import "github.com/magefile/mage/sh"

// Lookup prints the source cards whose name contains a query. Each hit carries
// its set code, collector number, house/type/rarity, printed text, and a
// ready-made card.Provenance(...) call — the details a card.New(...) definition
// needs. Run `mage lookup "ether spider"`.
func Lookup(query string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "lookup", query)
}

// Missing lists a set's source cards that are still to implement. Those are the
// ones no implemented card tags with a provenance Ref yet. Pass a set slug, e.g.
// `mage missing callofthearchons` (slugs match the files in
// internal/cards/provenance, minus the .json).
func Missing(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "missing", setSlug)
}

// Coverage reports how many cards of each set are implemented. A card
// counts as covered once an implemented card tags it with a provenance Ref.
func Coverage() error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "coverage")
}

// Stub scaffolds a stub for each unimplemented card in a set. Each stub is
// build-excluded (`//go:build todo`) and carries the card's printed text and a
// TODO marker. Excluded stubs do not compile or register, so the card database and
// coverage numbers stay honest until a card is actually implemented; to implement
// one, remove the build tag and write the real ability. Pass a set slug, e.g.
// `mage stub callofthearchons`.
func Stub(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "stub", setSlug)
}
