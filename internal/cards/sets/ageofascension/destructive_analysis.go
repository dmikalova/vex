package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Destructive Analysis
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 2 damage to a creature and purge any number of cards from your archives to deal an additional 2 damage to it for each card purged this way.
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
			Target: card.Target.Creature,
			Then: card.PurgeArchivesForDamage{
				Amount: 2,
				Target: card.Target.Triggering,
			},
		}),
)

// TODO: why is damage then needed instead of a sequence. PurgeArchivesForDamage is ridiculous
