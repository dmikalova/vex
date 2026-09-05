//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Smite
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ready and fight with a friendly creature. Deal 2D to the attacked creature's neighbors.
var Smite = card.New(
	"Smite",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 224),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
