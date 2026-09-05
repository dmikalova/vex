//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfGreed
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Gain 1A for each friendly Shard.
var ShardOfGreed = card.New(
	"Shard of Greed",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 315),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
