//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SignalFire
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Omni: Sacrifice Signal Fire. For the remainder of the turn, friendly Brobnar creatures may fight as though they belonged to the active house.
var SignalFire = card.New(
	"Signal Fire",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 49),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
