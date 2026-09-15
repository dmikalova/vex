//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Desire
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
//	Keys cost +4A.
//	Reap: Forge a key at current cost, reduced by 1A for each friendly Sin creature.
var Desire = set.New(
	"Desire",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "053"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
