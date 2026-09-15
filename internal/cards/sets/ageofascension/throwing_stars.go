package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Throwing Stars
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 1 damage to up to 3 creatures. For each creature destroyed this way, gain 1 Æmber.
var ThrowingStars = set.New(
	"Throwing Stars",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "279"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DealDamage{Spread: card.UpToCreatures{Count: 3, Amount: 1}},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CreaturesDestroyed{},
			},
		}},
	),
)
