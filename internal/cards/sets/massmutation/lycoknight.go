//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LycoKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Skirmish.
var LycoKnight = set.New(
	"Lyco-Knight",
	card.House.Sanctum,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "179"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
