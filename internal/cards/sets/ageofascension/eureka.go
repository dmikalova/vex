package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eureka!
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Alpha.
//	Play: Gain 2 Æmber. Archive 2 random cards from your hand.
var Eureka = set.New(
	"Eureka!",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "128"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.GainAember{Player: card.Controller, Amount: 2},
				card.ArchiveCard{
					Zone:      card.Hand,
					Selection: card.Random{},
					Quantity:  card.Takes{N: card.Fixed(2)},
				},
			},
		}),
)
