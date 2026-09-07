//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SaurianPlant
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	After a player chooses Saurian as their active house, gain 1A.
var SaurianPlant = card.New(
	"Saurian Plant",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "289"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
