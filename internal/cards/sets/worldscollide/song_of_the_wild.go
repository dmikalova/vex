//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SongOfTheWild
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each friendly creature gains, "Reap: Gain 1A."
var SongOfTheWild = card.New(
	"Song of the Wild",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 364),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
