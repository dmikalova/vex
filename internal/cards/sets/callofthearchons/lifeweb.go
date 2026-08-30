//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lifeweb
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent played 3 or more creatures on their previous turn, steal 2 Aember.
var Lifeweb = card.New(
	"Lifeweb",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 326),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
