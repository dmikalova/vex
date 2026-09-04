package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dr. Escotera
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Play: For each forged key your opponent has, gain 1 Æmber.
var DrEscotera = card.New(
	"Dr. Escotera",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 140),
	card.WithPower(4),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per:    card.OpponentForgedKeys{},
		}),
)
