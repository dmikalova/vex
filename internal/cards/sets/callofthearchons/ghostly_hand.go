package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ghostly Hand
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber Æmber
//
//	Play: If your opponent has exactly 1 Æmber, steal 1 Æmber.
var GhostlyHand = set.New(
	"Ghostly Hand",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "270"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.Exactly,
				Amount: 1,
			},
			Then: card.StealAember{Amount: 1},
		}),
)
