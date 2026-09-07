//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HuggerMugger
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Play: Capture 1A. Then, if your opponent has more forged keys than you, steal 1A.
var HuggerMugger = card.New(
	"Hugger-Mugger",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "240"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
