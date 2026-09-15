//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// OddClawde
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an odd amount of A, steal 1A.
var OddClawde = set.New(
	"Odd Clawde",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "121"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
