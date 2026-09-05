//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ScowlyCaper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Skirmish.
//	Scowly Caper enters play under your opponent's control and can be used as if it belonged to any house.
//	At the end of your turn, destroy one of
//	Scowly Caper's neighbors.
var ScowlyCaper = card.New(
	"Scowly Caper",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 313),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
