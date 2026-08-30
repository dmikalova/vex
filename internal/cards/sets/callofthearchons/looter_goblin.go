//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LooterGoblin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: For the remainder of the turn, gain 1 Aember each time an enemy creature is destroyed.
var LooterGoblin = card.New(
	"Looter Goblin",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 41),
	card.WithPower(2),
	card.WithTraits("Goblin"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
