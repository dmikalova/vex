//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TreasureMap
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you have not played any other cards this turn, gain 3 Aember. For the remainder of the turn, you cannot play cards.
var TreasureMap = card.New(
	"Treasure Map",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 284),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
