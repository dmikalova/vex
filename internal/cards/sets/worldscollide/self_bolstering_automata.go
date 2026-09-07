//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SelfBolsteringAutomata
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Destroyed: If you have any other creatures in play, instead of destroying Self-Bolstering Automata, fully heal it, exhaust it, and move it to a flank. If you do, give it two +1 power counters.
var SelfBolsteringAutomata = card.New(
	"Self-Bolstering Automata",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 176),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
