//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Symon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Alien • Thief
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	Fight: Put the creature Symon fights on top of its owner's deck.
var Symon = card.New(
	"Symon",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 247),
	card.WithPower(1),
	card.WithTraits(card.Traits.Alien, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
