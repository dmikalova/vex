//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// UmbraBeast
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Mutant • Beast
//
//	Elusive.
//	Skirmish.
var UmbraBeast = set.New(
	"Umbra-Beast",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "420"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
