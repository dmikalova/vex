package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Evil Eye
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Keys cost +3 Æmber during your opponent's next turn.
var TheEvilEye = set.New(
	"The Evil Eye",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "84"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.RaiseKeyCost{
			Player:   card.Opponent,
			Amount:   3,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
