//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MiniGroupthinkTank
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Robot • Experiment
//
//	Play/Fight/Reap: Deal 8D to a creature that shares a house with 2 of its neighbors.
var MiniGroupthinkTank = card.New(
	"Mini Groupthink Tank",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "155"),
	card.WithPower(3),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
