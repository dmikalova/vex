//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// EntropicManipulator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a player. You may redistribute the damage on the creatures that player controls among that player's creatures. (You may cause more damage to a creature than it has power.)
var EntropicManipulator = card.New(
	"Entropic Manipulator",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 195),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
