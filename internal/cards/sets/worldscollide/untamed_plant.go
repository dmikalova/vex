//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UntamedPlant
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Variant
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	After a player chooses Untamed as their active house, gain 1A.
var UntamedPlant = card.New(
	"Untamed Plant",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 291),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
