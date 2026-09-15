//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Bonesaw
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	If a friendly creature was destroyed this turn, Bonesaw enters play ready.
var Bonesaw = set.New(
	"Bonesaw",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "002"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
