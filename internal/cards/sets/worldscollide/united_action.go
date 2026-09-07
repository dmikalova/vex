//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UnitedAction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For the remainder of the turn, you may play cards from any house for which you have a card in play. You cannot use cards this turn.
var UnitedAction = card.New(
	"United Action",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 343),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
