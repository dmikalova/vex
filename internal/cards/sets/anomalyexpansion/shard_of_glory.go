package anomalyexpansion

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// Shard of Glory
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Connected
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, exalt an enemy creature.
var ShardOfGlory = set.New(
	"Shard of Glory",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Connected,
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.ForEach{
			Times: card.CardsInPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
			Do: card.Exalt{Target: card.Target.EnemyCreature, Amount: 1},
		}),
)
