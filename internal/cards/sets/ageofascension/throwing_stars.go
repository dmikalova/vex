//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ThrowingStars
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 1D to up to 3 creatures. Gain 1A for each creature destroyed this way.
var ThrowingStars = card.New(
	"Throwing Stars",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 279),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
