//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HarbingerOfDoom
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  3
//	Traits: Demon
//
//	Destroyed: Destroy each creature.
var HarbingerOfDoom = card.New(
	"Harbinger of Doom",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 76),
	card.WithPower(2),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
