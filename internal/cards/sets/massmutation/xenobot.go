//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoBot
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
//	Reap: Discard a card from your hand. If you do, draw a card.
var XenoBot = set.New(
	"Xeno-Bot",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "124"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
