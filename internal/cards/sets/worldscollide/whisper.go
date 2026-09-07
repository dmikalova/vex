//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Whisper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Lose 1A. If you do, destroy a creature.
var Whisper = card.New(
	"Whisper",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 265),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
