package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Grump Buggy
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Vehicle
//
//	Your opponent's keys cost +1 Æmber for each friendly creature with power 5 or higher.
//	Your keys cost +1 Æmber for each enemy creature with power 5 or higher.
var GrumpBuggy = card.New(
	"Grump Buggy",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "24"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Vehicle),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.InPlay{
		Player:   card.Controller,
		Type:     card.Type.Creature,
		MinPower: 5,
	})),
	card.WithKeyCost(card.KeyCostChange(card.Controller, 1).Per(card.InPlay{
		Player:   card.Opponent,
		Type:     card.Type.Creature,
		MinPower: 5,
	})),
)
