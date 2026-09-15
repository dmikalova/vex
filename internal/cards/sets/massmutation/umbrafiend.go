//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// UmbraFiend
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  2
//	Traits: Mutant • Demon
//
//	Elusive.
//	Destroyed: Steal 1A.
var UmbraFiend = set.New(
	"Umbra-Fiend",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "063"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
