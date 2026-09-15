//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroFiend
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Armor:  2
//	Traits: Mutant • Demon
//
//	Destroyed: Steal 1A.
var SacroFiend = set.New(
	"Sacro-Fiend",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "061"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
