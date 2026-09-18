package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Quicksand
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy the most powerful creature controlled by each player who does not have a friendly ready Untamed creature in play.
var Quicksand = set.New(
	"Quicksand",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "364"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.BatchDestroy{Gather: card.EachPlayerUnless{
			Spare: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.Houses.Named(card.House.Self),
				Ready:  true,
			},
			Take: card.MostPowerfulN(1),
		}}),
)
