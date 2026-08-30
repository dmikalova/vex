//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TakeThatSmartypants
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Steal 2 Aember if your opponent has 3 or more Logos cards in play.
var TakeThatSmartypants = card.New(
	"Take that, Smartypants",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 11),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
