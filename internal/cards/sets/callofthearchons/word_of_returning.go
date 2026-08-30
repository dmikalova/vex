//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// WordOfReturning
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 1 Damage to each enemy creature for each Aember on it. Return all Aember from those creatures to your pool.
var WordOfReturning = card.New(
	"Word of Returning",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 339),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
