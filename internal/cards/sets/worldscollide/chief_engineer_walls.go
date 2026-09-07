//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ChiefEngineerWalls
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Play/Fight/Reap: You may return an upgrade or Robot card from your discard pile to your hand.
var ChiefEngineerWalls = card.New(
	"Chief Engineer Walls",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "293"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
