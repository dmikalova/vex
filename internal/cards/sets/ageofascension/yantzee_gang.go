package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Yantzee Gang
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Elf • Thief
//
//	Action: Steal 1 Æmber.
var YantzeeGang = card.New(
	"Yantzee Gang",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 282),
	card.WithPower(5),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Action, card.StealAember{Amount: 1}),
)
