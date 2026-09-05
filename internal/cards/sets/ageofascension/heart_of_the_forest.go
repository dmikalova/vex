//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// HeartOfTheForest
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Each player cannot forge keys while they have more forged keys than their opponent.
var HeartOfTheForest = card.New(
	"Heart of the Forest",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 355),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
