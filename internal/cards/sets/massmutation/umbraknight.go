//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// UmbraKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Elusive.
var UmbraKnight = set.New(
	"Umbra-Knight",
	card.House.Sanctum,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "181"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
