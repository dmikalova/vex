//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// UxlyxTheZookeeper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Martian • Scientist
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Put an enemy creature into your archives. If that creature leaves your archives, it is put into its owner's hand instead.
var UxlyxTheZookeeper = card.New(
	"Uxlyx the Zookeeper",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 201),
	card.WithPower(2),
	card.WithTraits("Martian", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
