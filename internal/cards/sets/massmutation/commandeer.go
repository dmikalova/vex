package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Commandeer
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, each time you play another card, a friendly creature captures 1 Æmber from your opponent.
var Commandeer = set.New(
	"Commandeer",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "131"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.CardPlayed,
			Do: card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
		}),
)
