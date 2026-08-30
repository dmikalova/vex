//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PhylyxTheDisintegrator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Martian • Soldier
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Action: Your opponent loses 1 Aember for each other friendly Mars creature.
var PhylyxTheDisintegrator = card.New(
	"Phylyx the Disintegrator",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 197),
	card.WithPower(1),
	card.WithTraits("Martian", "Soldier"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
