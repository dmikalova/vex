//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Extinction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature. Destroy that creature and each creature that shares a trait with it. Gain 1 chain.
var Extinction = card.New(
	"Extinction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 196),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
