//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Resurgence
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Enhance Draw.
//	Play: Return a creature from your discard pile to your hand. If that creature is a Mutant, return another creature from your discard pile to your hand.
var Resurgence = set.New(
	"Resurgence",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "375"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
