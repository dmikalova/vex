package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Gamgee
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: If your opponent has more Æmber than you, steal 1 Æmber.
var Gamgee = set.New(
	"Gamgee",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "270"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.MoreThanYou,
			},
			Then: card.StealAember{Amount: 1},
		}),
)
