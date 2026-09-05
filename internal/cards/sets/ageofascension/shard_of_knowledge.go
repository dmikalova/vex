//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfKnowledge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Draw a card for each friendly Shard.
var ShardOfKnowledge = card.New(
	"Shard of Knowledge",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 155),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
