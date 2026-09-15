//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroBeast
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Armor:  2
//	Traits: Mutant • Beast
//
//	Skirmish.
var SacroBeast = set.New(
	"Sacro-Beast",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "418"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
