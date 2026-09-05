package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bingle Bangbang
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Goblin
//
//	Before Fight: Deal 5 damage to each neighbor of the creature Bingle Bangbang fights.
var BingleBangbang = card.New(
	"Bingle Bangbang",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 2),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 5,
			Target: card.Target.CreatureFought.NeighborsOf(),
		}),
)
