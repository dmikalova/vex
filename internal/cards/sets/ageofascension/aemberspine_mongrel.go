package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

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
var AemberspineMongrel = card.New(
	"Aemberspine Mongrel",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 335),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithHazardous(3),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.GainAember{
			Player: card.Controller,
			Amount: 1,
		}),
)
