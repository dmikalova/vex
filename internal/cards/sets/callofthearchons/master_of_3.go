//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MasterOf3
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Variant
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with 3 power.
var MasterOf3 = card.New(
	"Master of 3",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Variant,
	card.Provenance(card.CotA, 91),
	card.WithPower(4),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
