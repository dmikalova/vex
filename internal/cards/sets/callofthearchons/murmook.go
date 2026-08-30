package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Your opponent's keys cost +1 Æmber.
var Murmook = card.New(
	"Murmook",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 361),
	card.WithPower(3),
	card.WithTraits("Beast"),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1)),
)
