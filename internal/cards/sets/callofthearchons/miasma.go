//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Miasma
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent skips the "forge a key" step on their next turn.
var Miasma = card.New(
	"Miasma",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 275),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
