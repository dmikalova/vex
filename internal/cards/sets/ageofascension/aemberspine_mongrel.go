package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Aemberspine Mongrel
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Hazardous 3.
//	After an enemy creature reaps, gain 1 Æmber.
var AemberspineMongrel = set.New(
	"Aemberspine Mongrel",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "335"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithHazardous(3),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Conditional{
			Cond: card.ItIsEnemy{},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
