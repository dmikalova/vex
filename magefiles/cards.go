//go:build mage

package main

import "github.com/magefile/mage/sh"

// Lookup prints every source card whose name contains the query, with its set
// code, collector number, house/type/rarity, printed text, and a ready-made
// card.Provenance(...) call — the details a card.New(...) definition needs. Run
// `mage lookup "ether spider"`.
func Lookup(query string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "lookup", query)
}

// Missing lists the source cards in a set that no implemented card tags with a
// provenance Ref yet — the cards still to implement. Pass a set slug, e.g.
// `mage missing callofthearchons` (slugs match the files in
// internal/cards/provenance, minus the .json).
func Missing(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "missing", setSlug)
}

// Coverage prints, per source set, how many of its cards are covered by an
// implemented card's provenance Ref.
func Coverage() error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "coverage")
}

// Stub generates a build-excluded (`//go:build todo`) stub file for every
// unimplemented card in a set, each carrying the card's printed text and a TODO
// marker. Excluded stubs do not compile or register, so the card database and
// coverage numbers stay honest until a card is actually implemented; to implement
// one, remove the build tag and write the real ability. Pass a set slug, e.g.
// `mage stub callofthearchons`.
func Stub(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "stub", setSlug)
}
