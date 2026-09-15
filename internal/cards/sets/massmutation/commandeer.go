//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Commandeer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For the remainder of the turn, after you play another card, a friendly creature captures 1A.
var Commandeer = set.New(
	"Commandeer",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "131"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
