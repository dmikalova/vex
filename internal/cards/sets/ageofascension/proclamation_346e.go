package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Proclamation 346E
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	While your opponent does not control creatures from 3 or more different houses, your opponent's keys cost +2 Æmber.
var Proclamation346E = card.New(
	"Proclamation 346E",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "261"),
	card.WithTraits(card.Traits.Law),
	card.WithKeyCost(
		card.KeyCostChange(card.Opponent, 2).While(
			card.PlayerControlsFewerHousesThan{Player: card.Opponent, Amount: 3})),
)
