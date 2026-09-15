//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Vandalize
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Look at the top 3 cards of your opponent's deck. Discard 1 and put the others back in any order.
var Vandalize = set.New(
	"Vandalize",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "260"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
