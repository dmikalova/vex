//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MarsNeedsAember
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Each damaged enemy non-Mars creature captures 1A from their own side.
var MarsNeedsAember = card.New(
	"Mars Needs Aember",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 166),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
