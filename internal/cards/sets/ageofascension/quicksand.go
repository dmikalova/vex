package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Quicksand
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy the most powerful Creature controlled by each player who does not have a friendly ready Untamed Creature in play.
var Quicksand = set.New(
	"Quicksand",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "364"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.BatchDestroy{Gather: card.EachPlayerUnless{
			Spare: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
				Ready:  true,
			},
			Take: card.MostPowerfulN(1),
		}}),
)
