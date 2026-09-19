package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dodger's 10
//
//	House:  Shadows
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  11
//	Traits: Elf • Thief
//
//	Play/Fight/Reap: Steal half of your opponent's Æmber, rounded down.
var Dodgers10 = set.Gigantic(
	"Dodger's 10",
	card.House.Shadows,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "258"),
	card.WithPower(11),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.PlayFightReap, card.StealAember{By: card.HalfRoundedDown}),
)
