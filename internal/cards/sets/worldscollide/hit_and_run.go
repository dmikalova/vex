//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HitAndRun
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2D to a creature. Return a friendly creature to your hand.
var HitAndRun = card.New(
	"Hit and Run",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 238),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
