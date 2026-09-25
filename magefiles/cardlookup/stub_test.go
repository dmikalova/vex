package main

import (
	"testing"

	"github.com/dmikalova/vex/internal/cards"
	"github.com/dmikalova/vex/internal/cards/provenance"
)

// Master of X is one implemented card standing in for three printings, which it
// tags with a provenance Ref each. Nothing shares its name, so only those Refs
// keep Master of 1/2/3 out of `mage tool:missing` and stop `mage tool:stub` scaffolding a
// stub file for a card that is already implemented.
func TestMasterOfXCoversItsNumberedPrintings(t *testing.T) {
	_ = cards.All()
	covered := coveredNumbers()[provenance.CallOfTheArchons.Slug]

	for number, name := range map[string]string{
		"089": "Master of 1",
		"090": "Master of 2",
		"091": "Master of 3",
	} {
		if !covered[number] {
			t.Errorf(
				"CotA #%s (%s) is uncovered, so mage tool:stub would scaffold it",
				number,
				name,
			)
		}
	}
}
