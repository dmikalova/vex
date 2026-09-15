//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TechnoThief
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Mutant • Thief
//
//	Elusive.
//	Reap: Discard a card from your hand. If you do, draw a card.
var TechnoThief = set.New(
	"Techno-Thief",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "300"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
