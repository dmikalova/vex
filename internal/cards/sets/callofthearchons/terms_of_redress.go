package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Terms of Redress
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: A friendly creature captures 2 Æmber from your opponent.
var TermsOfRedress = card.New(
	"Terms of Redress",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 227),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 2,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
		}),
)
