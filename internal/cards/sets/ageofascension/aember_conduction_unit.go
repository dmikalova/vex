//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AemberConductionUnit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After an enemy creature reaps, if it is the first time a creature has reaped this turn, stun it.
var AemberConductionUnit = card.New(
	"Aember Conduction Unit",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 176),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
