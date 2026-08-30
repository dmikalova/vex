//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RoutineJob
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Steal 1 Aember. Then, steal 1 Aember for each copy of Routine Job in your discard pile.
var RoutineJob = card.New(
	"Routine Job",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 282),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
