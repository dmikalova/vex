//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Envy
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Demon • Sin
//
//	Elusive.
//	Reap: If there are 2 or more friendly Sin creatures, capture all of your opponent's A.
var Envy = set.New(
	"Envy",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "056"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
