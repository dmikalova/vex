//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BloodshardImp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After a creature reaps, its controller must sacrifice it.
var BloodshardImp = card.New(
	"Bloodshard Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 70),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
