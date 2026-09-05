//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShardOfStrength
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: Give a friendly creature a +1 power counter for each friendly Shard.
var ShardOfStrength = card.New(
	"Shard of Strength",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 48),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
