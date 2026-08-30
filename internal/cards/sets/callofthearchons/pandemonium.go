//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pandemonium
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Each undamaged creature captures 1 Aember from its opponent.
var Pandemonium = card.New(
	"Pandemonium",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 68),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
