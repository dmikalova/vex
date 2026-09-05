//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Foozle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Reap: If an enemy creature has been destroyed this turn, gain 1A.
var Foozle = card.New(
	"Foozle",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 8),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
