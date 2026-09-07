//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FavorOfRex
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Trigger the play effect of a creature as if you had just played it.
var FavorOfRex = card.New(
	"Favor of Rex",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "219"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
