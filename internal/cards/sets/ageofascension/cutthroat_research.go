package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Cutthroat Research
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has 8 Æmber or more, steal 2 Æmber.
var CutthroatResearch = card.New(
	"Cutthroat Research",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 110),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{
				Is:     card.AtLeast,
				Amount: 8,
			},
			Then: card.StealAember{Amount: 2},
		}),
)
