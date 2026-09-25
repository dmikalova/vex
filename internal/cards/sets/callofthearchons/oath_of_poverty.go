package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Oath of Poverty
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each friendly artifact. For each card destroyed this way, gain 2 Æmber.
var OathOfPoverty = set.New(
	"Oath of Poverty",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "222"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Sequence{Effects: []card.Effect{
		card.Destroy{Target: card.Target.EachFriendlyArtifact},
		card.GainAember{
			Player: card.Controller,
			Amount: 2,
			Per:    card.CardsDestroyed{},
		},
	}}),
)
