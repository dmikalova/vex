//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LavaBall
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 4 Damage to a creature with 2 Damage splash.
var LavaBall = card.New(
	"Lava Ball",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 9),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
