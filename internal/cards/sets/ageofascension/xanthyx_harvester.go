package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Xanthyx Harvester
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	You cannot use this card unless Xanthyx Harvester has no non-Mars neighbor.
//	Reap: Gain 1 Æmber.
var XanthyxHarvester = set.New(
	"Xanthyx Harvester",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "173"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithRestrictions(card.Restrictions{
		UseCondition: card.SourceHasNoNeighbor{
			House: card.Houses.Except(card.House.Self),
		},
	}),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{
			Player: card.Controller,
			Amount: 1,
		}),
)
