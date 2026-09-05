//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SirMarrows
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	After your opponent gains A by reaping, Sir Marrows captures it.
var SirMarrows = card.New(
	"Sir Marrows",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 223),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
