//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Restringuntus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Play: Choose a house. Your opponent cannot choose that house as their active house until Restringuntus leaves play.
var Restringuntus = card.New(
	"Restringuntus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 94),
	card.WithPower(1),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
