package ageofascension

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
)

// Shard of Strength
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, give a friendly creature three +1 power counters.
var ShardOfStrength = set.New(
	"Shard of Strength",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "48"),
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.AddPowerCounter{
			Target: card.Target.FriendlyCreature,
			Amount: 3,
			Per: card.CardsInPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
		}),
)
