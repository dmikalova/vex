package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Take Hostages
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time a friendly creature fights, it captures 1 Æmber from your opponent.
var TakeHostages = card.New(
	"Take Hostages",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 226),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.Fight,
			Do: card.CaptureAember{
				Amount: 1,
				Target: card.Target.Triggering,
				Source: card.Opponent,
			},
		}),
)
