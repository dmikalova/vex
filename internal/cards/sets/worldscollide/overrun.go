//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Overrun
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If 3 or more enemy creatures have been destroyed this turn, your opponent loses 2A.
var Overrun = card.New(
	"Overrun",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 25),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
