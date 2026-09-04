package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dust Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: Gain 2 Æmber.
var DustImp = card.New(
	"Dust Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 83),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithAbility(
		card.Trigger.Destroyed, card.GainAember{
			Player: card.Controller,
			Amount: 2,
		}),
)
