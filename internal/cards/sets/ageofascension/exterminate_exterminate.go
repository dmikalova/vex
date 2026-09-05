//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ExterminateExterminate
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each friendly Mars creature you control, destroy a non-Mars creature with lower power.
var ExterminateExterminate = card.New(
	"Exterminate! Exterminate!",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 180),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
