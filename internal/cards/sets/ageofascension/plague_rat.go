//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PlagueRat
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Rat
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Play: Each non-Rat creature is dealt 1D for each Rat creature in play.
var PlagueRat = card.New(
	"Plague Rat",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 308),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Rat),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
