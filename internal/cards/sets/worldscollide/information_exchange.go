package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// InformationExchange
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1A. If your opponent stole A from you on their previous turn, steal 2A instead.
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
