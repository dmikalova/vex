package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dark Faerie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Skirmish.
//	Fight: Gain 2 Æmber.
var DarkFaerie = set.New(
	"Dark Faerie",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "364"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.Controller,
			Amount: 2,
		}),
)
