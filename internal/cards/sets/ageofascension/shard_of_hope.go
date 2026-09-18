package ageofascension

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// Shard of Hope
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, a friendly creature captures 1 Æmber from your opponent.
var ShardOfHope = set.New(
	"Shard of Hope",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "264"),
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 1,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
			Times: card.CardsInPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
		}),
)
