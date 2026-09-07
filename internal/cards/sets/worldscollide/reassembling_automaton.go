//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ReassemblingAutomaton
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot • Experiment
//
//	Destroyed: If you have any other creatures in play, instead of destroying Reassembling Automaton, fully heal it, exhaust it, and move it to a flank.
var ReassemblingAutomaton = card.New(
	"Reassembling Automaton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 158),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
