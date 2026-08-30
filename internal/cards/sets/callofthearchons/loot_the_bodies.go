//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LootTheBodies
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, gain 1 Aember each time an enemy creature is destroyed.
var LootTheBodies = card.New(
	"Loot the Bodies",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 10),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
