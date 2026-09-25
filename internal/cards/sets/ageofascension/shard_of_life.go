package ageofascension

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
)

// Shard of Life
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, shuffle a card from your discard pile into your deck.
var ShardOfLife = set.New(
	"Shard of Life",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "366"),
	card.InCluster(clusters.Shard),
	card.OneCopyPerDeck(),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.ShuffleIntoDeck{
			Player: card.Controller, From: []card.Zone{card.Discard},
			Selection: card.Chosen{},
			Quantity: card.Takes{
				N: card.CardsInPlay{
					Player: card.Controller,
					Trait:  card.Traits.Shard,
				},
			},
		}),
)
