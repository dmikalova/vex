package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Headhunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Fight: Gain 1 Æmber.
var Headhunter = card.New(
	"Headhunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 35),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.Controller,
			Amount: 1,
		}),
)
