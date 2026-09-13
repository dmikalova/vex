package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Equalize
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Redistribute the Æmber on friendly Creatures among friendly Creatures. Redistribute the Æmber on enemy Creatures among enemy Creatures.
var Equalize = card.New(
	"Equalize",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "232"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.RedistributeCapturedAember{Side: card.Controller},
			card.RedistributeCapturedAember{Side: card.Opponent},
		}}),
)
