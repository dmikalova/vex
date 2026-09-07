//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// WeCanALLWin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each player's keys cost -2A until the end of your next turn.
var WeCanALLWin = card.New(
	"We Can ALL Win",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 344),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
