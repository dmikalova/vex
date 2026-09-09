package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Information Exchange
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber, or 2 if your opponent stole Æmber from you on their previous turn.
var InformationExchange = card.New(
	"Information Exchange",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "136"),
	card.WithAbility(
		card.Trigger.Play,
		card.StealAember{
			Amount: 1,
			Or: card.OrAmount{
				Amount: 2,
				When:   card.AemberStolenFromYou{},
			},
		},
	),
)
