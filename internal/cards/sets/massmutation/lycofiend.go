//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LycoFiend
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Mutant • Demon
//
//	Skirmish.
//	Destroyed: Steal 1A.
var LycoFiend = set.New(
	"Lyco-Fiend",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "059"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
