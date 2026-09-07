//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// StarAlliancePlant
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
//	After a player chooses Star Alliance as their active house, gain 1A.
var StarAlliancePlant = card.New(
	"Star Alliance Plant",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "290"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
