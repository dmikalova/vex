//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ogopogo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant
//
//	After Ogopogo attacks and destroys a creature, you may deal 2D to a creature.
var Ogopogo = card.New(
	"Ogopogo",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 26),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
