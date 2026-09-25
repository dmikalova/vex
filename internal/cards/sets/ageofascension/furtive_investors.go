package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Furtive Investors
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: If your opponent has more Æmber than you, for each forged key your opponent has, gain 1 Æmber.
var FurtiveInvestors = set.New(
	"Furtive Investors",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "269"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.MoreThanYou,
			},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.ForgedKeys{Player: card.Opponent},
			},
		}),
)
