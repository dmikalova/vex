//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// FreeMarkets
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Gain 1A (to a maximum of 6) for each house represented among cards in play, except for Sanctum.
var FreeMarkets = card.New(
	"Free Markets",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 233),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
