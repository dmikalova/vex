package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Information Exchange
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber. If your opponent stole Æmber from you on their previous turn, steal 1 Æmber.
var InformationExchange = set.New(
	"Information Exchange",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "136"),
	card.WithAbility(
		card.Trigger.Play,
		card.Sequence{Effects: []card.Effect{
			card.StealAember{Amount: 1},
			card.Conditional{
				Cond: card.AemberStolenFromYou{},
				Then: card.StealAember{Amount: 1},
			},
		}},
	),
)
