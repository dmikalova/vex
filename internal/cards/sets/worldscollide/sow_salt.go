//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SowSalt
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Alpha.
//	Play: Until the start of your next turn, creatures cannot be used to reap.
var SowSalt = card.New(
	"Sow Salt",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "230"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
