package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Neuro Syphon
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, steal 1 Æmber, and draw a card.
var NeuroSyphon = card.New(
	"Neuro Syphon",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 116),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAemberMoreThanYou{},
			Then: card.Sequence{
				Effects: []card.Effect{
					card.StealAember{Amount: 1},
					card.Draw{Amount: 1},
				},
			},
		}),
)
