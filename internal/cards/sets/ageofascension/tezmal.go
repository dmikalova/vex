//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Tezmal
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Choose a house. Your opponent cannot choose that house as their active house on their next turn.
var Tezmal = card.New(
	"Tezmal",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 66),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
