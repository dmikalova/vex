//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// NyzykResonator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Martian • Soldier
//
//	For each neighbor Nyzyk Resonator has, your opponent's keys cost +2A.
var NyzykResonator = card.New(
	"Nyzyk Resonator",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 184),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
