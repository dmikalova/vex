//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Adaptoid
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	Enhance Capture Damage Draw.
//	After you play a card with a bonus icon, for the remainder of the turn, Adaptoid gains (choose one): +2 armor, assault 2, or "Fight: Steal 1A."
var Adaptoid = set.New(
	"Adaptoid",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "099"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
