//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Camouflage
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	Creatures not on a flank cannot fight this creature.
var Camouflage = card.New(
	"Camouflage",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 337),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
