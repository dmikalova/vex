//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Zenzizenzizenzic
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Cyborg • Leader
//
//	During your "draw cards" step, if Zenzizenzizenzic is in the center of your battleline, refill your hand to 2 additional cards.
var Zenzizenzizenzic = card.New(
	"Zenzizenzizenzic",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 180),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Leader),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
