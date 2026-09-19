package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Binate Rupture
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha.
//	Play: Each player gains Æmber equal to the Æmber in their pool.
var BinateRupture = set.New(
	"Binate Rupture",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "109"),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player:  card.EachPlayer,
			EqualTo: card.AemberInPool{Player: card.Controller},
		}),
)
