//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PanpacaAnga
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Creatures to the right of Panpaca, Anga in the battleline get +2 power.
var PanpacaAnga = card.New(
	"Panpaca, Anga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 347),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
