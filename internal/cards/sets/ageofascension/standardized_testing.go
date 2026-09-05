//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// StandardizedTesting
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each creature with the lowest power and each creature with the highest power.
var StandardizedTesting = card.New(
	"Standardized Testing",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 119),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
