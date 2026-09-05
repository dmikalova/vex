//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Exhume
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a creature in your discard pile. You may play that creature as if it belonged to the active house and was in your hand.
var Exhume = card.New(
	"Exhume",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 59),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
