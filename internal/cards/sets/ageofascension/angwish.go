//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Angwish
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	For each damage on Angwish,
//	your opponent's keys cost +1A.
var Angwish = card.New(
	"Angwish",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 69),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
