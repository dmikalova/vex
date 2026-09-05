package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Storm Crawler
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Robot
//
//	Storm Crawler deals 1 Damage when fighting.
//	After an enemy creature reaps, stun it.
var StormCrawler = card.New(
	"Storm Crawler",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 189),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 1,
		Fixed:  true,
	}),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
