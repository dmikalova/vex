//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Quicksand
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy the most powerful creature controlled by each player who does not control a ready Untamed creature.
var Quicksand = card.New(
	"Quicksand",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 364),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
