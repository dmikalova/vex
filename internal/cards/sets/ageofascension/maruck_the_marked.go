//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MaruckTheMarked
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Spirit • Knight
//
//	After Maruck the Marked prevents damage with its armor, capture 1A for each damage just prevented.
var MaruckTheMarked = card.New(
	"Maruck the Marked",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 220),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
