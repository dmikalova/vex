package ageofascension

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
)

// Shard of Pain
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, deal 1 damage to an enemy creature.
var ShardOfPain = set.New(
	"Shard of Pain",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "104"),
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.DealDamage{
			Amount: 1,
			Per: card.CardsInPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
			Target: card.Target.EnemyCreature,
		}),
)
