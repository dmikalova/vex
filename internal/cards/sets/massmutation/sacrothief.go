//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroThief
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Armor:  2
//	Traits: Mutant • Thief
//
//	Elusive.
var SacroThief = set.New(
	"Sacro-Thief",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "299"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
