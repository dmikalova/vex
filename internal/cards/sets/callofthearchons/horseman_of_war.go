//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// HorsemanOfWar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: FIXED
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: For the remainder of the turn, each friendly creature can be used as if they were in the active house, but can only fight.
var HorsemanOfWar = card.New(
	"Horseman of War",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.CotA, 249),
	card.WithPower(5),
	card.WithTraits("Horseman", "Spirit"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
