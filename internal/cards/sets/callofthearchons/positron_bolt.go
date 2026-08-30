//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PositronBolt
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3 Damage to a flank creature. Deal 2 Damage to its neighbor. Deal 1 Damage to the second creature's other neighbor.
var PositronBolt = card.New(
	"Positron Bolt",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 118),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
