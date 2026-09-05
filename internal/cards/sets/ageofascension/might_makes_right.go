//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MightMakesRight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: You may sacrifice any number of creatures with total power of 25 or more. If you do, forge a key at no cost.
var MightMakesRight = card.New(
	"Might Makes Right",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 43),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
