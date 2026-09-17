package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Sir Marrows
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	After an enemy creature reaps, Sir Marrows captures 1 Æmber from your opponent.
var SirMarrows = set.New(
	"Sir Marrows",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "223"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Conditional{
			Cond: card.ItIsEnemy{},
			Then: card.CaptureAember{
				Amount: 1,
				Target: card.Target.This,
				Source: card.Opponent,
			},
		}),
)
