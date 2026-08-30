//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Truebaru
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	You must lose 3 Aember in order to play Truebaru.
//	Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)
//	Destroyed: Gain 5 Aember.
var Truebaru = card.New(
	"Truebaru",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 104),
	card.WithPower(7),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
