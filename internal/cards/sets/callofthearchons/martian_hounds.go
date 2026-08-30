//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MartianHounds
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature. For each damaged creature, give the chosen creature two +1 power counters.
var MartianHounds = card.New(
	"Martian Hounds",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 167),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
