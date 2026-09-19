package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Destructive Analysis
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature and purge any number of cards from your archives. For each card purged this way, deal 2 damage to the same creature.
var DestructiveAnalysis = set.New(
	"Destructive Analysis",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "194"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			After:  card.Always,
			Target: card.Target.Creature,
			Then: card.Sentences{Effects: []card.Effect{
				card.PurgeCard{
					Zones:     []card.Zone{card.Archives},
					Player:    card.Controller,
					Selection: card.Chosen{Optional: true},
					Quantity:  card.AnyNumber{},
				},
				card.DealDamage{
					Target: card.Target.TheSameCreature,
					Amount: 2,
					Per:    card.CardsPurged{},
				},
			}},
		}),
)
