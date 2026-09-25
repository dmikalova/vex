package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Jammer Pack
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Your opponent's keys cost +2 Æmber."
var JammerPack = set.New(
	"Jammer Pack",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "210"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{KeyCostChange: card.KeyCostChange(card.Opponent, 2)}),
)
