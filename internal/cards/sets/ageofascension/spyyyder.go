//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Spyyyder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Demon
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	Spyyyder gains poison while attacking an enemy flank creature.
var Spyyyder = card.New(
	"Spyyyder",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 84),
	card.WithPower(2),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
