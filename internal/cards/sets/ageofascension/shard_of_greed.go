package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Greed
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, gain 1 Æmber.
var ShardOfGreed = card.New(
	"Shard of Greed",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 315),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
		}),
)
