//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AbondTheArmorsmith
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Other friendly creatures get +1 armor.
//	Action: For the remainder of the turn, other friendly creatures get +1 armor.
var AbondTheArmorsmith = card.New(
	"Abond the Armorsmith",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 212),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
