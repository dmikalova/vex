//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// LittleNiff
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Elf • Thief
//
//	Omega. Deploy. Elusive.
//	After a neighbor of Little Niff is used to fight, steal 1A.
var LittleNiff = card.New(
	"Little Niff",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 289),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
