//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GleefulMayhem
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house, deal 5D to a creature of that house.
var GleefulMayhem = card.New(
	"Gleeful Mayhem",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 90),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
