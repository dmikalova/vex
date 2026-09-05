package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hecatomb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each Dis creature. For each creature they controlled that was destroyed this way, each player gains 1 Æmber.
var Hecatomb = card.New(
	"Hecatomb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 63),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Destroy{
				Target: card.Target.EachCreature.OfHouse(card.House.Self),
			},
			card.GainAember{
				Player: card.EachPlayer,
				Amount: 1,
				Per:    card.CreaturesDestroyedThisWay{Player: card.Controller},
			},
		}}),
)
