//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheShadowCouncil
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Leader • Thief
//
//	Elusive.
//	While The Shadow Council is in the center of your battleline, it gains, "Action: Steal 2A."
var TheShadowCouncil = card.New(
	"The Shadow Council",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 283),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Leader, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
