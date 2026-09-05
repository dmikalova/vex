package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Binate Rupture
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha.
//	Play: For each Æmber in your pool, gain 1 Æmber, and for each Æmber in your opponent's pool, your opponent gains 1 Æmber.
var BinateRupture = card.New(
	"Binate Rupture",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 109),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.GainAember{
					Player: card.Controller,
					Amount: 1,
					Per:    card.AemberInPool{Player: card.Controller},
				},
				card.GainAember{
					Player: card.Opponent,
					Amount: 1,
					Per:    card.AemberInPool{Player: card.Opponent},
				},
			},
		}),
)
