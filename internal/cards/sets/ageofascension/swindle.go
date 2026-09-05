package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Swindle
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha, Omega.
//	Play: Steal 3 Æmber.
var Swindle = card.New(
	"Swindle",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 278),
	card.WithKeywords(card.Keyword.Alpha, card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{Amount: 3}),
)
