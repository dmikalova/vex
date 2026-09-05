//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eureka
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Alpha.
//	Play: Gain 2A. Archive 2 random cards from your hand.
var Eureka = card.New(
	"Eureka!",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 128),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
