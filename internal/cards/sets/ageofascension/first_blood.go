//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// FirstBlood
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha. (You can only play this card before doing anything else this step.)
//	Play: Deal 2D for each friendly Brobnar creature. You may divide this damage among any number of creatures.
var FirstBlood = card.New(
	"First Blood",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 7),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
