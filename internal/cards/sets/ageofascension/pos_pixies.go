//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PosPixies
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Faerie
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	A stolen or captured from your pool is taken from the common supply instead.
var PosPixies = card.New(
	"Po's Pixies",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 362),
	card.WithPower(1),
	card.WithTraits(card.Traits.Faerie),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
