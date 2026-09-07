//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// InkyGloom
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use creatures to reap on their next turn.
var InkyGloom = card.New(
	"Inky Gloom",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 241),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
