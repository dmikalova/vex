package anomalyexpansion

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// Shard of Unity
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Connected
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, use a friendly creature.
var ShardOfUnity = set.New(
	"Shard of Unity",
	card.House.StarAlliance,
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
			Do: card.OnChooseCreature{
				Target: card.Target.FriendlyCreature,
				Verbs:  []card.CreatureVerb{card.UseVerb{}},
			},
		}),
)
