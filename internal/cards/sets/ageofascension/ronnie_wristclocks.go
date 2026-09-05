//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// RonnieWristclocks
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Steal 1A. If your opponent has 7A or more, steal 2A instead.
var RonnieWristclocks = card.New(
	"Ronnie Wristclocks",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 276),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
