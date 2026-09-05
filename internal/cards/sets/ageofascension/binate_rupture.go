//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BinateRupture
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha. (You can only play this card before doing anything else this step.)
//	Play: Each player gains A equal to
//	the A in their pool.
var BinateRupture = card.New(
	"Binate Rupture",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 109),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
