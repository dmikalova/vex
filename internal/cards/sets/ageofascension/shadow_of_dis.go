//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ShadowOfDis
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Until your next turn, enemy creatures' text boxes are considered blank (except for traits).
var ShadowOfDis = card.New(
	"Shadow of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 103),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
