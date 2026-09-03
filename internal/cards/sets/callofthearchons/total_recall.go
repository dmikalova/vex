package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Total Recall
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For each friendly ready creature in play, gain 1 Æmber. Put each friendly creature into its owner's hand.
var TotalRecall = card.New(
	"Total Recall",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 179),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sentences{Effects: []card.Effect{
		card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Ready:  true,
			},
		},
		card.PutFromPlay{
			Target:      card.Target.EachFriendlyCreature,
			Destination: card.To.Hand,
		},
	}}),
)
