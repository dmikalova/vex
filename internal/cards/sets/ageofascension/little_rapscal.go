//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// LittleRapscal
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Goblin
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Creatures must fight when used, if able.
var LittleRapscal = card.New(
	"Little Rapscal",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 25),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
