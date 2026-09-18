package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Soldiers to Flowers
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Purge each Untamed creature from each player's discard pile. For each card purged this way, its owner gains 1 Æmber.
var SoldiersToFlowers = set.New(
	"Soldiers to Flowers",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "349"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PurgeCard{
			Zone:   card.Discard,
			Player: card.EachPlayer,
			Selection: card.Each{
				House: card.Houses.Named(card.House.Self),
				Type:  card.Type.Creature,
			},
			GainOwnerAember: true,
		}),
)
