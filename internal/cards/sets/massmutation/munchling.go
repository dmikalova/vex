//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Munchling
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant
//
//	Skirmish.
//	Fight: You may discard a Logos card from your hand or archives. If you do, gain 1A.
var Munchling = set.New(
	"Munchling",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "076"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
