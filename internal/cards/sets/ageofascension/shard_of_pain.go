//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfPain
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Deal 1D to an enemy creature for each friendly Shard.
var ShardOfPain = card.New(
	"Shard of Pain",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 104),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
