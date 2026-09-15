//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Sloth
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Demon • Sin
//
//	At the end of your turn, if you did not use any creatures this turn, gain 1A for each friendly Sin creature.
var Sloth = set.New(
	"Sloth",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "062"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
