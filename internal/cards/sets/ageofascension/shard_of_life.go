package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Life
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, shuffle a card from your discard pile into your deck.
var ShardOfLife = card.New(
	"Shard of Life",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 366),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.ShuffleCardsFromDiscard{
			Count: card.InPlay{Player: card.Controller, Trait: card.Traits.Shard},
		}),
)
