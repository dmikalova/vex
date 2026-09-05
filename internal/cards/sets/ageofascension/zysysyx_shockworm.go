//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ZysysyxShockworm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Martian • Soldier
//
//	After an enemy creature reaps, stun it.
var ZysysyxShockworm = card.New(
	"Zysysyx Shockworm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 175),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
