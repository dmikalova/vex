//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// JONCargo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Reap: Discard the top card of your deck and reveal your hand. Archive each card that shares a house with the discarded card.
var JONCargo = set.New(
	"J.O.N. Cargo",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "347"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
