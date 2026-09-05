package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Nyzyk Resonator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Martian • Soldier
//
//	For each neighbor Nyzyk Resonator has, your opponent's keys cost +2 Æmber.
var NyzykResonator = card.New(
	"Nyzyk Resonator",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 184),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 2).Per(card.NeighborsOfThis{})),
)
