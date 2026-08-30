//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mimicry
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	When you play this card, treat it as a copy of an action card in your opponent's discard pile.
var Mimicry = card.New(
	"Mimicry",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 328),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
