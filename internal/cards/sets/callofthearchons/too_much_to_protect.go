//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TooMuchToProtect
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Steal all but 6 of your opponent's Aember.
var TooMuchToProtect = card.New(
	"Too Much to Protect",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 283),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
