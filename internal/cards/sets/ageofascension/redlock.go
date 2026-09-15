package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Redlock
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	At the end of your turn, if you did not play any creatures this turn, gain 1 Æmber.
var Redlock = set.New(
	"Redlock",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "294"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Conditional{
			Cond: card.NoCreaturesPlayedThisTurn{},
			Then: card.GainAember{Player: card.Controller, Amount: 1},
		}),
)
