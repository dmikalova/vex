//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SpecialAgentFingers
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 1A.
var SpecialAgentFingers = card.New(
	"Special Agent \"Fingers\"",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 339),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
