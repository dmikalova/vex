//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// EvenIvan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an even amount of A, steal 1A.
var EvenIvan = set.New(
	"Even Ivan",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "073"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
