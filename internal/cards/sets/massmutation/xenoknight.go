//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var XenoKnight = set.New(
	"Xeno-Knight",
	card.House.Sanctum,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "182"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
