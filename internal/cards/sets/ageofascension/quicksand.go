package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Quicksand
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy the most powerful creature controlled by each player who does not control a ready Untamed creature.
var Quicksand = card.New(
	"Quicksand",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "364"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.DestroyMostPowerfulUnlessReadyHouse{House: card.House.Self}),
)
