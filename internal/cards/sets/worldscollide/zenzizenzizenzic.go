package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Zenzizenzizenzic
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Cyborg • Leader
//
//	While Zenzizenzizenzic is in the center of the battleline, your hand size is 2 more.
var Zenzizenzizenzic = set.New(
	"Zenzizenzizenzic",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "180"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Leader),
	card.WithDrawModifierInCenter(card.Controller, 2),
)
