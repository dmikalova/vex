//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mindworm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Beast
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Before Fight: The creature Mindworm fights deals damage equal to its power to each of its neighbors.
var Mindworm = card.New(
	"Mindworm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 168),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
