//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gongoozle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3 Damage to a creature. If it is not destroyed, its owner discards a random card from their hand.
var Gongoozle = card.New(
	"Gongoozle",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 60),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
