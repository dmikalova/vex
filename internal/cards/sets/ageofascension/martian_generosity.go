//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MartianGenerosity
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Lose all of your A. Draw 2 cards for each A lost.
var MartianGenerosity = card.New(
	"Martian Generosity",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 202),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
