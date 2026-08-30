//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pitlord
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  9
//	Æmber:  2
//	Traits: Demon
//
//	Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)
//	While Pitlord is in play you must choose Dis as your active house.
var Pitlord = card.New(
	"Pitlord",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 93),
	card.WithPower(9),
	card.WithAemberBonus(2),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
