//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Etaromme
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Reap: Destroy a creature of the house with the most creatures in play.
var Etaromme = card.New(
	"Etaromme",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 73),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
