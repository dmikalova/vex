//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfHope
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: A friendly creature captures
//	1A for each friendly Shard.
var ShardOfHope = card.New(
	"Shard of Hope",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 264),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
