//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// GreyAberrant
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Monk • Mutant
//
//	Each creature loses each of its traits.
var GreyAberrant = set.New(
	"Grey Aberrant",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "181"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Monk, card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
