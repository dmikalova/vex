//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BreakerHill
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Each of Breaker Hill's neighbors gains, "Action: Steal 1A."
var BreakerHill = card.New(
	"Breaker Hill",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 237),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
