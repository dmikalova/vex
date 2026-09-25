package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Proclamation 346E
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	While fewer than 3 houses are represented among enemy creatures, your opponent's keys cost +2 Æmber.
var Proclamation346E = set.New(
	"Proclamation 346E",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "261"),
	card.WithTraits(card.Traits.Law),
	card.WithKeyCost(
		card.KeyCostChange(card.Opponent, 2).While(
			card.Not{Cond: card.CountIs{
				Count: card.HousesAmong{
					Player: card.Opponent,
					Type:   card.Type.Creature,
				},
				Is:     card.AtLeast,
				Amount: 3,
			}})),
)
