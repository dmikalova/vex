//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Proclamation346E
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	While your opponent does not control creatures from 3 different houses, their keys cost +2A.
var Proclamation346E = card.New(
	"Proclamation 346E",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 261),
	card.WithTraits(card.Traits.Law),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
