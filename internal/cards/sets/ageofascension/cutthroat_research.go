package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Cutthroat Research
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: If your opponent has 8 Æmber or more, steal 2 Æmber.
var CutthroatResearch = set.New(
	"Cutthroat Research",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "110"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.AtLeast,
				Amount: 8,
			},
			Then: card.StealAember{Amount: 2},
		}),
)
