//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ScramblerStorm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent cannot play action cards on their next turn.
var ScramblerStorm = card.New(
	"Scrambler Storm",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 122),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
