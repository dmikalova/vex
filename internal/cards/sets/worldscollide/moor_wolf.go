//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MoorWolf
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast • Wolf
//
//	Skirmish.
//	Play: Ready each other Wolf creature.
var MoorWolf = card.New(
	"Moor Wolf",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 393),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Wolf),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
