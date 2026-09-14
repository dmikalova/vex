package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Jammer Pack
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Your opponent's keys cost +2 Æmber."
var JammerPack = set.New(
	"Jammer Pack",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "210"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{KeyCostChange: card.KeyCostChange(card.Opponent, 2)}),
)
