//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Pride
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Demon • Sin
//
//	Reap: Ward each friendly Sin creature.
var Pride = set.New(
	"Pride",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "060"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
