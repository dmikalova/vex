//go:build mage

package main

import (
	"fmt"

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
// interactive ↑/↓ picker; pass -set=<slug> to name one directly, e.g.
// `mage tool:missing -set=callofthearchons` (slugs match the files in
// internal/cards/provenance, minus the .json).
func (Tool) Missing(set *string) error {
	args := []string{"run", "./magefiles/cardlookup", "missing"}
	if set != nil && *set != "" {
		args = append(args, *set)
	}
	return sh.RunV("go", args...)
}

// Coverage reports how many cards of each set are implemented. A card
// counts as covered once an implemented card tags it with a provenance Ref.
// Pass -new to count only the cards a set introduces, excluding the ones it
// reprints from an earlier set.
func (Tool) Coverage(new *bool) error {
	args := []string{"run", "./magefiles/cardlookup", "coverage"}
	if new != nil && *new {
		args = append(args, "-new")
	}
	return sh.RunV("go", args...)
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

// NextCard prints the next unimplemented card whose stub file still carries the
// `//go:build todo` constraint, walking a set's missing cards in collector-number
// order and stopping at the first one on disk. It is the pick-the-next-card step
// of the implement-cards workflow: build the card it names, drop the build tag,
// and run it again for the next. With no set chosen it opens the interactive ↑/↓
// picker; pass -set=<slug> to name one directly, e.g. `mage tool:nextCard
// -set=ageofascension`.
func (Tool) NextCard(set *string) error {
	args := []string{"run", "./magefiles/cardlookup", "next-card"}
	if set != nil && *set != "" {
		args = append(args, *set)
	}
	return sh.RunV("go", args...)
}

// NodeUsage reports how widely each card-facade node is used. It prints every
// exported name in internal/card, grouped by the category its declaration block
// documents, with the number of card definitions, total occurrences, and sets
// that name it — rarest first, then a summary of the whole facade. Low usage is
// not a defect on its own (half the card pool is unimplemented); it marks the
// nodes to check are built from reusable atoms rather than hard-coding one card.
// Pass -max=<n> to show only the nodes at most n cards use, and
// -category=<substring> to narrow to one group, e.g.
// `mage tool:nodeUsage -max=1 -category=damage`.
func (Tool) NodeUsage(max *int, category *string) error {
	args := []string{"run", "./magefiles/cardlookup", "node-usage"}
	if max != nil && *max >= 0 {
		args = append(args, fmt.Sprintf("-max=%d", *max))
	}
	if category != nil && *category != "" {
		args = append(args, "-category="+*category)
	}
	return sh.RunV("go", args...)
}

// ImportProvenance rebuilds a set's source catalog from the Master Vault decks
// feed. It pages the feed for the set's expansion, folds each linked card into the
// catalog shape (ASCII-folded name and text, expanded amber/damage markup,
// normalized house and rarity, classified anomalies), and writes
// internal/cards/provenance/<slug>.json. With no set chosen it opens an
// interactive picker offering every set plus "All sets"; pass -set=<slug> to name
// one, or -set=all to rebuild every set in release order (each fetched on its own).
func (Tool) ImportProvenance(set *string) error {
	args := []string{"run", "./magefiles/cardlookup", "import-provenance"}
	if set != nil && *set != "" {
		args = append(args, *set)
	}
	return sh.RunV("go", args...)
}
