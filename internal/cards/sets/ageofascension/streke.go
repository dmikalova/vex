//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Streke
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
//	While Streke is not on a flank, your opponent refills their hand to 1 less card during their "draw cards" step.
var Streke = card.New(
	"Streke",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 65),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
