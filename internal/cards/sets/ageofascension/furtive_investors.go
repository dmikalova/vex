package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Furtive Investors
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, for each forged key your opponent has, gain 1 Æmber.
var FurtiveInvestors = card.New(
	"Furtive Investors",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 269),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.MoreThanYou},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.OpponentForgedKeys{},
			},
		}),
)
