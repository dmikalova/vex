//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MassAbduction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Put up to 3 damaged enemy creatures into your archives. If any of these creatures leave your archives, they are put into their owner's hand instead.
var MassAbduction = card.New(
	"Mass Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 169),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
