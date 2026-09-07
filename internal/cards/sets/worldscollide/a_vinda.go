//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// AVinda
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Elf • Thief
//
//	Reap: Deal 1D to a creature. If this damage destroys that creature, your opponent discards a random card from their hand.
var AVinda = card.New(
	"A. Vinda",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 235),
	card.WithPower(4),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
