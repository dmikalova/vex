//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// RecklessRizzo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 2A. Until the start of your next turn, Reckless Rizzo loses elusive.
var RecklessRizzo = set.New(
	"Reckless Rizzo",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "273"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
