//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// UmbraBot
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Elusive.
//	Reap: Discard a card from your hand. If you do, draw a card.
var UmbraBot = set.New(
	"Umbra-Bot",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "123"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
