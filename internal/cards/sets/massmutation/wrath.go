//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Wrath
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Armor:  3
//	Traits: Demon • Sin
//
//	Taunt. Poison. Skirmish.
//	Fight: For each friendly Sin creature, enrage an enemy creature.
var Wrath = set.New(
	"Wrath",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "064"),
	card.WithPower(3),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
