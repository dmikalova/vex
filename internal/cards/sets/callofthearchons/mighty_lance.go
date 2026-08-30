//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MightyLance
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 3 Damage to a creature and 3 Damage to a neighbor of that creature.
var MightyLance = card.New(
	"Mighty Lance",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 221),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
