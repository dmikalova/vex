package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dew Faerie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Faerie
//
//	Elusive.
//	Reap: Gain 1 Æmber.
var DewFaerie = card.New(
	"Dew Faerie",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 350),
	card.WithPower(2),
	card.WithTraits(card.Traits.Faerie),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{
			Player: card.Controller,
			Amount: 1,
		}),
)
