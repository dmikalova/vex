//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gub
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Demon
//
//	While Gub is not on a flank, it gets +5 power and gains taunt.
var Gub = card.New(
	"Gub",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 60),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
