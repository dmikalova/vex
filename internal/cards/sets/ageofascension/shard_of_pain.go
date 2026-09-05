package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Pain
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, deal 1 damage to an enemy creature.
var ShardOfPain = card.New(
	"Shard of Pain",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 104),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.DealDamage{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
			Target: card.Target.EnemyCreature,
		}),
)
