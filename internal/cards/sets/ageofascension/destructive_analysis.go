//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// DestructiveAnalysis
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 2D to a creature. You may purge any number of cards from your archives to deal an additional 2D to the same creature for each card purged this way.
var DestructiveAnalysis = card.New(
	"Destructive Analysis",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 194),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
