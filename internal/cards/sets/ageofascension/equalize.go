package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Equalize
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Redistribute the Æmber on friendly creatures among friendly creatures. Redistribute the Æmber on enemy creatures among enemy creatures.
var Equalize = set.New(
	"Equalize",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "232"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RedistributeCapturedAember{Side: card.Controller},
			card.RedistributeCapturedAember{Side: card.Opponent},
		}}),
)
