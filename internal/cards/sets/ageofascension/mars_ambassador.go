//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MarsAmbassador
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Variant
//	Power:  1
//	Traits: Human
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Fight/Reap: You may play or use a Mars card this turn.
var MarsAmbassador = card.New(
	"Mars Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Variant,
	card.Provenance(card.AoA, 238),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
