//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Picaroon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  0
//	Traits: Mutant • Changeling
//
//	Deploy.
//	X is the combined power of Picaroon's non-Changeling neighbors.
var Picaroon = set.New(
	"Picaroon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "028"),
	card.WithPower(0),
	card.WithTraits(card.Traits.Mutant, card.Traits.Changeling),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
