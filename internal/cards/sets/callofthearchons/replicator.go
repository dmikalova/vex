//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Replicator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Mutant
//
//	Reap: Trigger the reap effect of another creature in play as if you controlled that creature. (That creature does not exhaust.)
var Replicator = card.New(
	"Replicator",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 150),
	card.WithPower(2),
	card.WithTraits("Mutant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
