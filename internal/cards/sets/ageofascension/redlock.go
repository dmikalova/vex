//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Redlock
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	At the end of your turn, if you did not play any creatures this turn, gain 1A.
var Redlock = card.New(
	"Redlock",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 294),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
