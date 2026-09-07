//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DeepwoodDruid
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Witch
//
//	Deploy. (This creature can enter play anywhere in your battleline.)
//	Play/Reap: Fully heal a neighboring creature.
var DeepwoodDruid = card.New(
	"Deepwood Druid",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 355),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
