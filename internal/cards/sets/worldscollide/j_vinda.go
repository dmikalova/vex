//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// JVinda
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Deal 1D to a creature. If this damage destroys that creature, steal 1A.
var JVinda = card.New(
	"J. Vinda",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 242),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
