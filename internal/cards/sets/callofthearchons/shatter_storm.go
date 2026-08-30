//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ShatterStorm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Lose all your Aember. Then, your opponent loses triple the amount of Aember you lost this way.
var ShatterStorm = card.New(
	"Shatter Storm",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 176),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
