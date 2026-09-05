//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// IntoTheFray
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, a friendly Brobnar creature gains, "Fight: Ready this creature."
var IntoTheFray = card.New(
	"Into the Fray",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 13),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
