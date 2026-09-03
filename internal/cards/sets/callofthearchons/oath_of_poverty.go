package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Oath of Poverty
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each friendly artifact. For each card destroyed this way, gain 2 Æmber.
var OathOfPoverty = card.New(
	"Oath of Poverty",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 222),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sentences{Effects: []card.Effect{
		card.Destroy{Target: card.Target.EachFriendlyArtifact},
		card.GainAember{
			Player: card.Controller,
			Amount: 2,
			Per:    card.CardsDestroyed{},
		},
	}}),
)
