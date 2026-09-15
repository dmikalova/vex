//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Survey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Enhance Draw.
//	Play: Look at the top 2 cards of your deck. Discard 1 of them.
var Survey = set.New(
	"Survey",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "316"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
