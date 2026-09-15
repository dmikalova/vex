//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TechnoBeast
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
//	Reap: Discard a card from your hand. If you do, draw a card.
var TechnoBeast = set.New(
	"Techno-Beast",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "419"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
