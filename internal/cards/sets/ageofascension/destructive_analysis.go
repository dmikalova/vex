package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Destructive Analysis
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 2 damage to a Creature and purge any number of cards from your archives, and for each card purged this way, deal 2 damage to it.
var DestructiveAnalysis = card.New(
	"Destructive Analysis",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "194"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			After:  card.Always,
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.PurgeArchives{},
				card.DealDamage{
					Target: card.Target.Triggering,
					Amount: 2,
					Per:    card.CardsPurged{},
				},
			}},
		}),
)
