//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Card12Punch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Stun an enemy creature.
//	If that creature was already stunned, destroy it instead.
var Card12Punch = card.New(
	"1-2 Punch",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 1),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
