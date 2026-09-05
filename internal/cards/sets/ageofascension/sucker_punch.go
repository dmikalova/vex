//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SuckerPunch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha. (You can only play this card before doing anything else this step.)
//	Play: Deal 2D to an enemy creature.
//	If that creature is destroyed by this effect, archive Sucker Punch.
var SuckerPunch = card.New(
	"Sucker Punch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 277),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
