//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hecatomb
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each Dis creature. Each player gains 1 Aember for each creature they controlled that was destroyed this way.
var Hecatomb = card.New(
	"Hecatomb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 63),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
