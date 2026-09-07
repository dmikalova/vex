package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Harbinger of Doom
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
	card.WithAbility(
		card.Trigger.Destroyed, card.Destroy{Target: card.Target.EachCreature}),
)
