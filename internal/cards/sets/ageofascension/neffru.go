//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Neffru
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Each time a creature is destroyed,
//	its owner gains 1A.
var Neffru = card.New(
	"Neffru",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 94),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
