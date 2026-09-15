//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroBot
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Armor:  2
//	Traits: Mutant • Scientist
//
//	Reap: Discard a card from your hand. If you do, draw a card.
var SacroBot = set.New(
	"Sacro-Bot",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "122"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
