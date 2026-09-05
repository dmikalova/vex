//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GoldKeyImp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Variant
//	Power:  2
//	Traits: Imp
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Players cannot forge their third key.
var GoldKeyImp = card.New(
	"Gold Key Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Variant,
	card.Provenance(card.AoA, 72),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
