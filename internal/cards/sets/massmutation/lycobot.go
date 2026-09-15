//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LycoBot
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
//	Skirmish.
//	Reap: Discard a card from your hand. If you do, draw a card.
var LycoBot = set.New(
	"Lyco-Bot",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "120"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
