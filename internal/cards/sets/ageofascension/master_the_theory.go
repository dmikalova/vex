//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MasterTheTheory
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, you may archive a card for each enemy creature.
var MasterTheTheory = card.New(
	"Master the Theory",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 148),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
