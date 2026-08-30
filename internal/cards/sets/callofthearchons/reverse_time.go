//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ReverseTime
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Swap your deck and your discard pile. Then, shuffle your deck.
var ReverseTime = card.New(
	"Reverse Time",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 121),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
