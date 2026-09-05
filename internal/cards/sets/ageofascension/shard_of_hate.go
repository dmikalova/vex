//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfHate
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Stun an enemy creature for each friendly Shard.
var ShardOfHate = card.New(
	"Shard of Hate",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 205),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
