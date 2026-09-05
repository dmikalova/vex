//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Nightforge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have not forged a key
//	this turn, you may forge a key at
//	+4A current cost.
var Nightforge = card.New(
	"Nightforge",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 291),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
