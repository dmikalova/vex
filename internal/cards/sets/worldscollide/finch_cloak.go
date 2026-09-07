//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FinchCloak
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Elf • Thief
//
//	Fight/Reap: If you have less A than your opponent, steal 1A. Otherwise, each player gains 1A.
var FinchCloak = card.New(
	"Finch Cloak",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 267),
	card.WithPower(4),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
