//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PawnSacrifice
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Sacrifice a friendly creature. If you do, deal 3 Damage each to 2 creatures.
var PawnSacrifice = card.New(
	"Pawn Sacrifice",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 279),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
