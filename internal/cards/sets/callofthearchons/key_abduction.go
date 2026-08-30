//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// KeyAbduction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Return each Mars creature to its owner's hand. Then, you may forge a key at +9 Aember current cost, reduced by 1 Aember for each card in your hand.
var KeyAbduction = card.New(
	"Key Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 166),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
