//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// KeyHammer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent forged a key on their previous turn, unforge it. Your opponent gains 6 Aember.
var KeyHammer = card.New(
	"Key Hammer",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 66),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
