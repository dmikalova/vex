//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ThreeFates
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy the 3 most powerful creatures.
var ThreeFates = card.New(
	"Three Fates",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 71),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
