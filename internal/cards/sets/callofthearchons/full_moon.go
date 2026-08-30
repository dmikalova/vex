package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Full Moon
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time you play a creature, gain 1 Æmber.
var FullMoon = card.New(
	"Full Moon",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 323),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.CreaturePlayed,
			Do: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
