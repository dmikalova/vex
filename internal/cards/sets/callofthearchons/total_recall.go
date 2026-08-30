//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TotalRecall
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For each friendly ready creature, gain 1 Aember. Return each friendly creature to your hand.
var TotalRecall = card.New(
	"Total Recall",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 179),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
