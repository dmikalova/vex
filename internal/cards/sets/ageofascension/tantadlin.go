//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Tantadlin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Tree
//
//	Tantadlin only deals 2D when fighting.
//	Fight: Discard a random card from your opponent's archives.
var Tantadlin = card.New(
	"Tantadlin",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 333),
	card.WithPower(9),
	card.WithTraits(card.Traits.Tree),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
