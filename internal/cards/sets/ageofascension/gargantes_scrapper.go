package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gargantes Scrapper
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Æmber:  1
//	Traits: Giant
//
//	Alpha.
//	Play: For each Æmber in your pool, deal 3 damage to an enemy creature.
var GargantesScrapper = card.New(
	"Gargantes Scrapper",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 23),
	card.WithAemberBonus(1),
	card.WithPower(3),
	card.WithTraits(card.Traits.Giant),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Per:    card.AemberInPool{Player: card.Controller},
			Target: card.Target.EnemyCreature,
		}),
)
