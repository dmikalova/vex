//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Gluttony
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Demon • Sin
//
//	Play: Exalt Gluttony once for each friendly Sin creature.
//	Reap: Move each A from friendly creatures to your pool.
var Gluttony = set.New(
	"Gluttony",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "057"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
