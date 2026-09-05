package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Before Fight: Deal 2 damage to each neighbor of the creature Cowfyne fights.
var Cowfyne = card.New(
	"Cowfyne",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 5),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 2,
			Target: card.Target.CreatureFought.NeighborsOf(),
		}),
)
