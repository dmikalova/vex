package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Gravid Cycle
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Omega.
//	Play: Put a card from your discard pile into your hand.
var GravidCycle = set.New(
	"Gravid Cycle",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "354"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play,
		card.PutCard{
			Zones:       []card.Zone{card.Discard},
			Selection:   card.Chosen{},
			Destination: card.To.Hand,
		},
	),
)
