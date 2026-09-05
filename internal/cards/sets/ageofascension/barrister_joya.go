//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BarristerJoya
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Human • Knight
//
//	Enemy creatures cannot reap.
var BarristerJoya = card.New(
	"Barrister Joya",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 228),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
