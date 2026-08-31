package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ghostly Hand
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  2
//
//	Play: If your opponent has exactly 1 Æmber, steal 1 Æmber.
var GhostlyHand = card.New(
	"Ghostly Hand",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 270),
	card.WithAemberBonus(2),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.Exactly, Amount: 1},
			Then: card.StealAember{Amount: 1},
		}),
)
