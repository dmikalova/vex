package ageofascension

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
)

// Shard of Hate
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, stun an enemy creature.
var ShardOfHate = set.New(
	"Shard of Hate",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "205"),
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.ForEach{
			Times: card.CardsInPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
			Do: card.Stun{Target: card.Target.EnemyCreature},
		}),
)
