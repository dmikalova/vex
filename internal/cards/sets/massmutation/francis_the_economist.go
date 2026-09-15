package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Francis the "Economist"
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	Fight: Each player gains 1 Æmber.
var FrancisTheEconomist = set.New(
	"Francis the \"Economist\"",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "248"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.EachPlayer,
			Amount: 1,
		}),
)
