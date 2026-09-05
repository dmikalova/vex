//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// DuskChronicles
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has more A than you, draw a card. If you have more A than your opponent, archive a card.
var DuskChronicles = card.New(
	"Dusk Chronicles",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 268),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
