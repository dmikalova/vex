//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Equalize
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Redistribute the A on friendly creatures among friendly creatures. Then, redistribute the A on enemy creatures among enemy creatures.
var Equalize = card.New(
	"Equalize",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 232),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
