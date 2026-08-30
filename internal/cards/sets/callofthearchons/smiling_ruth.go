//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SmilingRuth
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: If you forged a key this turn, take control of an enemy flank creature.
var SmilingRuth = card.New(
	"Smiling Ruth",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 312),
	card.WithPower(1),
	card.WithTraits("Elf", "Thief"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
