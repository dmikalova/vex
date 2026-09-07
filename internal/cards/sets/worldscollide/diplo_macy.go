//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DiploMacy
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Alpha.
//	Play: Until the start of your next turn, each creature gains, "Before Fight: Exalt this creature."
var DiploMacy = card.New(
	"Diplo-Macy",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "218"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
