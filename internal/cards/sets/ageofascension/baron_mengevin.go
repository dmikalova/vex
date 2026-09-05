//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BaronMengevin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Human • Knight
//
//	After you discard a Sanctum card from your hand, Baron Mengevin captures 1A.
var BaronMengevin = card.New(
	"Baron Mengevin",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 227),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
