//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MarkOfDis
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2D to a creature. If it is not destroyed, its controller must choose that creature's house as their active house on their next turn.
var MarkOfDis = set.New(
	"Mark of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "011"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
