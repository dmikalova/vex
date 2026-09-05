//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BindingIrons
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent gains 3 chains.
var BindingIrons = card.New(
	"Binding Irons",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 55),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
