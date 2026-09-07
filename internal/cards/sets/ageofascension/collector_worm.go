package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Collector Worm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  5
//	Traits: Beast
//
//	Fight: Put the creature Collector Worm fought into your archives.
var CollectorWorm = card.New(
	"Collector Worm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "162"),
	card.WithPower(2),
	card.WithArmor(5),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Fight, card.PutFromPlay{
			Target:      card.Target.CreatureFought,
			Destination: card.To.Archives.Yours(),
		}),
)
