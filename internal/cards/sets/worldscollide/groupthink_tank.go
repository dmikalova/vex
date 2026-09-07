//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GroupthinkTank
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  3
//	Traits: Robot • Experiment
//
//	Action: Deal 4D to each creature that shares a house with at least 1 of its neighbors.
var GroupthinkTank = card.New(
	"Groupthink Tank",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "151"),
	card.WithPower(4),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
