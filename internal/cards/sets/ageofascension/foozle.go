package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Foozle
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Reap: If an enemy creature has been destroyed this turn, gain 1 Æmber.
var Foozle = card.New(
	"Foozle",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 8),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Reap, card.Conditional{
			Cond: card.EnemyCreatureDestroyed{},
			Then: card.GainAember{Player: card.Controller, Amount: 1},
		}),
)
