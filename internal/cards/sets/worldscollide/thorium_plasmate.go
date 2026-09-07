//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ThoriumPlasmate
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Move an enemy creature anywhere in its controller's battleline. Deal 2D to that creature for each of its neighbors that shares a house with it.
var ThoriumPlasmate = card.New(
	"Thorium Plasmate",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 140),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
