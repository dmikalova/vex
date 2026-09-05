//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Archimedes
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Beast
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Each of Archimedes's neighbors gains, "Destroyed: Archive this creature."
var Archimedes = card.New(
	"Archimedes",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 108),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
