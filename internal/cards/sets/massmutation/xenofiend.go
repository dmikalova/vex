//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoFiend
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
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
//	Destroyed: Steal 1A.
var XenoFiend = set.New(
	"Xeno-Fiend",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "065"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
