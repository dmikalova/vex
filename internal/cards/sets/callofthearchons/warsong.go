//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Warsong
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, gain 1 Aember each time a friendly creature fights.
var Warsong = card.New(
	"Warsong",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 18),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
