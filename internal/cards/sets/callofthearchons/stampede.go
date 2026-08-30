//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Stampede
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you used 3 or more creatures this turn, steal 2 Aember.
var Stampede = card.New(
	"Stampede",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 335),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
