package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bait and Switch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: If your opponent has more Æmber than you, steal 1 Æmber -> repeat this effect.
var BaitAndSwitch = card.New(
	"Bait and Switch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 267),
	card.WithAbility(
		card.Trigger.Play, card.RepeatWhile{
			Cond: card.OpponentAember{Is: card.MoreThanYou},
			Do:   card.StealAember{Amount: 1},
		}),
)
