//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LibraryAccess
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time you play another card, draw a card.
var LibraryAccess = card.New(
	"Library Access",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 115),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
