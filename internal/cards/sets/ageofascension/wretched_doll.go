//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// WretchedDoll
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: If there is a doom counter in play, destroy all creatures with doom counters. Otherwise, put a doom counter on a creature.
var WretchedDoll = card.New(
	"Wretched Doll",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 107),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
