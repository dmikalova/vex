package main

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// Master of X is one implemented card standing in for three printings, which it
// tags with a provenance Ref each. Nothing shares its name, so only those Refs
// keep Master of 1/2/3 out of `mage missing` and stop `mage stub` scaffolding a
// stub file for a card that is already implemented.
func TestMasterOfXCoversItsNumberedPrintings(t *testing.T) {
	_ = cards.All()
	covered := coveredNumbers()[provenance.CallOfTheArchons.Slug]

	for number, name := range map[int]string{89: "Master of 1", 90: "Master of 2", 91: "Master of 3"} {
		if !covered[number] {
			t.Errorf("CotA #%d (%s) is uncovered, so mage stub would scaffold it", number, name)
		}
	}
}
