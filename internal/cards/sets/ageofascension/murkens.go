//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Murkens
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
//	Play: Choose a random card in your opponent's archives or the top card of your opponent's deck. Play that card as if it were yours.
var Murkens = card.New(
	"Murkens",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 290),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
