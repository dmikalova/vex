package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bait and Switch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber -> if your opponent has more Æmber than you, repeat this effect.
var BaitAndSwitch = set.New(
	"Bait and Switch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "267"),
	card.WithAbility(
		card.Trigger.Play, card.Repeat{
			Do:   card.StealAember{Amount: 1},
			Gate: card.While{Cond: card.PoolAember{Player: card.Opponent, Is: card.MoreThanYou}},
		}),
)
