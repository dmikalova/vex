//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Malison
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
//	Fight: You may move an enemy creature anywhere in its controller's battleline. Then, if it is on a flank, it captures 1A from its own side.
var Malison = card.New(
	"Malison",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "80"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
