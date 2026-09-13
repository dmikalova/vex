package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Hate
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, stun an enemy Creature.
var ShardOfHate = card.New(
	"Shard of Hate",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "205"),
	card.InCluster(shardCluster),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.ForEach{
			Times: card.InPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
			Do: card.Stun{Target: card.Target.EnemyCreature},
		}),
)
