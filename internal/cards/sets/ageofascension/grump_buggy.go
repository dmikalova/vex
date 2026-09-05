//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GrumpBuggy
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Vehicle
//
//	Your opponent's keys cost +1A for each friendly creature with power 5 or higher.
//	Your keys cost +1A for each enemy creature with power 5 or higher.
var GrumpBuggy = card.New(
	"Grump Buggy",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 24),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Vehicle),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
