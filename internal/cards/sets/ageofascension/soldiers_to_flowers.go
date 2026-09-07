package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Soldiers to Flowers
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Purge each Untamed creature from each player's discard pile. For each card purged this way, its owner gains 1 Æmber.
var SoldiersToFlowers = card.New(
	"Soldiers to Flowers",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "349"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PurgeEachFromDiscard{
			House:           card.House.Self,
			Type:            card.Type.Creature,
			GainOwnerAember: true,
		}),
)
