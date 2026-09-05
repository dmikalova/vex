//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfLife
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Shuffle a card from your discard pile into your deck for each friendly Shard.
var ShardOfLife = card.New(
	"Shard of Life",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 366),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
