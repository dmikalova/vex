//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Scout
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For the remainder of the turn, up to 2 friendly creatures gain skirmish. Then, fight with those creatures one at a time.
var Scout = card.New(
	"Scout",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 334),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
