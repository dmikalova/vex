//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// LethalDistraction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a creature. For the remainder of the turn, whenever this creature takes damage, it takes an additional 2D.
var LethalDistraction = card.New(
	"Lethal Distraction",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 305),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
