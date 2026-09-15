//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoThief
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Mutant • Thief
//
//	Elusive.
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var XenoThief = set.New(
	"Xeno-Thief",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "301"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
