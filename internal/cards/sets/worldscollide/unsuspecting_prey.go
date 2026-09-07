//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UnsuspectingPrey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2D to up to 3 undamaged creatures.
var UnsuspectingPrey = card.New(
	"Unsuspecting Prey",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 368),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
