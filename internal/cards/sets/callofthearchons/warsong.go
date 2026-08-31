package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Warsong
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time a friendly creature fights, gain 1 Æmber.
var Warsong = card.New(
	"Warsong",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 18),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.Fight,
			Do: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
