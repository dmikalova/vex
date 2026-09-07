//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FesteringTouch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose up to 2 creatures. Deal 1D to each chosen creature. If that creature was already damaged, deal 3D instead.
var FesteringTouch = card.New(
	"Festering Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "75"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
