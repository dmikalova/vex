//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tremor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun a creature and each of its neighbors.
var Tremor = card.New(
	"Tremor",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 16),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
