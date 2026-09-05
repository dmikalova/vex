//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PanpacaJaga
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	Creatures to the left of Panpaca,
//	Jaga in the battleline gain skirmish.
var PanpacaJaga = card.New(
	"Panpaca, Jaga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 348),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
