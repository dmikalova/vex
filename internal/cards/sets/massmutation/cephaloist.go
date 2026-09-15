//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Cephaloist
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	While you have 4A or more, your A cannot be stolen.
var Cephaloist = set.New(
	"Cephaloist",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "362"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
