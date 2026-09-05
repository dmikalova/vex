package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ogopogo
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Ogopogo, you may deal 2 damage to a creature.
var Ogopogo = card.New(
	"Ogopogo",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 26),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.May{
			Do: card.DealDamage{
				Target: card.Target.Creature,
				Amount: 2,
			},
		}),
)
