package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Angwish
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	Your opponent's keys cost +1 Æmber for each damage on it.
var Angwish = set.New(
	"Angwish",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "69"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.DamageOnThis{})),
)
