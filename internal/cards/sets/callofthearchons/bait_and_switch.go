package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bait and Switch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber -> if your opponent has more Æmber than you, repeat this effect.
var BaitAndSwitch = card.New(
	"Bait and Switch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "267"),
	card.WithAbility(
		card.Trigger.Play, card.RepeatOnCondition{
			Do:   card.StealAember{Amount: 1},
			Cond: card.PoolAember{Player: card.Opponent, Is: card.MoreThanYou},
		}),
)
