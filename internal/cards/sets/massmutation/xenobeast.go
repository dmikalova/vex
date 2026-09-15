//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoBeast
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Mutant • Beast
//
//	Skirmish.
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var XenoBeast = set.New(
	"Xeno-Beast",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "421"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
