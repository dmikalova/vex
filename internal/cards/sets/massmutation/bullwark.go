//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Bullwark
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Mutant • Knight
//
//	Assault 2. (Before this creature attacks, deal 2D to the attacked enemy.)
//	Each of Bull-wark's neighbors gains assault 2.
var Bullwark = set.New(
	"Bull-wark",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "127"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
