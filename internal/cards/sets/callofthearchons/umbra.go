package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Umbra
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Skirmish.
//	Fight: Steal 1 Æmber.
var Umbra = card.New(
	"Umbra",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 314),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.StealAember{Amount: 1}),
)
