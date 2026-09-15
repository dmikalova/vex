//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mastermindy
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
//	Elusive.
//	At the end of your turn, put a scheme counter on Mastermindy.
//	Action: Remove each scheme counter from Mastermindy. Steal 1A for each scheme counter removed this way.
var Mastermindy = set.New(
	"Mastermindy",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "285"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
