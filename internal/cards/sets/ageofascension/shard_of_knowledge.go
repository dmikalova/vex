package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Knowledge
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, draw a card.
var ShardOfKnowledge = card.New(
	"Shard of Knowledge",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 155),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.Draw{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
		}),
)
