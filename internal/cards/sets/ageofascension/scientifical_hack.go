//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ScientificalHack
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Equation
//
//	Omni: Sacrifice Scientifical Hack. For the remainder of the turn, you may use friendly artifacts as if they belonged to the active house.
var ScientificalHack = card.New(
	"Scientifical Hack",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 154),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Equation),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
