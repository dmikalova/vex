package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Free Markets
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each house represented among cards in play, except for Sanctum, gain 1 Æmber.
var FreeMarkets = card.New(
	"Free Markets",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "233"),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per:    card.HousesInPlay{Except: card.House.Self},
		}),
)
